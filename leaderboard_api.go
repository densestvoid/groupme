// Package groupme defines a client capable of executing API commands for the GroupMe chat service
package groupme

import (
	"context"
	"fmt"
	"net/http"
)

// GroupMe documentation: https://dev.groupme.com/docs/v3#leaderboard

/*//////// Endpoints ////////*/
const (
	// Used to build other endpoints
	leaderboardEndpointRoot = groupEndpointRoot + "/likes"

	// Actual Endpoints
	indexLeaderboardEndpoint   = leaderboardEndpointRoot             // GET
	myLikesLeaderboardEndpoint = leaderboardEndpointRoot + "/mine"   // GET
	myHitsLeaderboardEndpoint  = leaderboardEndpointRoot + "/for_me" // GET
)

/*//////// API Requests ////////*/

// Index

type period string

// Define acceptable period values
const (
	PeriodDay   = "day"
	PeriodWeek  = "week"
	PeriodMonth = "month"
)

// IndexLeaderboard - A list of the liked messages in the group for a given period of
// time. Messages are ranked in order of number of likes.
func (c *Client) IndexLeaderboard(ctx context.Context, groupID ID, p period) ([]*Message, error) {
	url := fmt.Sprintf(c.endpointBase+indexLeaderboardEndpoint, groupID)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	URL := httpReq.URL
	query := URL.Query()
	query.Set("period", string(p))
	URL.RawQuery = query.Encode()

	var resp struct {
		Messages []*Message `json:"messages"`
	}
	err = c.doWithAuthToken(ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Messages, nil
}

// My Likes

/*
MyLikesLeaderboard -

A list of messages you have liked. Messages are returned in
reverse chrono-order. Note that the payload includes a liked_at
timestamp in ISO-8601 format.

Parameters:
	groupID - required, ID(string)
*/
func (c *Client) MyLikesLeaderboard(ctx context.Context, groupID ID) ([]*Message, error) {
	url := fmt.Sprintf(c.endpointBase+myLikesLeaderboardEndpoint, groupID)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Messages []*Message `json:"messages"`
	}
	err = c.doWithAuthToken(ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Messages, nil
}

// My Hits

/*
MyHitsLeaderboard -

A list of messages you have liked. Messages are returned in
reverse chrono-order. Note that the payload includes a liked_at
timestamp in ISO-8601 format.

Parameters:
	groupID - required, ID(string)
*/
func (c *Client) MyHitsLeaderboard(ctx context.Context, groupID ID) ([]*Message, error) {
	url := fmt.Sprintf(c.endpointBase+myHitsLeaderboardEndpoint, groupID)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Messages []*Message `json:"messages"`
	}
	err = c.doWithAuthToken(ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Messages, nil
}
