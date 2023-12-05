package groupme

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GroupMe documentation: https://dev.groupme.com/docs/v3#sms_mode

/*//////// Endpoints ////////*/
const (
	// Used to build other endpoints
	smsModeEndpointRoot = usersEndpointRoot + "/sms_mode"

	// Actual Endpoints
	createSMSModeEndpoint = smsModeEndpointRoot             // POST
	deleteSMSModeEndpoint = smsModeEndpointRoot + "/delete" // POST
)

/*//////// API Requests ////////*/

// Create

/*
CreateSMSMode -
Enables SMS mode for N hours, where N is at most 48. After N
hours have elapsed, user will receive push notfications.

Parameters:

	duration - required, integer
	registration_id - string; The push notification ID/token
		that should be suppressed during SMS mode. If this is
		omitted, both SMS and push notifications will be
		delivered to the device.
*/
func (c *Client) CreateSMSMode(ctx context.Context, duration int, registrationID *ID) error {
	URL := fmt.Sprintf(c.apiEndpointBase + createSMSModeEndpoint)

	var data = struct {
		Duration       int `json:"duration"`
		RegistrationID *ID `json:"registration_id,omitempty"`
	}{
		duration,
		registrationID,
	}

	jsonBytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}

	err = c.doWithAuthToken(ctx, httpReq, nil)
	if err != nil {
		return err
	}

	return nil
}

// Delete

/*
DeleteSMSMode -

Disables SMS mode
*/
func (c *Client) DeleteSMSMode(ctx context.Context) error {
	url := fmt.Sprintf(c.apiEndpointBase + deleteSMSModeEndpoint)

	httpReq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	return c.doWithAuthToken(ctx, httpReq, nil)
}
