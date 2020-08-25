package groupme

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GroupMe documentation: https://dev.groupme.com/docs/v3#members

////////// Endpoints //////////
const (
	// Used to build other endpoints
	membersEndpointRoot = groupEndpointRoot + "/members"

	// Actual Endpoints
	addMembersEndpoint        = membersEndpointRoot + "/add"              // POST
	addMembersResultsEndpoint = membersEndpointRoot + "/results/%s"       // GET
	removeMemberEndpoint      = membersEndpointRoot + "/%s/remove"        // POST
	updateMemberEndpoint      = groupEndpointRoot + "/memberships/update" // POST
)

///// Add /////

/*
AddMembers -

Add members to a group.

Multiple members can be added in a single request, and results
are fetchedwith a separate call (since memberships are processed
asynchronously). The response includes a results_id that's used
in the results request.

In order to correlate request params with resulting memberships,
GUIDs can be added to the members parameters. These GUIDs will
be reflected in the membership JSON objects.

Parameters:
	groupID - required, ID(string)
	See Member.
		Nickname - required
		One of the following identifiers must be used:
			UserID - ID(string)
			PhoneNumber - PhoneNumber(string)
			Email - string
*/
func (c *Client) AddMembers(groupID ID, members ...*Member) (string, error) {
	URL := fmt.Sprintf(c.endpointBase+addMembersEndpoint, groupID)

	var data = struct {
		Members []*Member `json:"members"`
	}{
		members,
	}

	jsonBytes, err := json.Marshal(&data)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}

	var resp struct {
		ResultsID string `json:"results_id"`
	}

	err = c.doWithAuthToken(httpReq, &resp)
	if err != nil {
		return "", err
	}

	return resp.ResultsID, nil
}

///// Results /////

/*
AddMembersResults -
Get the membership results from an add call.

Successfully created memberships will be returned, including
any GUIDs that were sent up in the add request. If GUIDs were
absent, they are filled in automatically. Failed memberships
and invites are omitted.

Keep in mind that results are temporary -- they will only be
available for 1 hour after the add request.

Parameters:
	groupID - required, ID(string)
	resultID - required, string
*/
func (c *Client) AddMembersResults(groupID ID, resultID string) ([]*Member, error) {
	URL := fmt.Sprintf(c.endpointBase+addMembersResultsEndpoint, groupID, resultID)

	httpReq, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Members []*Member `json:"members"`
	}

	err = c.doWithAuthToken(httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Members, nil
}

///// Remove /////

/*
RemoveMember -

Remove a member (or yourself) from a group.

Note: The creator of the group cannot be removed or exit.

Parameters:
	groupID - required, ID(string)
	membershipID - required, ID(string). Not the same as userID
*/
func (c *Client) RemoveMember(groupID, membershipID ID) error {
	URL := fmt.Sprintf(c.endpointBase+removeMemberEndpoint, groupID, membershipID)

	httpReq, err := http.NewRequest("POST", URL, nil)
	if err != nil {
		return err
	}

	return c.doWithAuthToken(httpReq, nil)
}

///// Update /////

/*
UpdateMember -

Update your nickname in a group. The nickname must be
between 1 and 50 characters.
*/
func (c *Client) UpdateMember(groupID ID, nickname string) (*Member, error) {
	URL := fmt.Sprintf(c.endpointBase+updateMemberEndpoint, groupID)

	type Nickname struct {
		Nickname string `json:"nickname"`
	}
	var data = struct {
		Membership Nickname `json:"membership"`
	}{
		Nickname{nickname},
	}

	jsonBytes, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	var resp Member

	err = c.doWithAuthToken(httpReq, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
