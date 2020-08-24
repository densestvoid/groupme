package groupme

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Endpoints are added on to the GroupMeAPIBase to get the full URI.
// Overridable for testing
const GroupMeAPIBase = "https://api.groupme.com/v3"

// Client communicates with the GroupMe API to perform actions
// on the basic types, i.e. Listing, Creating, Destroying
type Client struct {
	httpClient         *http.Client
	endpointBase       string
	authorizationToken string
}

// NewClient creates a new GroupMe API Client
func NewClient(authToken string) *Client {
	return &Client{
		// TODO: enable transport information passing in
		httpClient:         &http.Client{},
		endpointBase:       GroupMeAPIBase,
		authorizationToken: authToken,
	}
}

// Close safely shuts down the Client
func (c *Client) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}

// String returns a json formatted string
func (c Client) String() string {
	return marshal(&c)
}

///// Handle parsing of nested interface type response /////
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

func (c Client) do(req *http.Request, i interface{}) error {
	getResp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer getResp.Body.Close()

	bytes, err := ioutil.ReadAll(getResp.Body)
	if err != nil {
		return err
	}

	// Check Status Code is 1XX or 2XX
	if getResp.StatusCode/100 > 2 {
		return fmt.Errorf("%s: %s", getResp.Status, string(bytes))
	}

	if i == nil {
		return nil
	}

	resp := newJSONResponse(i)
	if err := json.Unmarshal(bytes, &resp); err != nil {
		return err
	}

	// Check Status Code is 1XX or 2XX
	if resp.Meta.Code/100 > 2 {
		return &resp.Meta
	}

	return nil
}

func (c Client) doWithAuthToken(req *http.Request, i interface{}) error {
	URL := req.URL
	query := URL.Query()
	query.Set("token", c.authorizationToken)
	URL.RawQuery = query.Encode()

	return c.do(req, i)
}
