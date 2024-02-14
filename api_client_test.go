package groupme

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ClientSuite struct{ APISuite }

func (s *ClientSuite) SetupSuite() {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("error"))
		s.Require().NoError(err)
	})
	s.handler = serverMux
	s.setupSuite()
}

func (s *ClientSuite) SetupTest() {
	s.client = NewClient("")
	s.Require().NotNil(s.client)

	s.client.apiEndpointBase = s.addr
}

func (s *ClientSuite) TestClient_Close() {
	s.Assert().NoError(s.client.Close())
}

func (s *ClientSuite) TestClient_do_PostContentType() {
	req, err := http.NewRequest("POST", "", nil)
	s.Require().NoError(err)

	s.Assert().Error(s.client.do(context.Background(), req, struct{}{}))
	s.Assert().EqualValues(req.Header.Get("Content-Type"), "application/json")
}

func (s *ClientSuite) TestClient_do_DoError() {
	req, err := http.NewRequest("", "", nil)
	s.Require().NoError(err)

	s.Assert().Error(s.client.do(context.Background(), req, struct{}{}))
}

func (s *ClientSuite) TestClient_do_UnmarshalError() {
	req, err := http.NewRequest("GET", s.addr, nil)
	s.Require().NoError(err)

	s.Assert().Error(s.client.do(context.Background(), req, struct{}{}))
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}
