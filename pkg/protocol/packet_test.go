package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPacketBodyNotValid(t *testing.T) {
	packet := packet{}
	err := packet.validate()
	assert.NoError(t, err)
	assert.Equal(t, int32(10), packet.size())
}

func TestPacketSizeGreaterThanAllowed(t *testing.T) {
	packet := packet{}
	b := make([]byte, 4097)
	for i := range b {
		b[i] = 'a'
	}
	packet.Body = string(b)
	err := packet.validate()
	assert.Error(t, err)
}

func TestDefaultPacketSize(t *testing.T) {
	packet := packet{}
	raw, err := packet.serialize()
	assert.Nil(t, err)
	assert.NotNil(t, raw)
	assert.Equal(t, 14, len(raw)) // 10 packet + 1 size header + 1 empty string termination
	assert.Equal(t, int32(10), packet.size())
}

func TestSerializePacket(t *testing.T) {
	packet := newPacket(serverDataAuth, "aaaa")

	raw, err := packet.serialize()
	assert.NoError(t, packet.validate())
	assert.Nil(t, err)
	assert.NotNil(t, raw)
	assert.Equal(t, int32(14), packet.size())
	assert.Equal(t, 18, len(raw))
}

func TestDeserializePacket(t *testing.T) {
	p1 := newPacket(serverDataAuth, "aaaa")

	raw, err := p1.serialize()
	assert.NoError(t, p1.validate())
	assert.Nil(t, err)
	assert.NotNil(t, raw)

	p2 := packet{}
	err = p2.deserialize(raw)
	assert.NoError(t, err)
	assert.Equal(t, "aaaa", p2.Body)
	assert.Equal(t, int32(14), p2.size())
	assert.Equal(t, p1.ID, p2.ID)
	assert.Equal(t, p1.Type, p2.Type)
}
