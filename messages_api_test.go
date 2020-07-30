package groupme

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type MessagesAPISuite struct{ APISuite }

func (s *MessagesAPISuite) SetupSuite() {
	s.handler = messagesTestRouter()
	s.setupSuite()
}

func (s *MessagesAPISuite) TestMessagesIndex() {
	resp, err := s.client.IndexMessages(
		ID("123"),
		&IndexMessagesQuery{
			BeforeID: "0123456789",
			SinceID:  "9876543210",
			AfterID:  "0246813579",
			Limit:    20,
		},
	)
	s.Require().NoError(err)
	s.Require().NotZero(resp)
	for _, message := range resp.Messages {
		s.Assert().NotZero(message)
	}
}

func (s *MessagesAPISuite) TestMessagesCreate() {
	message, err := s.client.CreateMessage(
		ID("123"),
		&Message{
			Text: "Test",
		},
	)
	s.Require().NoError(err)
	s.Require().NotNil(message)
	s.Assert().NotZero(*message)
}

func TestMessagesAPISuite(t *testing.T) {
	suite.Run(t, new(MessagesAPISuite))
}

func messagesTestRouter() *mux.Router {
	router := mux.NewRouter()

	// Index
	router.Path("/groups/{id:[0-9]+}/messages").
		Methods("GET").
		Name("IndexMessages").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": {
					"count": 123,
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

	// Create
	router.Path("/groups/{id:[0-9]+}/messages").
		Methods("POST").
		Name("CreateMessages").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(201)
			fmt.Fprint(w, `{
				"response": {
					"message": {
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
