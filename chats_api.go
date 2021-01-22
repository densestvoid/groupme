// Package groupme defines a client capable of executing API commands for the GroupMe chat service
package groupme

import (
	"context"
	"net/http"
	"strconv"
)

// GroupMe documentation: https://dev.groupme.com/docs/v3#chats

/*//////// Endpoints ////////*/
const (
	// Used to build other endpoints
	chatsEndpointRoot = "/chats"

	indexChatsEndpoint = chatsEndpointRoot // GET
)

// IndexChatsQuery defines the optional URL parameters for IndexChats
type IndexChatsQuery struct {
	// Page Number
	Page int `json:"page"`
	// Number of chats per page
	PerPage int `json:"per_page"`
}

// IndexChats - Returns a paginated list of direct message chats, or
// conversations, sorted by updated_at descending.
func (c *Client) IndexChats(ctx context.Context, req *IndexChatsQuery) ([]*Chat, error) {
	httpReq, err := http.NewRequest("GET", c.endpointBase+indexChatsEndpoint, nil)
	if err != nil {
		return nil, err
	}

	URL := httpReq.URL
	query := URL.Query()
	if req != nil {
		if req.Page > 0 {
			query.Set("page", strconv.Itoa(req.Page))
		}
		if req.PerPage > 0 {
			query.Set("per_page", strconv.Itoa(req.PerPage))
		}
	}
	URL.RawQuery = query.Encode()

	var resp []*Chat
	err = c.doWithAuthToken(ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
