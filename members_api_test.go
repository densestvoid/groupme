// Package groupme defines a client capable of executing API commands for the GroupMe chat service
package groupme

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type MembersAPISuite struct{ APISuite }

func (s *MembersAPISuite) SetupSuite() {
	s.handler = membersTestRouter()
	s.setupSuite()
}

func (s *MembersAPISuite) TestMembersAdd() {
	_, err := s.client.AddMembers(
		context.Background(),
		"1",
		&Member{Nickname: "test"},
	)
	s.Require().NoError(err)
}

func (s *MembersAPISuite) TestMembersResults() {
	_, err := s.client.AddMembersResults(context.Background(), "1", "123")
	s.Require().NoError(err)
}

func (s *MembersAPISuite) TestMembersRemove() {
	err := s.client.RemoveMember(context.Background(), "1", "123")
	s.Require().NoError(err)
}

func (s *MembersAPISuite) TestMembersUpdate() {
	_, err := s.client.UpdateMember(context.Background(), "1", "nickname")
	s.Require().NoError(err)
}

func TestMembersAPISuite(t *testing.T) {
	suite.Run(t, new(MembersAPISuite))
}

func membersTestRouter() *mux.Router {
	router := mux.NewRouter().Queries("token", "").Subrouter()

	// Add
	router.Path("/groups/{id:[0-9]+}/members/add").
		Methods("POST").
		Name("AddMembers").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(202)
			fmt.Fprint(w, `{
				"response": {
					"results_id": "GUID"
				},
				"meta": {
					"code": 202,
					"errors": []
				}
			}`)
		})

	// Results
	router.Path("/groups/{id:[0-9]+}/members/results/{result_id}").
		Methods("GET").
		Name("AddMembersResults").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": {
					"members": [
						{
							"id": "1000",
							"user_id": "10000",
							"nickname": "John",
							"muted": false,
							"image_url": "https://i.groupme.com/AVATAR",
							"autokicked": false,
							"app_installed": true,
							"guid": "GUID-1"
						},
						{
							"id": "2000",
							"user_id": "20000",
							"nickname": "Anne",
							"muted": false,
							"image_url": "https://i.groupme.com/AVATAR",
							"autokicked": false,
							"app_installed": true,
							"guid": "GUID-2"
						}
					]
				},
				"meta": {
					"code": 200,
					"errors": []
				}
			}`)
		})

	// Remove
	router.Path("/groups/{id:[0-9]+}/members/{membership_id}/remove").
		Methods("POST").
		Name("RemoveMember").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
		})

	// Update
	router.Path("/groups/{id:[0-9]+}/memberships/update").
		Methods("POST").
		Name("UpdateMember").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": {
					"id": "MEMBERSHIP ID",
					"user_id": "USER ID",
					"nickname": "NEW NICKNAME",
					"muted": false,
					"image_url": "AVATAR URL",
					"autokicked": false,
					"app_installed": true
				},
				"meta": {
					"code": 200,
					"errors": []
				}
			}`)
		})

	/*// Return test router //*/
	return router
}
