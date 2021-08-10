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

type DirectMessagesAPISuite struct{ APISuite }

func (s *DirectMessagesAPISuite) SetupSuite() {
	s.handler = directMessagesTestRouter()
	s.setupSuite()
}

func (s *DirectMessagesAPISuite) TestDirectMessagesIndex() {
	resp, err := s.client.IndexDirectMessages(
		context.Background(),
		"123",
		&IndexDirectMessagesQuery{
			BeforeID: "0123456789",
			SinceID:  "9876543210",
		},
	)
	s.Require().NoError(err)
	s.Require().NotZero(resp)
	for _, message := range resp.Messages {
		s.Assert().NotZero(message)
	}
}

func (s *DirectMessagesAPISuite) TestDirectMessagesCreate() {
	message, err := s.client.CreateDirectMessage(
		context.Background(),
		&Message{
			RecipientID: ID("123"),
			Text:        "Test",
		},
	)
	s.Require().NoError(err)
	s.Require().NotNil(message)
	s.Assert().NotZero(*message)
}

func TestDirectMessagesAPISuite(t *testing.T) {
	suite.Run(t, new(DirectMessagesAPISuite))
}

// nolint // not duplicate code
func directMessagesTestRouter() *mux.Router {
	router := mux.NewRouter().Queries("token", "").Subrouter()

	// Index
	router.Path("/direct_messages").
		Methods("GET").
		Name("IndexDirectMessages").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": {
					"count": 123,
					"direct_messages": [
						{
							"id": "1234567890",
							"source_guid": "GUID",
							"recipient_id": "20",
							"user_id": "1234567890",
							"created_at": 1302623328,
							"name": "John",
							"avatar_url": "https://i.groupme.com/123456789",
							"text": "Hello world ☃☃",
							"favorited_by": [
								"101"
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

	// Create
	router.Path("/direct_messages").
		Methods("POST").
		Name("CreateDirectMessage").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(201)
			fmt.Fprint(w, `{
				"response": {
					"direct_message": {
						"id": "1234567890",
						"source_guid": "GUID",
						"recipient_id": "20",
						"user_id": "1234567890",
						"created_at": 1302623328,
						"name": "John",
						"avatar_url": "https://i.groupme.com/123456789",
						"text": "Hello world ☃☃",
						"favorited_by": [
							"101"
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
