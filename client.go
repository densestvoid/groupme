package groupme

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// GroupMeAPIBase - Endpoints are added on to this to get the full URI.
// Overridable for testing
const GroupMeAPIBase = "https://api.groupme.com/v3"
const GroupMeImageBase = "https://image.groupme.com"

// client communicates with the GroupMe API to perform actions
// on the basic types, i.e. Listing, Creating, Destroying
type client struct {
	httpClient        *http.Client
	apiEndpointBase   string
	imageEndpointBase string
}

type ClientOption func(client *client)

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(client *client) {
		client.httpClient = httpClient
	}
}

// Close safely shuts down the Client
func (c *client) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
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

func (c *client) do(ctx context.Context, req *http.Request, i interface{}) error {
	req = req.WithContext(ctx)
	if req.Method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
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
