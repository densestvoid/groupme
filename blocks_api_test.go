package groupme

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type BlocksAPISuite struct{ APISuite }

func (s *BlocksAPISuite) SetupSuite() {
	s.handler = blocksTestRouter()
	s.setupSuite()
}

func (s *BlocksAPISuite) TestBlocksIndex() {
	blocks, err := s.client.IndexBlock("1")
	s.Require().NoError(err)
	s.Require().NotZero(blocks)
	for _, block := range blocks {
		s.Assert().NotZero(block)
	}
}

func (s *BlocksAPISuite) TestBlocksBetween() {
	between, err := s.client.BlockBetween("1", "2")
	s.Require().NoError(err)
	s.Assert().True(between)
}

func (s *BlocksAPISuite) TestBlocksCreate() {
	block, err := s.client.CreateBlock("1", "2")
	s.Require().NoError(err)
	s.Assert().NotZero(block)
}

func (s *BlocksAPISuite) TestBlocksUnblock() {
	s.Assert().NoError(s.client.Unblock("1", "2"))
}

func TestBlocksAPISuite(t *testing.T) {
	suite.Run(t, new(BlocksAPISuite))
}

func blocksTestRouter() *mux.Router {
	router := mux.NewRouter()

	// Index
	router.Path("/blocks").
		Queries("user", "").
		Methods("GET").
		Name("IndexBlock").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": {
					"blocks": [
						{
							"user_id": "1234567890",
							"blocked_user_id": "1234567890",
							"created_at": 1302623328
						}
					]
				},
				"meta": {
					"code": 200,
					"errors": []
				}
			}`)
		})

	// Block Between
	router.Path("/blocks/between").
		Queries("user", "", "otherUser", "").
		Methods("GET").
		Name("BlockBetween").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": {
					"between": true
				},
				"meta": {
					"code": 200,
					"errors": []
				}
			}`)
		})

	// Create Block
	router.Path("/blocks").
		Queries("user", "", "otherUser", "").
		Methods("POST").
		Name("CreateBlock").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
				"response": {
					"block": {
						"user_id": "1234567890",
						"blocked_user_id": "1234567890",
						"created_at": 1302623328
					}
				},
				"meta": {
					"code": 200,
					"errors": []
				}
			}`)
		})

	// Unblock
	router.Path("/blocks").
		Queries("user", "", "otherUser", "").
		Methods("DELETE").
		Name("Unblock").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
		})

	/*// Return test router //*/
	return router
}
