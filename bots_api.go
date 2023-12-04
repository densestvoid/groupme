package groupme

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GroupMe documentation: https://dev.groupme.com/docs/v3#bots

/*//////// Endpoints ////////*/
const (
	// Used to build other endpoints
	botsEndpointRoot = "/bots"

	// Actual Endpoints
	createBotEndpoint      = botsEndpointRoot              // POST
	postBotMessageEndpoint = botsEndpointRoot + "/post"    // POST
	indexBotsEndpoint      = botsEndpointRoot              // GET
	destroyBotEndpoint     = botsEndpointRoot + "/destroy" // POST
)

/*//////// API Requests ////////*/

// CreateBot - Create a bot. See the Bots Tutorial (https://dev.groupme.com/tutorials/bots)
// for a full walkthrough.
func (c *Client) CreateBot(ctx context.Context, bot *Bot) (*Bot, error) {
	URL := c.endpointBase + createBotEndpoint

	var data = struct {
		Bot *Bot `json:"bot,omitempty"`
	}{
		bot,
	}

	jsonBytes, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	var resp Bot
	err = c.doWithAuthToken(ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// PostBotMessage - Post a message from a bot
// TODO: Move PostBotMessage to bot object, since it doesn't require access token
func (c *Client) PostBotMessage(ctx context.Context, botID ID, text string, pictureURL *string) error {
	URL := fmt.Sprintf(c.endpointBase + postBotMessageEndpoint)

	var data = struct {
		BotID      ID      `json:"bot_id"`
		Text       string  `json:"text"`
		PictureURL *string `json:",omitempty"`
	}{
		botID,
		text,
		pictureURL,
	}

	jsonBytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}

	return c.do(ctx, httpReq, nil)
}

// IndexBots - list bots that you have created
func (c *Client) IndexBots(ctx context.Context) ([]*Bot, error) {
	httpReq, err := http.NewRequest("GET", c.endpointBase+indexBotsEndpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp []*Bot
	err = c.doWithAuthToken(ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DestroyBot - Remove a bot that you have created
func (c *Client) DestroyBot(ctx context.Context, botID ID) error {
	URL := fmt.Sprintf(c.endpointBase + destroyBotEndpoint)

	var data = struct {
		BotID ID `json:"bot_id"`
	}{
		botID,
	}

	jsonBytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}

	return c.doWithAuthToken(ctx, httpReq, nil)
}
