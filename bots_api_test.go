package groupme

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type BotsAPISuite struct{ APISuite }

func (s *BotsAPISuite) SetupSuite() {
	s.handler = botsTestRouter()
	s.setupSuite()
}

func (s *BotsAPISuite) TestBotsCreate() {
	bot, err := s.client.CreateBot(&Bot{
		Name:           "test",
		GroupID:        "1",
		AvatarURL:      "url.com",
		CallbackURL:    "otherURL.com",
		DMNotification: true,
	})
	s.Require().NoError(err)
	s.Require().NotZero(bot)
}

func (s *BotsAPISuite) TestBotsPostMessage() {
	err := s.client.PostBotMessage("1", "test message", nil)
	s.Require().NoError(err)
}

func (s *BotsAPISuite) TestBotsIndex() {
	bots, err := s.client.IndexBots()
	s.Require().NoError(err)
	s.Require().NotZero(bots)
	for _, bot := range bots {
		s.Assert().NotZero(bot)
	}
}

func (s *BotsAPISuite) TestBotsDestroy() {
	s.Require().NoError(s.client.DestroyBot("1"))
}

func TestBotsAPISuite(t *testing.T) {
	suite.Run(t, new(BotsAPISuite))
}

func botsTestRouter() *mux.Router {
	router := mux.NewRouter()

	// Create
	router.Path("/bots").
		Methods("POST").
		Name("CreateBot").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(201)
			fmt.Fprint(w, `{
				"response": {
					"bot_id": "1234567890",
					"group_id": "1234567890",
					"name": "hal9000",
					"avatar_url": "https://i.groupme.com/123456789",
					"callback_url": "https://example.com/bots/callback",
					"dm_notification": false
				},
				"meta": {
					"code": 201,
					"errors": []
				}
			}`)
		})

	// Post Message
	router.Path("/bots/post").
		Methods("POST").
		Name("PostBotMessage").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(201)
		})

	// Index
	router.Path("/bots").
		Methods("GET").
		Name("IndexBots").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": [
					{
						"bot_id": "1234567890",
						"group_id": "1234567890",
						"name": "hal9000",
						"avatar_url": "https://i.groupme.com/123456789",
						"callback_url": "https://example.com/bots/callback",
						"dm_notification": false
					}
				],
				"meta": {
					"code": 200,
					"errors": []
				}
			}`)
		})

	// Destroy
	router.Path("/bots/destroy").
		Methods("POST").
		Name("DestroyBot").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(201)
		})

	/*// Return test router //*/
	return router
}
