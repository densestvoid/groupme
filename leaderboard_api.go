package groupme

import (
	"fmt"
	"net/http"
)

// GroupMe documentation: https://dev.groupme.com/docs/v3#leaderboard

////////// Endpoints //////////
const (
	// Used to build other endpoints
	leaderboardEndpointRoot = groupEndpointRoot + "/likes"

	// Actual Endpoints
	indexLeaderboardEndpoint   = leaderboardEndpointRoot             // GET
	myLikesLeaderboardEndpoint = leaderboardEndpointRoot + "/mine"   // GET
	myHitsLeaderboardEndpoint  = leaderboardEndpointRoot + "/for_me" // GET
)

////////// API Requests //////////

// Index

type period string

func (p period) String() string {
	return string(p)
}

// Define acceptable period values
const (
	Period_Day   = "day"
	Period_Week  = "week"
	Period_Month = "month"
)

/*
IndexLeaderboard -

A list of the liked messages in the group for a given period of
time. Messages are ranked in order of number of likes.

Parameters:
	groupID - required, ID(string)
	p - required, period(string)
*/
func (c *Client) IndexLeaderboard(groupID ID, p period) ([]*Message, error) {
	url := fmt.Sprintf(c.endpointBase+indexLeaderboardEndpoint, groupID)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	URL := httpReq.URL
	query := URL.Query()
	query.Set("period", p.String())
	URL.RawQuery = query.Encode()

	var resp struct {
		Messages []*Message `json:"messages"`
	}
	err = c.do(httpReq, &resp)
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
func (c *Client) MyLikesLeaderboard(groupID ID) ([]*Message, error) {
	url := fmt.Sprintf(c.endpointBase+myLikesLeaderboardEndpoint, groupID)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Messages []*Message `json:"messages"`
	}
	err = c.do(httpReq, &resp)
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
func (c *Client) MyHitsLeaderboard(groupID ID) ([]*Message, error) {
	url := fmt.Sprintf(c.endpointBase+myHitsLeaderboardEndpoint, groupID)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Messages []*Message `json:"messages"`
	}
	err = c.do(httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Messages, nil
}
