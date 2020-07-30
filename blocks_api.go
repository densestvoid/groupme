package groupme

import (
	"net/http"
)

// GroupMe documentation: https://dev.groupme.com/docs/v3#blocks

////////// Endpoints //////////
const (
	// Used to build other endpoints
	blocksEndpointRoot = "/blocks"

	// Actual Endpoints
	indexBlocksEndpoint  = blocksEndpointRoot              // GET
	blockBetweenEndpoint = blocksEndpointRoot + "/between" // GET
	createBlockEndpoint  = blocksEndpointRoot              // POST
	unblockEndpoint      = blocksEndpointRoot              // DELETE
)

////////// API Requests //////////

// Index

/*
IndexBlock -

A list of contacts you have blocked. These people cannot DM you

Parameters:
	userID - required, ID(string)
*/
func (c *Client) IndexBlock(userID ID) ([]*Block, error) {
	httpReq, err := http.NewRequest("GET", c.endpointBase+indexBlocksEndpoint, nil)
	if err != nil {
		return nil, err
	}

	URL := httpReq.URL
	query := URL.Query()
	query.Set("user", userID.String())
	URL.RawQuery = query.Encode()

	var resp struct {
		Blocks []*Block `json:"blocks"`
	}
	err = c.do(httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Blocks, nil
}

// Between

/*
BlockBetween -

Asks if a block exists between you and another user id

Parameters:
	otherUserID - required, ID(string)
*/
func (c *Client) BlockBetween(userID, otherUserID ID) (bool, error) {
	httpReq, err := http.NewRequest("GET", c.endpointBase+blockBetweenEndpoint, nil)
	if err != nil {
		return false, err
	}

	URL := httpReq.URL
	query := URL.Query()
	query.Set("user", userID.String())
	query.Set("otherUser", otherUserID.String())
	URL.RawQuery = query.Encode()

	var resp struct {
		Between bool `json:"between"`
	}
	err = c.do(httpReq, &resp)
	if err != nil {
		return false, err
	}

	return resp.Between, nil
}

// Create

/*
CreateBlock -

Creates a block between you and the contact

Parameters:
	userID - required, ID(string)
	otherUserID - required, ID(string)
*/
func (c *Client) CreateBlock(userID, otherUserID ID) (*Block, error) {
	httpReq, err := http.NewRequest("POST", c.endpointBase+createBlockEndpoint, nil)
	if err != nil {
		return nil, err
	}

	URL := httpReq.URL
	query := URL.Query()
	query.Set("user", userID.String())
	query.Set("otherUser", otherUserID.String())
	URL.RawQuery = query.Encode()

	var resp struct {
		Block *Block `json:"block"`
	}
	err = c.do(httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Block, nil
}

// Unblock

/*
Unblock -

Removes block between you and other user

Parameters:
	userID - required, ID(string)
	otherUserID - required, ID(string)
*/
func (c *Client) Unblock(userID, otherUserID ID) error {
	httpReq, err := http.NewRequest("DELETE", c.endpointBase+unblockEndpoint, nil)
	if err != nil {
		return err
	}

	URL := httpReq.URL
	query := URL.Query()
	query.Set("user", userID.String())
	query.Set("otherUser", otherUserID.String())
	URL.RawQuery = query.Encode()

	err = c.do(httpReq, nil)
	if err != nil {
		return err
	}

	return nil
}
