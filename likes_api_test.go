package groupme

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type LikesAPISuite struct{ APISuite }

func (s *LikesAPISuite) SetupSuite() {
	s.handler = likesTestRouter()
	s.setupSuite()
}

func (s *LikesAPISuite) TestLikesCreate() {
	err := s.client.CreateLike("1", "1")
	s.Require().NoError(err)
}

func (s *LikesAPISuite) TestLikesDestroy() {
	err := s.client.DestroyLike("1", "1")
	s.Require().NoError(err)
}

func TestLikesAPISuite(t *testing.T) {
	suite.Run(t, new(LikesAPISuite))
}
func likesTestRouter() *mux.Router {
	router := mux.NewRouter().Queries("token", "").Subrouter()

	// Create
	router.Path(`/messages/{conversation_id}/{message_id}/like`).
		Methods("POST").
		Name("CreateLike").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
		})

	// Destroy
	router.Path(`/messages/{conversation_id}/{message_id}/unlike`).
		Methods("POST").
		Name("DestroyLike").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
		})

	/*// Return test router //*/
	return router
}
