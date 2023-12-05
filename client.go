package groupme

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

// GroupMeAPIBase - Endpoints are added on to this to get the full URI.
// Overridable for testing
const GroupMeAPIBase = "https://api.groupme.com/v3"
const GroupMeImageBase = "https://image.groupme.com"

type logger struct {
	*slog.Logger
	requestToAttr  func(*http.Request) slog.Attr
	responseToAttr func(*http.Response) slog.Attr
}

func requestToAttr(req *http.Request) slog.Attr {
	return slog.Group(
		"request",
		"method", req.Method,
		"url", req.URL.String(),
	)
}

func responseToAttr(resp *http.Response) slog.Attr {
	return slog.Group(
		"response",
		slog.Group(
			"header",
			"content-length", resp.Header.Get("Content-Length"),
			"content-type", resp.Header.Get("Content-Type"),
		),
		slog.Group(
			"request",
			"method", resp.Request.Method,
			"url", resp.Request.URL.String(),
		),
	)
}

// Client communicates with the GroupMe API to perform actions
// on the basic types, i.e. Listing, Creating, Destroying
type Client struct {
	httpClient         *http.Client
	logger             *logger
	apiEndpointBase    string
	imageEndpointBase  string
	authorizationToken string
}

// NewClient creates a new GroupMe API Client
func NewClient(authToken string, options ...ClientOption) *Client {
	client := &Client{
		httpClient: &http.Client{},
		logger: &logger{
			Logger:         slog.New(new(noopSlogHandler)),
			requestToAttr:  requestToAttr,
			responseToAttr: responseToAttr,
		},
		apiEndpointBase:    GroupMeAPIBase,
		imageEndpointBase:  GroupMeImageBase,
		authorizationToken: authToken,
	}

	for _, option := range options {
		option(client)
	}

	return client
}

type ClientOption func(client *Client)

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

func WithLogHander(handler slog.Handler) ClientOption {
	return func(client *Client) {
		client.logger.Logger = slog.New(handler)
	}
}

func WithRequestToAttr(requestToValue func(*http.Request) slog.Attr) ClientOption {
	return func(client *Client) {
		client.logger.requestToAttr = requestToValue
	}
}

func WithResponseToAttr(responseToValue func(*http.Response) slog.Attr) ClientOption {
	return func(client *Client) {
		client.logger.responseToAttr = responseToValue
	}
}

type noopSlogHandler struct{}

func (h *noopSlogHandler) Enabled(context.Context, slog.Level) bool  { return true }
func (h *noopSlogHandler) Handle(context.Context, slog.Record) error { return nil }
func (h *noopSlogHandler) WithAttrs(attra []slog.Attr) slog.Handler  { return h }
func (h *noopSlogHandler) WithGroup(name string) slog.Handler        { return h }

// Close safely shuts down the Client
func (c *Client) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}

// String returns a json formatted string
func (c Client) String() string {
	return marshal(&c)
}

/*/// Handle parsing of nested interface type response ///*/
type jsonResponse struct {
	Response response `json:"response"`
	Meta     `json:"meta"`
}

func newJSONResponse(i interface{}) *jsonResponse {
	return &jsonResponse{Response: response{i}}
}

type response struct {
	i interface{}
}

func (r response) UnmarshalJSON(bs []byte) error {
	return json.NewDecoder(bytes.NewBuffer(bs)).Decode(r.i)
}

const errorStatusCodeMin = 300

func (c Client) do(ctx context.Context, req *http.Request, i interface{}) error {
	req = req.WithContext(ctx)
	if req.Method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	c.logger.Info("sending", c.logger.requestToAttr(req))
	resp, err := c.httpClient.Do(req)
	c.logger.Info("received response", c.logger.responseToAttr(resp))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var readBytes []byte
	// Check Status Code is 1XX or 2XX
	if resp.StatusCode >= errorStatusCodeMin {
		readBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			// We couldn't read the output.  Oh well; generate the appropriate error type anyway.
			return &Meta{
				Code: resp.StatusCode,
			}
		}

		jsonResp := newJSONResponse(nil)
		if err = json.Unmarshal(readBytes, &jsonResp); err != nil {
			// We couldn't parse the output.  Oh well; generate the appropriate error type anyway.
			return &Meta{
				Code: resp.StatusCode,
			}
		}
		return &jsonResp.Meta
	}

	if i == nil {
		return nil
	}

	readBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	jsonResp := newJSONResponse(i)
	if err := json.Unmarshal(readBytes, &jsonResp); err != nil {
		return err
	}

	return nil
}

func (c Client) doWithAuthToken(ctx context.Context, req *http.Request, i interface{}) error {
	URL := req.URL
	query := URL.Query()
	query.Set("token", c.authorizationToken)
	URL.RawQuery = query.Encode()

	return c.do(ctx, req, i)
}
