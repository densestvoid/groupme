// Package groupme defines a client capable of executing API commands for the GroupMe chat service
package groupme

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// GroupMe documentation: https://dev.groupme.com/docs/v3#groups

/*//////// Endpoints ////////*/
const (
	// Used to build other endpoints
	groupsEndpointRoot = "/groups"
	groupEndpointRoot  = "/groups/%s"

	// Actual Endpoints
	indexGroupsEndpoint      = groupsEndpointRoot                    // GET
	formerGroupsEndpoint     = groupsEndpointRoot + "/former"        // GET
	showGroupEndpoint        = groupEndpointRoot                     // GET
	createGroupEndpoint      = groupsEndpointRoot                    // POST
	updateGroupEndpoint      = groupEndpointRoot + "/update"         // POST
	destroyGroupEndpoint     = groupEndpointRoot + "/destroy"        // POST
	joinGroupEndpoint        = groupEndpointRoot + "/join/%s"        // POST
	rejoinGroupEndpoint      = groupsEndpointRoot + "/join"          // POST
	changeGroupOwnerEndpoint = groupsEndpointRoot + "/change_owners" // POST
)

/*//////// Common Request Parameters ////////*/

// GroupSettings is the settings for a group, used by CreateGroup and UpdateGroup
type GroupSettings struct {
	// Required. Primary name of the group. Maximum 140 characters
	Name string `json:"name"`
	// A subheading for the group. Maximum 255 characters
	Description string `json:"description"`
	// GroupMe Image Service URL
	ImageURL string `json:"image_url"`
	// Defaults false. If true, disables notifications for all members.
	// Documented for use only for UpdateGroup
	OfficeMode bool `json:"office_mode"`
	// Defaults false. If true, generates a share URL.
	// Anyone with the URL can join the group
	Share bool `json:"share"`
}

func (gss GroupSettings) String() string {
	return marshal(&gss)
}

/*//////// API Requests ////////*/

/*/// Index ///*/

// GroupsQuery defines optional URL parameters for IndexGroups
type GroupsQuery struct {
	// Fetch a particular page of results. Defaults to 1.
	Page int `json:"page"`
	// Define page size. Defaults to 10.
	PerPage int `json:"per_page"`
	// Comma separated list of data to omit from output.
	// Currently supported value is only "memberships".
	// If used then response will contain empty (null) members field.
	Omit string `json:"omit"`
}

func (q GroupsQuery) String() string {
	return marshal(&q)
}

