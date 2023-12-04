package groupme

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type GroupsAPISuite struct{ APISuite }

func (s *GroupsAPISuite) SetupSuite() {
	s.handler = groupsTestRouter()
	s.setupSuite()
}

func (s *GroupsAPISuite) TestGroupsIndex() {
	groups, err := s.client.IndexGroups(
		context.Background(),
		&GroupsQuery{
			Page:    5,
			PerPage: 20,
			Omit:    "memberships",
		},
	)
	s.Require().NoError(err)
	s.Require().NotZero(groups)
	for _, group := range groups {
		s.Assert().NotZero(group)
	}
}

func (s *GroupsAPISuite) TestGroupsFormer() {
	groups, err := s.client.FormerGroups(context.Background())
	s.Require().NoError(err)
	s.Require().NotZero(groups)
	for _, group := range groups {
		s.Assert().NotZero(group)
	}
}

func (s *GroupsAPISuite) TestGroupsShow() {
	group, err := s.client.ShowGroup(context.Background(), "1")
	s.Require().NoError(err)
	s.Assert().NotZero(group)
}

func (s *GroupsAPISuite) TestGroupsCreate() {
	group, err := s.client.CreateGroup(
		context.Background(),
		GroupSettings{
			"Test",
			"This is a test group",
			"www.blank.com/image",
			false,
			true,
		},
	)
	s.Require().NoError(err)
	s.Assert().NotZero(group)
}

func (s *GroupsAPISuite) TestGroupsUpdate() {
	group, err := s.client.UpdateGroup(context.Background(), "1", GroupSettings{
		"Test",
		"This is a test group",
		"www.blank.com/image",
		true,
		true,
	})
	s.Require().NoError(err)
	s.Assert().NotZero(group)
}

func (s *GroupsAPISuite) TestGroupsDestroy() {
	err := s.client.DestroyGroup(context.Background(), "1")
	s.Require().NoError(err)
}

func (s *GroupsAPISuite) TestGroupsJoin() {
	group, err := s.client.JoinGroup(context.Background(), "1", "please")
	s.Require().NoError(err)
	s.Assert().NotZero(group)
}

func (s *GroupsAPISuite) TestGroupsRejoin() {
	group, err := s.client.RejoinGroup(context.Background(), "1")
	s.Require().NoError(err)
	s.Assert().NotZero(group)
}

func (s *GroupsAPISuite) TestGroupsChangeOwner() {
	result, err := s.client.ChangeGroupOwner(
		context.Background(),
		ChangeOwnerRequest{
			"1",
			"123",
		},
	)
	s.Require().NoError(err)
	s.Assert().NotZero(result)
}
func TestGroupsAPISuite(t *testing.T) {
	suite.Run(t, new(GroupsAPISuite))
}

/*//////// Test Groups Router ////////*/

