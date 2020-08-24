package groupme

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GroupMe documentation: https://dev.groupme.com/docs/v3#users

////////// Endpoints //////////
const (
	// Used to build other endpoints
	usersEndpointRoot = "/users"

	// Actual Endpoints
	myUserEndpoint       = usersEndpointRoot + "/me"     // GET
	updateMyUserEndpoint = usersEndpointRoot + "/update" // POST
)

////////// API Requests //////////

// Me

/*
MyUser -

Loads a specific group.

Parameters:
	groupID - required, ID(string)
*/
func (c *Client) MyUser() (*User, error) {
	URL := fmt.Sprintf(c.endpointBase + myUserEndpoint)

	httpReq, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	var resp User
	err = c.doWithAuthToken(httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Update

type UserSettings struct {
	// URL to valid JPG/PNG/GIF image. URL will be converted into
	// an image service link (https://i.groupme.com/....)
	AvatarURL string `json:"avatar_url"`
	// Name must be of the form FirstName LastName
	Name string `json:"name"`
	// Email address. Must be in name@domain.com form
	Email   string `json:"email"`
	ZipCode string `json:"zip_code"`
}

/*
UpdateMyUser -

Update attributes about your own account

Parameters: See UserSettings
*/
func (c *Client) UpdateMyUser(us UserSettings) (*User, error) {
	URL := fmt.Sprintf(c.endpointBase + updateMyUserEndpoint)

	jsonBytes, err := json.Marshal(&us)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	var resp User
	err = c.doWithAuthToken(httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
