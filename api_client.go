package groupme

import (
	"context"
	"net/http"
)

// Client communicates with the GroupMe API to perform actions
// on the basic types, i.e. Listing, Creating, Destroying
type Client struct {
	client
	authorizationToken string
}

// NewClient creates a new GroupMe API Client
func NewClient(authToken string, options ...ClientOption) *Client {
	client := &Client{
		client: client{
			httpClient:        &http.Client{},
			apiEndpointBase:   GroupMeAPIBase,
			imageEndpointBase: GroupMeImageBase,
		},
		authorizationToken: authToken,
	}

	for _, option := range options {
		option(&client.client)
	}

	return client
}

func (c Client) doWithAuthToken(ctx context.Context, req *http.Request, i interface{}) error {
	URL := req.URL
	query := URL.Query()
	query.Set("token", c.authorizationToken)
	URL.RawQuery = query.Encode()

	return c.do(ctx, req, i)
}