// nolint // not duplicate code
func groupsTestRouter() *mux.Router {
	router := mux.NewRouter().Queries("token", "").Subrouter()

	// Index
	router.Path("/groups").
		Methods("GET").
		Name("IndexGroups").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": [
					{
						"id": "1234567890",
						"name": "Family",
						"type": "private",
						"description": "Coolest Family Ever",
						"image_url": "https://i.groupme.com/123456789",
						"creator_user_id": "1234567890",
						"created_at": 1302623328,
						"updated_at": 1302623328,
						"members": [
							{
								"user_id": "1234567890",
								"nickname": "Jane",
								"muted": false,
								"image_url": "https://i.groupme.com/123456789"
							}
						],
						"share_url": "https://groupme.com/join_group/1234567890/SHARE_TOKEN",
						"messages": {
							"count": 100,
							"last_message_id": "1234567890",
							"last_message_created_at": 1302623328,
							"preview": {
								"nickname": "Jane",
								"text": "Hello world",
								"image_url": "https://i.groupme.com/123456789",
								"attachments": [
									{
										"type": "image",
										"url": "https://i.groupme.com/123456789"
									},
									{
										"type": "image",
										"url": "https://i.groupme.com/123456789"
									},
									{
										"type": "location",
										"lat": "40.738206",
										"lng": "-73.993285",
										"name": "GroupMe HQ"
									},
									{
										"type": "split",
										"token": "SPLIT_TOKEN"
									},
									{
										"type": "emoji",
										"placeholder": "☃",
										"charmap": [
											[
												1,
												42
											],
											[
												2,
												34
											]
										]
									}
								]
							}
						}
					}
				],
				"meta": {
					"code": 201,
					"errors": []
				}
			}`)
		})

	// Former
	router.Path("/groups/former").
		Methods("GET").
		Name("FormerGroups").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": [
					{
						"id": "1234567890",
						"name": "Family",
						"type": "private",
						"description": "Coolest Family Ever",
						"image_url": "https://i.groupme.com/123456789",
						"creator_user_id": "1234567890",
						"created_at": 1302623328,
						"updated_at": 1302623328,
						"members": [
							{
								"user_id": "1234567890",
								"nickname": "Jane",
								"muted": false,
								"image_url": "https://i.groupme.com/123456789"
							}
						],
						"share_url": "https://groupme.com/join_group/1234567890/SHARE_TOKEN",
						"messages": {
							"count": 100,
							"last_message_id": "1234567890",
							"last_message_created_at": 1302623328,
							"preview": {
								"nickname": "Jane",
								"text": "Hello world",
								"image_url": "https://i.groupme.com/123456789",
								"attachments": [
									{
										"type": "image",
										"url": "https://i.groupme.com/123456789"
									},
									{
										"type": "image",
										"url": "https://i.groupme.com/123456789"
									},
									{
										"type": "location",
										"lat": "40.738206",
										"lng": "-73.993285",
										"name": "GroupMe HQ"
									},
									{
										"type": "split",
										"token": "SPLIT_TOKEN"
									},
									{
										"type": "emoji",
										"placeholder": "☃",
										"charmap": [
											[
												1,
												42
											],
											[
												2,
												34
											]
										]
									}
								]
							}
						}
					}
				],
				"meta": {
					"code": 200,
					"errors": []
				}
			}`)
		})

	// Show
	router.Path("/groups/{id:[0-9]+}").
		Methods("GET").
		Name("ShowGroup").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": {
					"id": "1234567890",
					"name": "Family",
					"type": "private",
					"description": "Coolest Family Ever",
					"image_url": "https://i.groupme.com/123456789",
					"creator_user_id": "1234567890",
					"created_at": 1302623328,
					"updated_at": 1302623328,
					"members": [
						{
							"user_id": "1234567890",
							"nickname": "Jane",
							"muted": false,
							"image_url": "https://i.groupme.com/123456789"
						}
					],
					"share_url": "https://groupme.com/join_group/1234567890/SHARE_TOKEN",
					"messages": {
						"count": 100,
						"last_message_id": "1234567890",
						"last_message_created_at": 1302623328,
						"preview": {
							"nickname": "Jane",
							"text": "Hello world",
							"image_url": "https://i.groupme.com/123456789",
							"attachments": [
								{
									"type": "image",
									"url": "https://i.groupme.com/123456789"
								},
								{
									"type": "image",
									"url": "https://i.groupme.com/123456789"
								},
								{
									"type": "location",
									"lat": "40.738206",
									"lng": "-73.993285",
									"name": "GroupMe HQ"
								},
								{
									"type": "split",
									"token": "SPLIT_TOKEN"
								},
								{
									"type": "emoji",
									"placeholder": "☃",
									"charmap": [
										[
											1,
											42
										],
										[
											2,
											34
										]
									]
								}
							]
						}
					}
				},
				"meta": {
					"code": 200,
					"errors": []
				}
			}`)
		})

	// Create
	router.Path("/groups").
		Methods("POST").
		Name("CreateGroup").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(201)
			fmt.Fprint(w, `{
				"response": {
					"id": "1234567890",
					"name": "Family",
					"type": "private",
					"description": "Coolest Family Ever",
					"image_url": "https://i.groupme.com/123456789",
					"creator_user_id": "1234567890",
					"created_at": 1302623328,
					"updated_at": 1302623328,
					"members": [
						{
							"user_id": "1234567890",
							"nickname": "Jane",
							"muted": false,
							"image_url": "https://i.groupme.com/123456789"
						}
					],
					"share_url": "https://groupme.com/join_group/1234567890/SHARE_TOKEN",
					"messages": {
						"count": 100,
						"last_message_id": "1234567890",
						"last_message_created_at": 1302623328,
						"preview": {
							"nickname": "Jane",
							"text": "Hello world",
							"image_url": "https://i.groupme.com/123456789",
							"attachments": [
								{
									"type": "image",
									"url": "https://i.groupme.com/123456789"
								},
								{
									"type": "image",
									"url": "https://i.groupme.com/123456789"
								},
								{
									"type": "location",
									"lat": "40.738206",
									"lng": "-73.993285",
									"name": "GroupMe HQ"
								},
								{
									"type": "split",
									"token": "SPLIT_TOKEN"
								},
								{
									"type": "emoji",
									"placeholder": "☃",
									"charmap": [
										[
											1,
											42
										],
										[
											2,
											34
										]
									]
								}
							]
						}
					}
				},
				"meta": {
					"code": 201,
					"errors": []
				}
			}`)
		})

	// Update
	router.Path("/groups/{id:[0-9]+}/update").
		Methods("POST").
		Name("UpdateGroup").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": {
					"id": "1234567890",
					"name": "Family",
					"type": "private",
					"description": "Coolest Family Ever",
					"image_url": "https://i.groupme.com/123456789",
					"creator_user_id": "1234567890",
					"created_at": 1302623328,
					"updated_at": 1302623328,
					"members": [
						{
							"user_id": "1234567890",
							"nickname": "Jane",
							"muted": false,
							"image_url": "https://i.groupme.com/123456789"
						}
					],
					"share_url": "https://groupme.com/join_group/1234567890/SHARE_TOKEN",
					"messages": {
						"count": 100,
						"last_message_id": "1234567890",
						"last_message_created_at": 1302623328,
						"preview": {
							"nickname": "Jane",
							"text": "Hello world",
							"image_url": "https://i.groupme.com/123456789",
							"attachments": [
								{
									"type": "image",
									"url": "https://i.groupme.com/123456789"
								},
								{
									"type": "image",
									"url": "https://i.groupme.com/123456789"
								},
								{
									"type": "location",
									"lat": "40.738206",
									"lng": "-73.993285",
									"name": "GroupMe HQ"
								},
								{
									"type": "split",
									"token": "SPLIT_TOKEN"
								},
								{
									"type": "emoji",
									"placeholder": "☃",
									"charmap": [
										[
											1,
											42
										],
										[
											2,
											34
										]
									]
								}
							]
						}
					}
				},
				"meta": {
					"code": 201,
					"errors": []
				}
			}`)
		})

	// Destroy
	router.Path("/groups/{id:[0-9]+}/destroy").
		Methods("POST").
		Name("DestroyGroup").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
		})

	// Join
	router.Path("/groups/{id:[0-9]+}/join/{share_token}").
		Methods("POST").
		Name("JoinGroup").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": {
					"group": {
						"id": "1234567890",
						"name": "Family",
						"type": "private",
						"description": "Coolest Family Ever",
						"image_url": "https://i.groupme.com/123456789",
						"creator_user_id": "1234567890",
						"created_at": 1302623328,
						"updated_at": 1302623328,
						"members": [
							{
								"user_id": "1234567890",
								"nickname": "Jane",
								"muted": false,
								"image_url": "https://i.groupme.com/123456789"
							}
						],
						"share_url": "https://groupme.com/join_group/1234567890/SHARE_TOKEN",
						"messages": {
							"count": 100,
							"last_message_id": "1234567890",
							"last_message_created_at": 1302623328,
							"preview": {
								"nickname": "Jane",
								"text": "Hello world",
								"image_url": "https://i.groupme.com/123456789",
								"attachments": [
									{
										"type": "image",
										"url": "https://i.groupme.com/123456789"
									},
									{
										"type": "image",
										"url": "https://i.groupme.com/123456789"
									},
									{
										"type": "location",
										"lat": "40.738206",
										"lng": "-73.993285",
										"name": "GroupMe HQ"
									},
									{
										"type": "split",
										"token": "SPLIT_TOKEN"
									},
									{
										"type": "emoji",
										"placeholder": "☃",
										"charmap": [
											[
												1,
												42
											],
											[
												2,
												34
											]
										]
									}
								]
							}
						}
					}
				},
				"meta": {
					"code": 201,
					"errors": []
				}
			}`)
		})

	// Rejoin
	router.Path("/groups/join").
		Methods("POST").
		Name("RejoinGroup").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": {
					"id": "1234567890",
					"name": "Family",
					"type": "private",
					"description": "Coolest Family Ever",
					"image_url": "https://i.groupme.com/123456789",
					"creator_user_id": "1234567890",
					"created_at": 1302623328,
					"updated_at": 1302623328,
					"members": [
						{
							"user_id": "1234567890",
							"nickname": "Jane",
							"muted": false,
							"image_url": "https://i.groupme.com/123456789"
						}
					],
					"share_url": "https://groupme.com/join_group/1234567890/SHARE_TOKEN",
					"messages": {
						"count": 100,
						"last_message_id": "1234567890",
						"last_message_created_at": 1302623328,
						"preview": {
							"nickname": "Jane",
							"text": "Hello world",
							"image_url": "https://i.groupme.com/123456789",
							"attachments": [
								{
									"type": "image",
									"url": "https://i.groupme.com/123456789"
								},
								{
									"type": "image",
									"url": "https://i.groupme.com/123456789"
								},
								{
									"type": "location",
									"lat": "40.738206",
									"lng": "-73.993285",
									"name": "GroupMe HQ"
								},
								{
									"type": "split",
									"token": "SPLIT_TOKEN"
								},
								{
									"type": "emoji",
									"placeholder": "☃",
									"charmap": [
										[
											1,
											42
										],
										[
											2,
											34
										]
									]
								}
							]
						}
					}
				},
				"meta": {
					"code": 201,
					"errors": []
				}
			}`)
		})

	// Change Owner
	router.Path("/groups/change_owners").
		Methods("POST").
		Name("ChangeGroupOwner").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": {
					"results": [
						{
							"group_id": "1234567890",
							"owner_id": "1234567890",
							"status": "200"
						},
						{
							"group_id": "1234567890",
							"owner_id": "1234567890",
							"status": "400"
						}
					]
				},
				"meta": {
					"code": 201,
					"errors": []
				}
			}`)
		})

	/*// Return test router //*/
	return router
}
