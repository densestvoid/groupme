package groupme

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

// GroupMe documentation: https://dev.groupme.com/docs/v3#direct_messages

////////// Endpoints //////////
const (
	// Used to build other endpoints
	directMessagesEndpointRoot = "/direct_messages"

	// Actual Endpoints
	indexDirectMessagesEndpoint = directMessagesEndpointRoot // GET
	createDirectMessageEndpoint = directMessagesEndpointRoot // POST
)

////////// API Requests //////////

// Index

// MessagesQuery defines the optional URL parameters for IndexDirectMessages
type IndexDirectMessagesQuery struct {
	// Returns 20 messages created before the given message ID
	BeforeID ID `json:"before_id"`
	// Returns 20 messages created after the given message ID
	SinceID ID `json:"since_id"`
}

func (q IndexDirectMessagesQuery) String() string {
	return marshal(&q)
}

// IndexDirectMessagesResponse contains the count and set of
// messages returned by the IndexDirectMessages API request
type IndexDirectMessagesResponse struct {
	Count    int        `json:"count"`
	Messages []*Message `json:"direct_messages"`
}

func (r IndexDirectMessagesResponse) String() string {
	return marshal(&r)
}

/*
IndexDirectMessages -

Fetch direct messages between two users.

DMs are returned in groups of 20, ordered by created_at
descending.

If no messages are found (e.g. when filtering with since_id) we
return code 304.

Note that for historical reasons, likes are returned as an array
of user ids in the favorited_by key.

Parameters:
	otherUserID - required, ID(string); the other participant in the conversation.
	See IndexDirectMessagesQuery
*/
func (c *Client) IndexDirectMessages(otherUserID ID, req *IndexDirectMessagesQuery) (IndexDirectMessagesResponse, error) {
	httpReq, err := http.NewRequest("GET", c.endpointBase+indexDirectMessagesEndpoint, nil)
	if err != nil {
		return IndexDirectMessagesResponse{}, err
	}

	query := httpReq.URL.Query()
	query.Set("other_user_id", otherUserID.String())
	if req != nil {
		if req.BeforeID != "" {
			query.Add("before_ID", req.BeforeID.String())
		}
		if req.SinceID != "" {
			query.Add("since_id", req.SinceID.String())
		}
	}

	var resp IndexDirectMessagesResponse
	err = c.do(httpReq, &resp)
	if err != nil {
		return IndexDirectMessagesResponse{}, err
	}

	return resp, nil
}

// Create

/*
CreateDirectMessage -

Send a DM to another user

If you want to attach an image, you must first process it
through our image service.

Attachments of type emoji rely on data from emoji PowerUps.

Clients use a placeholder character in the message text and
specify a replacement charmap to substitute emoji characters

The character map is an array of arrays containing rune data
([[{pack_id,offset}],...]).

Parameters:
	See Message.
		recipientID - required, ID(string); The GroupMe user ID of the recipient of this message.
		text - required, string. Can be ommitted if at least one
			attachment is present
		attachments - a polymorphic list of attachments (locations,
			images, etc). You may have You may have more than
			one of any type of attachment, provided clients can
			display it.
*/
func (c *Client) CreateDirectMessage(m *Message) (*Message, error) {
	URL := fmt.Sprintf(c.endpointBase + createDirectMessageEndpoint)
	httpReq, err := http.NewRequest("POST", URL, nil)
	if err != nil {
		return nil, err
	}

	m.SourceGUID = uuid.New().String()

	data := url.Values{}
	data.Set("direct_message", m.String())

	httpReq.PostForm = data

	var resp struct {
		*Message `json:"message"`
	}
	err = c.do(httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Message, nil
}