/*
IndexGroups -

List the authenticated user's active groups.

The response is paginated, with a default of 10 groups per page.

Please consider using of omit=memberships parameter. Not including
member lists might significantly improve user experience of your
app for users who are participating in huge groups.

Parameters: See GroupsQuery
*/
func (c *Client) IndexGroups(ctx context.Context, req *GroupsQuery) ([]*Group, error) {
	httpReq, err := http.NewRequest("GET", c.endpointBase+indexGroupsEndpoint, nil)
	if err != nil {
		return nil, err
	}

	URL := httpReq.URL
	query := URL.Query()
	if req != nil {
		if req.Page != 0 {
			query.Set("page", strconv.Itoa(req.Page))
		}
		if req.PerPage != 0 {
			query.Set("per_page", strconv.Itoa(req.PerPage))
		}
		if req.Omit != "" {
			query.Set("omit", req.Omit)
		}
	}
	URL.RawQuery = query.Encode()

	var resp []*Group
	err = c.doWithAuthToken(ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

/*/// Former ///*/

/*
FormerGroups -

List they groups you have left but can rejoin.
*/
func (c *Client) FormerGroups(ctx context.Context) ([]*Group, error) {
	httpReq, err := http.NewRequest("GET", c.endpointBase+formerGroupsEndpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp []*Group
	err = c.doWithAuthToken(ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

/*/// Show ///*/

/*
ShowGroup -

Loads a specific group.

Parameters:
	groupID - required, ID(string)
*/
func (c *Client) ShowGroup(ctx context.Context, groupID ID) (*Group, error) {
	URL := fmt.Sprintf(c.endpointBase+showGroupEndpoint, groupID)

	httpReq, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	var resp Group
	err = c.doWithAuthToken(ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

/*/// Create ///*/

/*
CreateGroup -

Create a new group

Parameters: See GroupSettings
*/
func (c *Client) CreateGroup(ctx context.Context, gs GroupSettings) (*Group, error) {
	URL := fmt.Sprintf(c.endpointBase + createGroupEndpoint)

	jsonBytes, err := json.Marshal(&gs)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	var resp Group
	err = c.doWithAuthToken(ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

/*/// Update ///*/

/*
UpdateGroup -

Update a group after creation

Parameters:
	groupID - required, ID(string)
	See GroupSettings
*/
func (c *Client) UpdateGroup(ctx context.Context, groupID ID, gs GroupSettings) (*Group, error) {
	URL := fmt.Sprintf(c.endpointBase+updateGroupEndpoint, groupID)

	jsonBytes, err := json.Marshal(&gs)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	var resp Group
	err = c.doWithAuthToken(ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

/*/// Destroy ///*/

/*
DestroyGroup -

Disband a group

This action is only available to the group creator

Parameters:
	groupID - required, ID(string)
*/
func (c *Client) DestroyGroup(ctx context.Context, groupID ID) error {
	url := fmt.Sprintf(c.endpointBase+destroyGroupEndpoint, groupID)

	httpReq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	return c.doWithAuthToken(ctx, httpReq, nil)
}

/*/// Join ///*/

/*
JoinGroup -

Join a shared group

Parameters:
	groupID - required, ID(string)
	shareToken - required, string
*/
func (c *Client) JoinGroup(ctx context.Context, groupID ID, shareToken string) (*Group, error) {
	URL := fmt.Sprintf(c.endpointBase+joinGroupEndpoint, groupID, shareToken)

	httpReq, err := http.NewRequest("POST", URL, nil)
	if err != nil {
		return nil, err
	}

	var resp Group
	err = c.doWithAuthToken(ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

/*/// Rejoin ///*/

/*
RejoinGroup -

Rejoin a group. Only works if you previously removed yourself.

Parameters:
	groupID - required, ID(string)
*/
func (c *Client) RejoinGroup(ctx context.Context, groupID ID) (*Group, error) {
	URL := fmt.Sprintf(c.endpointBase + rejoinGroupEndpoint)

	var data = struct {
		GroupID ID `json:"group_id"`
	}{
		groupID,
	}

	jsonBytes, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	var resp Group
	err = c.doWithAuthToken(ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

/*/// Change Owner ///*/

/*
ChangeGroupOwner - Change owner of requested groups.

This action is only available to the group creator.

Response is a result object which contain status field,
the result of change owner action for the request

Parameters: See ChangeOwnerRequest
*/
func (c *Client) ChangeGroupOwner(ctx context.Context, reqs ChangeOwnerRequest) (ChangeOwnerResult, error) {
	URL := fmt.Sprintf(c.endpointBase + changeGroupOwnerEndpoint)

	var data = struct {
		Requests []ChangeOwnerRequest `json:"requests"`
	}{
		[]ChangeOwnerRequest{reqs},
	}

	jsonBytes, err := json.Marshal(&data)
	if err != nil {
		return ChangeOwnerResult{}, err
	}

	httpReq, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return ChangeOwnerResult{}, err
	}

	var resp struct {
		Results []ChangeOwnerResult `json:"results"`
	}

	err = c.doWithAuthToken(ctx, httpReq, &resp)
	if err != nil {
		return ChangeOwnerResult{}, err
	}

	if len(resp.Results) < 1 {
		return ChangeOwnerResult{}, errors.New("failed to parse results")
	}

	return resp.Results[0], nil
}

type changeOwnerStatusCode string

// Change owner Status Codes
const (
	ChangeOwnerOk                changeOwnerStatusCode = "200"
	ChangeOwnerRequesterNewOwner changeOwnerStatusCode = "400"
	ChangeOwnerNotOwner          changeOwnerStatusCode = "403"
	ChangeOwnerBadGroupOrOwner   changeOwnerStatusCode = "404"
	ChangeOwnerBadRequest        changeOwnerStatusCode = "405"
)

// String returns the description of the status code according to GroupMe
func (c changeOwnerStatusCode) String() string {
	return map[changeOwnerStatusCode]string{
		ChangeOwnerOk:                "success",
		ChangeOwnerRequesterNewOwner: "requester is also a new owner",
		ChangeOwnerNotOwner:          "requester is not the owner of the group",
		ChangeOwnerBadGroupOrOwner:   "group or new owner not found or new owner is not member of the group",
		ChangeOwnerBadRequest:        "request object is missing required field or any of the required fields is not an ID",
	}[c]
}

// ChangeOwnerRequest defines the new owner of a group
type ChangeOwnerRequest struct {
	// Required
	GroupID string `json:"group_id"`
	// Required. UserId of the new owner of the group
	// who must be an active member of the group
	OwnerID string `json:"owner_id"`
}

func (r ChangeOwnerRequest) String() string {
	return marshal(&r)
}

// ChangeOwnerResult holds the status of the group owner change
type ChangeOwnerResult struct {
	GroupID string `json:"group_id"`
	// UserId of the new owner of the group who is
	// an active member of the group
	OwnerID string                `json:"owner_id"`
	Status  changeOwnerStatusCode `json:"status"`
}

func (r ChangeOwnerResult) String() string {
	return marshal(&r)
}
