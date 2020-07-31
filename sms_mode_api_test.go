package groupme

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type SMSModeAPISuite struct{ APISuite }

func (s *SMSModeAPISuite) SetupSuite() {
	s.handler = smsModeTestRouter()
	s.setupSuite()
}

func (s *SMSModeAPISuite) TestSMSModeCreate() {
	s.Assert().NoError(s.client.CreateSMSMode(10, nil))
}

func (s *SMSModeAPISuite) TestSMSModeDelete() {
	s.Assert().NoError(s.client.DeleteSMSMode())
}
func TestSMSModeAPISuite(t *testing.T) {
	suite.Run(t, new(SMSModeAPISuite))
}

func smsModeTestRouter() *mux.Router {
	router := mux.NewRouter().Queries("token", "").Subrouter()

	// Create
	router.Path("/users/sms_mode").
		Methods("POST").
		Name("CreateSMSMode").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(201)
		})

	// Delete
	router.Path("/users/sms_mode/delete").
		Methods("POST").
		Name("DeleteSMSMode").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
		})

	/*// Return test router //*/
	return router
}
