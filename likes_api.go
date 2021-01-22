// Package groupme defines a client capable of executing API commands for the GroupMe chat service
package groupme

import (
	"context"
	"fmt"
	"net/http"
)

// GroupMe documentation: https://dev.groupme.com/docs/v3#likes

/*//////// Endpoints ////////*/
const (
	// Used to build other endpoints
	likesEndpointRoot = "/messages/%s/%s"

	createLikeEndpoint  = likesEndpointRoot + "/like"   // POST
	destroyLikeEndpoint = likesEndpointRoot + "/unlike" // POST
)

/*//////// API Requests ////////*/

// Create

// CreateLike - Like a message.
func (c *Client) CreateLike(ctx context.Context, conversationID, messageID ID) error {
	url := fmt.Sprintf(c.endpointBase+createLikeEndpoint, conversationID, messageID)

	httpReq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	return c.doWithAuthToken(ctx, httpReq, nil)
}

// DestroyLike - Unlike a message.
func (c *Client) DestroyLike(ctx context.Context, conversationID, messageID ID) error {
	url := fmt.Sprintf(c.endpointBase+destroyLikeEndpoint, conversationID, messageID)

	httpReq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	return c.doWithAuthToken(ctx, httpReq, nil)
}
