package groupme

import (
	"fmt"
	"net/http"
)

// GroupMe documentation: https://dev.groupme.com/docs/v3#likes

////////// Endpoints //////////
const (
	// Used to build other endpoints
	likesEndpointRoot = "/messages/%s/%s"

	createLikeEndpoint  = likesEndpointRoot + "/like"   // POST
	destroyLikeEndpoint = likesEndpointRoot + "/unlike" // POST
)

////////// API Requests /////////

// Create

/*
CreateLike -

Like a message.

Parameters:
	conversationID - required, ID(string)
	messageID - required, ID(string)
*/
func (c *Client) CreateLike(conversationID, messageID ID) error {
	url := fmt.Sprintf(c.endpointBase+createLikeEndpoint, conversationID, messageID)

	httpReq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	return c.do(httpReq, nil)
}

// Destroy

/*
DestroyLike -

Unlike a message.

Parameters:
	conversationID - required, ID(string)
	messageID - required, ID(string)
*/
func (c *Client) DestroyLike(conversationID, messageID ID) error {
	url := fmt.Sprintf(c.endpointBase+destroyLikeEndpoint, conversationID, messageID)

	httpReq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	return c.do(httpReq, nil)
}
