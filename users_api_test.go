package groupme

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type UsersAPISuite struct{ APISuite }

func (s *UsersAPISuite) SetupSuite() {
	s.handler = usersTestRouter()
	s.setupSuite()
}

func (s *UsersAPISuite) TestUsersMe() {
	user, err := s.client.MyUser(context.Background())
	s.Require().NoError(err)
	s.Assert().NotZero(user)
}

func (s *UsersAPISuite) TestUsersUpdate() {
	user, err := s.client.UpdateMyUser(context.Background(), UserSettings{})
	s.Require().NoError(err)
	s.Assert().NotZero(user)
}

func TestUsersAPISuite(t *testing.T) {
	suite.Run(t, new(UsersAPISuite))
}

// nolint // not duplicate code
func usersTestRouter() *mux.Router {
	router := mux.NewRouter().Queries("token", "").Subrouter()

	// Me
	router.Path("/users/me").
		Methods("GET").
		Name("MyUser").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": {
					"id": "1234567890",
					"phone_number": "+1 2123001234",
					"image_url": "https://i.groupme.com/123456789",
					"name": "Ronald Swanson",
					"created_at": 1302623328,
					"updated_at": 1302623328,
					"email": "me@example.com",
					"sms": false
				},
				"meta": {
					"code": 200,
					"errors": []
				}
			}`)
		})

	// Update
	router.Path("/users/update").
		Methods("POST").
		Name("UpdateMyUser").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": {
					"id": "1234567890",
					"phone_number": "+1 2123001234",
					"image_url": "https://i.groupme.com/123456789",
					"name": "Ronald Swanson",
					"created_at": 1302623328,
					"updated_at": 1302623328,
					"email": "me@example.com",
					"sms": false
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
