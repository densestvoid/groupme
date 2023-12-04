package groupme

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	addr := "localhost:" + GeneratePort()

	var received bool
	mux := http.NewServeMux()
	mux.HandleFunc("/", HTTPHandlerFunc(func(m Message) { received = true }))
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go server.ListenAndServe()
	defer server.Shutdown(context.Background())

	// Wait until server has started listening
	// nolint // url is meant to be variable
	for _, err := http.Get("http://" + addr); err != nil; _, err = http.Get("http://" + addr) {
		continue
	}

	t.Run("Valid", func(t *testing.T) {
		received = false
		msgBytes, err := json.Marshal(Message{})
		require.NoError(t, err)

		resp, err := http.DefaultClient.Post("http://"+addr, "application/json", bytes.NewReader(msgBytes))
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, received)
	})

	t.Run("UnsupportedMediaType", func(t *testing.T) {
		received = false
		msgBytes, err := json.Marshal(Message{})
		require.NoError(t, err)

		resp, err := http.DefaultClient.Post("http://"+addr, "application/text", bytes.NewReader(msgBytes))
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
		assert.False(t, received)
	})

	t.Run("BadRequest", func(t *testing.T) {
		received = false
		msgBytes, err := json.Marshal(struct{ Bad string }{Bad: "bad"})
		require.NoError(t, err)

		resp, err := http.DefaultClient.Post("http://"+addr, "application/json", bytes.NewReader(msgBytes))
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.False(t, received)
	})
}
