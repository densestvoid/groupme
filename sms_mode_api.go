package groupme

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// GroupMe documentation: https://dev.groupme.com/docs/v3#sms_mode

////////// Endpoints //////////
const (
	// Used to build other endpoints
	smsModeEndpointRoot = usersEndpointRoot + "/sms_mode"

	// Actual Endpoints
	createSMSModeEndpoint = smsModeEndpointRoot             // POST
	deleteSMSModeEndpoint = smsModeEndpointRoot + "/delete" // POST
)

////////// API Requests //////////

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
func (c *Client) CreateSMSMode(duration int, registrationID *ID) error {
	httpReq, err := http.NewRequest("POST", c.endpointBase+createSMSModeEndpoint, nil)
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Add("duration", strconv.Itoa(duration))

	if registrationID != nil {
		data.Add("registration_id", registrationID.String())
	}

	httpReq.PostForm = data

	err = c.doWithAuthToken(httpReq, nil)
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
func (c *Client) DeleteSMSMode() error {
	url := fmt.Sprintf(c.endpointBase + deleteSMSModeEndpoint)

	httpReq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	return c.doWithAuthToken(httpReq, nil)
}
