package groupme

import (
	"context"
	"net/http"
)

// GroupMe documentation: https://dev.groupme.com/docs/v3#blocks

/*//////// Endpoints ////////*/
const (
	// Used to build other endpoints
	blocksEndpointRoot = "/blocks"

	// Actual Endpoints
	indexBlocksEndpoint  = blocksEndpointRoot              // GET
	blockBetweenEndpoint = blocksEndpointRoot + "/between" // GET
	createBlockEndpoint  = blocksEndpointRoot              // POST
	unblockEndpoint      = blocksEndpointRoot              // DELETE
)

/*//////// API Requests ////////*/

// IndexBlock - A list of contacts you have blocked. These people cannot DM you
func (c *Client) IndexBlock(ctx context.Context, userID string) ([]*Block, error) {
	httpReq, err := http.NewRequest("GET", c.apiEndpointBase+indexBlocksEndpoint, nil)
	if err != nil {
		return nil, err
	}

	URL := httpReq.URL
	query := URL.Query()
	query.Set("user", userID)
	URL.RawQuery = query.Encode()

	var resp struct {
		Blocks []*Block `json:"blocks"`
	}
	err = c.doWithAuthToken(ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Blocks, nil
}

// BlockBetween - Asks if a block exists between you and another user id
func (c *Client) BlockBetween(ctx context.Context, userID, otherUserID string) (bool, error) {
	httpReq, err := http.NewRequest("GET", c.apiEndpointBase+blockBetweenEndpoint, nil)
	if err != nil {
		return false, err
	}

	URL := httpReq.URL
	query := URL.Query()
	query.Set("user", userID)
	query.Set("otherUser", otherUserID)
	URL.RawQuery = query.Encode()

	var resp struct {
		Between bool `json:"between"`
	}
	err = c.doWithAuthToken(ctx, httpReq, &resp)
	if err != nil {
		return false, err
	}

	return resp.Between, nil
}

// CreateBlock - Creates a block between you and the contact
func (c *Client) CreateBlock(ctx context.Context, userID, otherUserID string) (*Block, error) {
	httpReq, err := http.NewRequest("POST", c.apiEndpointBase+createBlockEndpoint, nil)
	if err != nil {
		return nil, err
	}

	URL := httpReq.URL
	query := URL.Query()
	query.Set("user", userID)
	query.Set("otherUser", otherUserID)
	URL.RawQuery = query.Encode()

	var resp struct {
		Block *Block `json:"block"`
	}
	err = c.doWithAuthToken(ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Block, nil
}

// Unblock - Removes block between you and other user
func (c *Client) Unblock(ctx context.Context, userID, otherUserID string) error {
	httpReq, err := http.NewRequest("DELETE", c.apiEndpointBase+unblockEndpoint, nil)
	if err != nil {
		return err
	}

	URL := httpReq.URL
	query := URL.Query()
	query.Set("user", userID)
	query.Set("otherUser", otherUserID)
	URL.RawQuery = query.Encode()

	err = c.doWithAuthToken(ctx, httpReq, nil)
	if err != nil {
		return err
	}

	return nil
}
