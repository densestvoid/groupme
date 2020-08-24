package groupme

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

// GroupMe documentation: https://dev.groupme.com/docs/v3#messages

////////// Endpoints //////////
const (
	// Used to build other endpoints
	messagesEndpointRoot = groupEndpointRoot + "/messages"

	indexMessagesEndpoint  = messagesEndpointRoot // GET
	createMessagesEndpoint = messagesEndpointRoot // POST
)

// Index

// MessagesQuery defines the optional URL parameters for IndexMessages
type IndexMessagesQuery struct {
	// Returns messages created before the given message ID
	BeforeID ID
	// Returns most recent messages created after the given message ID
	SinceID ID
	// Returns messages created immediately after the given message ID
	AfterID ID
	// Number of messages returned. Default is 20. Max is 100.
	Limit int
}

func (q IndexMessagesQuery) String() string {
	return marshal(&q)
}

// MessagesIndexResponse contains the count and set of
// messages returned by the IndexMessages API request
type IndexMessagesResponse struct {
	Count    int        `json:"count"`
	Messages []*Message `json:"messages"`
}

func (r IndexMessagesResponse) String() string {
	return marshal(&r)
}

/*
IndexMessages -

Retrieve messages for a group.

By default, messages are returned in groups of 20, ordered by
created_at descending. This can be raised or lowered by passing
a limit parameter, up to a maximum of 100 messages.

Messages can be scanned by providing a message ID as either the
before_id, since_id, or after_id parameter. If before_id is
provided, then messages immediately preceding the given message
will be returned, in descending order. This can be used to
continually page back through a group's messages.

The after_id parameter will return messages that immediately
follow a given message, this time in ascending order (which
makes it easy to pick off the last result for continued
pagination).

Finally, the since_id parameter also returns messages created
after the given message, but it retrieves the most recent
messages. For example, if more than twenty messages are created
after the since_id message, using this parameter will omit the
messages that immediately follow the given message. This is a
bit counterintuitive, so take care.

If no messages are found (e.g. when filtering with before_id)
we return code 304.

Note that for historical reasons, likes are returned as an
array of user ids in the favorited_by key.

Parameters: See MessageQuery
*/
func (c *Client) IndexMessages(groupID ID, req *IndexMessagesQuery) (IndexMessagesResponse, error) {
	url := fmt.Sprintf(c.endpointBase+indexMessagesEndpoint, groupID)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return IndexMessagesResponse{}, err
	}

	URL := httpReq.URL
	query := URL.Query()
	if req != nil {
		if req.BeforeID != "" {
			query.Add("before_id", req.BeforeID.String())
		}
		if req.SinceID != "" {
			query.Add("since_id", req.SinceID.String())
		}
		if req.AfterID != "" {
			query.Add("after_id", req.AfterID.String())
		}
		if req.Limit != 0 {
			query.Add("limit", strconv.Itoa(req.Limit))
		}
	}
	URL.RawQuery = query.Encode()

	var resp IndexMessagesResponse
	err = c.doWithAuthToken(httpReq, &resp)
	if err != nil {
		return IndexMessagesResponse{}, err
	}

	return resp, nil
}

// Create

/*
CreateMessage -
Send a message to a group

If you want to attach an image, you must first process it
through our image service.

Attachments of type emoji rely on data from emoji PowerUps.

Clients use a placeholder character in the message text and
specify a replacement charmap to substitute emoji characters

The character map is an array of arrays containing rune data
([[{pack_id,offset}],...]).

The placeholder should be a high-point/invisible UTF-8 character.

Parameters:
	groupID - required, ID(String)
	See Message.
		text - required, string. Can be ommitted if at least one
			attachment is present
		attachments - a polymorphic list of attachments (locations,
			images, etc). You may have You may have more than
			one of any type of attachment, provided clients can
			display it.

*/
func (c *Client) CreateMessage(groupID ID, m *Message) (*Message, error) {
	URL := fmt.Sprintf(c.endpointBase+createMessagesEndpoint, groupID)

	m.SourceGUID = uuid.New().String()
	var data = struct {
		Message *Message `json:"message"`
	}{
		m,
	}

	jsonBytes, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	var resp struct {
		*Message `json:"message"`
	}
	err = c.doWithAuthToken(httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Message, nil
}
