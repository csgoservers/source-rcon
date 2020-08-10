package protocol

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsUrl(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	assert.NoError(t, err)
	defer listener.Close()

	opts := &Options{Host: "127.0.0.1", Port: 8080}
	conn := New(opts)
	defer conn.Close()
	err = conn.initialize()

	assert.NoError(t, err, "Found unexpected error on initialize")
	assert.Equal(t, "127.0.0.1:8080", conn.connection.RemoteAddr().String())
}
