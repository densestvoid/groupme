package groupme

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type ChatsAPISuite struct{ APISuite }

func (s *ChatsAPISuite) SetupSuite() {
	s.handler = chatsTestRouter()
	s.setupSuite()
}

func (s *ChatsAPISuite) TestChatsIndex() {
	chats, err := s.client.IndexChats(
		context.Background(),
		&IndexChatsQuery{
			Page:    1,
			PerPage: 20,
		},
	)
	s.Require().NoError(err)
	s.Require().NotZero(chats)
	for _, chat := range chats {
		s.Assert().NotZero(chat)
	}
}

func TestChatsAPISuite(t *testing.T) {
	suite.Run(t, new(ChatsAPISuite))
}

func chatsTestRouter() *mux.Router {
	router := mux.NewRouter().Queries("token", "").Subrouter()

	// Index
	router.Path("/chats").
		Methods("GET").
		Name("IndexChats").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": [
					{
						"created_at": 1352299338,
						"updated_at": 1352299338,
						"last_message": {
							"attachments": [
					
							],
							"avatar_url": "https://i.groupme.com/200x200.jpeg.abcdef",
							"conversation_id": "12345+67890",
							"created_at": 1352299338,
							"favorited_by": [
					
							],
							"id": "1234567890",
							"name": "John Doe",
							"recipient_id": "67890",
							"sender_id": "12345",
							"sender_type": "user",
							"source_guid": "GUID",
							"text": "Hello world",
							"user_id": "12345"
						},
						"messages_count": 10,
						"other_user": {
							"avatar_url": "https://i.groupme.com/200x200.jpeg.abcdef",
							"id": "12345",
							"name": "John Doe"
						}
					}
				],
				"meta": {
					"code": 200,
					"errors": []
				}
			}`)
		})

	/*// Return test router //*/
	return router
}
