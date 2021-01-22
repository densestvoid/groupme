// Package groupme defines a client capable of executing API commands for the GroupMe chat service
package groupme

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

/*//////// Base API Suite ////////*/
type APISuite struct {
	// Base attributes
	suite.Suite
	client *Client
	server *http.Server
	wg     sync.WaitGroup

	// Overridden by child Suite
	addr    string
	handler http.Handler
}

func (s *APISuite) setupSuite() {
	s.addr = "localhost:" + s.generatePort()

	s.client = NewClient("")
	s.client.endpointBase = "http://" + s.addr

	s.server = s.startServer(s.addr, s.handler)
}

func (s *APISuite) TearDownSuite() {
	s.client.Close()
	s.server.Close()
	s.wg.Wait()
}

/*/// Start Server ///*/
func (s *APISuite) startServer(addr string, handler http.Handler) *http.Server {
	server := &http.Server{
		Addr:     addr,
		Handler:  handler,
		ErrorLog: log.New(os.Stdout, "SERVER", log.Ltime),
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := server.ListenAndServe(); err.Error() != "http: Server closed" {
			s.Assert().NoError(err)
		}
	}()

	// Wait until server has started listening
	url := fmt.Sprintf("http://%s", addr)
	// nolint // url is meant to be variable
	for _, err := http.Get(url); err != nil; _, err = http.Get(url) {
		continue
	}

	return server
}

/*/// Generate Ephemeral Port ///*/
const (
	portMin   = 49152
	portMax   = 65535
	portRange = portMax - portMin
)

func (s *APISuite) generatePort() string {
	rand.Seed(time.Now().UnixNano())
	// nolint // weak random generator is ok for creating port number in a test
	return strconv.Itoa((rand.Intn(portRange) + portMin))
}

/*//////// Test Main ////////*/
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
