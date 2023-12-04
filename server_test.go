package groupme

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	var received bool
	server := httptest.NewServer(HTTPHandlerFunc(func(m Message) { received = true }))
	defer server.Close()

	client := server.Client()

	// Wait until server has started listening
	// nolint // url is meant to be variable
	for _, err := client.Get(server.URL); err != nil; _, err = client.Get(server.URL) {
		continue
	}

	t.Run("Valid", func(t *testing.T) {
		received = false
		msgBytes, err := json.Marshal(Message{})
		require.NoError(t, err)

		resp, err := client.Post(server.URL, "application/json", bytes.NewReader(msgBytes))
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, received)
	})

	t.Run("UnsupportedMediaType", func(t *testing.T) {
		received = false
		msgBytes, err := json.Marshal(Message{})
		require.NoError(t, err)

		resp, err := client.Post(server.URL, "application/text", bytes.NewReader(msgBytes))
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
		assert.False(t, received)
	})

	t.Run("BadRequest", func(t *testing.T) {
		received = false
		msgBytes, err := json.Marshal(struct{ Bad string }{Bad: "bad"})
		require.NoError(t, err)

		resp, err := client.Post(server.URL, "application/json", bytes.NewReader(msgBytes))
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.False(t, received)
	})
}
