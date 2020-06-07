package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsUrl(t *testing.T) {
	opts := Options{Host: "aa", Port: 1234}
	assert.Equal(t, "aa:1234", opts.url())
}
