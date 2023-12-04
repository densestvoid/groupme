package groupme

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type LeaderboardAPISuite struct{ APISuite }

func (s *LeaderboardAPISuite) SetupSuite() {
	s.handler = leaderboardTestRouter()
	s.setupSuite()
}

func (s *LeaderboardAPISuite) TestLeaderboardIndex() {
	messages, err := s.client.IndexLeaderboard(context.Background(), "1", PeriodDay)
	s.Require().NoError(err)
	s.Require().NotZero(messages)
	for _, message := range messages {
		s.Assert().NotZero(message)
	}
}

func (s *LeaderboardAPISuite) TestLeaderboardMyLikes() {
	messages, err := s.client.MyLikesLeaderboard(context.Background(), "1")
	s.Require().NoError(err)
	s.Require().NotZero(messages)
	for _, message := range messages {
		s.Assert().NotZero(message)
	}
}

func (s *LeaderboardAPISuite) TestLeaderboardMyHits() {
	messages, err := s.client.MyHitsLeaderboard(context.Background(), "1")
	s.Require().NoError(err)
	s.Require().NotZero(messages)
	for _, message := range messages {
		s.Assert().NotZero(message)
	}
}

func TestLeaderboardAPISuite(t *testing.T) {
	suite.Run(t, new(LeaderboardAPISuite))
}

// nolint // not duplicate code
func leaderboardTestRouter() *mux.Router {
	router := mux.NewRouter().Queries("token", "").Subrouter()

	// Index
	router.Path("/groups/{id:[0-9]+}/likes").
		Queries("period", "{period:day|week|month}").
		Methods("GET").
		Name("IndexLeaderboard").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": {
					"messages": [
						{
							"id": "1234567890",
							"source_guid": "GUID",
							"created_at": 1302623328,
							"user_id": "1234567890",
							"group_id": "1234567890",
							"name": "John",
							"avatar_url": "https://i.groupme.com/123456789",
							"text": "Hello world ☃☃",
							"system": true,
							"favorited_by": [
								"101",
								"66",
								"1234567890"
							],
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
						},
						{
							"id": "1234567890",
							"source_guid": "GUID",
							"created_at": 1302623328,
							"user_id": "1234567890",
							"group_id": "1234567890",
							"name": "John",
							"avatar_url": "https://i.groupme.com/123456789",
							"text": "Hello world ☃☃",
							"system": true,
							"favorited_by": [
								"1",
								"2"
							],
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
					]
				},
				"meta": {
					"code": 200,
					"errors": []
				}
			}`)
		})

	// My Likes
	router.Path("/groups/{id:[0-9]+}/likes/mine").
		Methods("GET").
		Name("MyLikesLeaderboard").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": {
					"messages": [
						{
							"id": "1234567890",
							"source_guid": "GUID",
							"created_at": 1302623328,
							"user_id": "1234567890",
							"group_id": "1234567890",
							"name": "John",
							"avatar_url": "https://i.groupme.com/123456789",
							"text": "Hello world ☃☃",
							"system": true,
							"favorited_by": [
								"101",
								"66",
								"1234567890"
							],
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
							],
							"liked_at": "2014-05-08T18:30:31.6617Z"
						}
					]
				},
				"meta": {
					"code": 200,
					"errors": []
				}
			}`)
		})

	// My Hits
	router.Path("/groups/{id:[0-9]+}/likes/for_me").
		Methods("GET").
		Name("MyHitsLeaderboard").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": {
					"messages": [
						{
							"id": "1234567890",
							"source_guid": "GUID",
							"created_at": 1302623328,
							"user_id": "1234567890",
							"group_id": "1234567890",
							"name": "John",
							"avatar_url": "https://i.groupme.com/123456789",
							"text": "Hello world ☃☃",
							"system": true,
							"favorited_by": [
								"101",
								"66",
								"1234567890"
							],
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
					]
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
