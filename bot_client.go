package groupme

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client posts bot messages
type BotClient struct {
	client
	botID string
}

// NewBotClient creates a new GroupMe API Client to post bot messages
func NewBotClient(botID string, options ...ClientOption) *BotClient {
	client := &BotClient{
		client: client{
			httpClient:        &http.Client{},
			apiEndpointBase:   GroupMeAPIBase,
			imageEndpointBase: GroupMeImageBase,
		},
		botID: botID,
	}

	for _, option := range options {
		option(&client.client)
	}

	return client
}

// PostBotMessage - Post a message from a bot
func (c *BotClient) PostBotMessage(ctx context.Context, text string, pictureURL *string) error {
	URL := fmt.Sprintf(c.apiEndpointBase + postBotMessageEndpoint)

	var data = struct {
		BotID      string  `json:"bot_id"`
		Text       string  `json:"text"`
		PictureURL *string `json:",omitempty"`
	}{
		c.botID,
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
