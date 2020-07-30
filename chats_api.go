package groupme

import (
	"net/http"
	"strconv"
)

// GroupMe documentation: https://dev.groupme.com/docs/v3#chats

////////// Endpoints //////////
const (
	// Used to build other endpoints
	chatsEndpointRoot = "/chats"

	indexChatsEndpoint = chatsEndpointRoot // GET
)

// Index

// ChatsQuery defines the optional URL parameters for IndexChats
type IndexChatsQuery struct {
	// Page Number
	Page int `json:"page"`
	// Number of chats per page
	PerPage int `json:"per_page"`
}

/*
IndexChats -

Returns a paginated list of direct message chats, or
conversations, sorted by updated_at descending.

Parameters: See ChatsQuery
*/
func (c *Client) IndexChats(req *IndexChatsQuery) ([]*Chat, error) {
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
	err = c.do(httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
