package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPacketBodyNotValid(t *testing.T) {
	packet := Packet{}
	err := packet.Validate()
	assert.NoError(t, err)
	assert.Equal(t, int32(10), packet.Size())
}

func TestPacketSizeGreaterThanAllowed(t *testing.T) {
	packet := Packet{}
	b := make([]byte, 4097)
	for i := range b {
		b[i] = 'a'
	}
	packet.Body = string(b)
	err := packet.Validate()
	assert.Error(t, err)
}

func TestDefaultPacketSize(t *testing.T) {
	packet := Packet{}
	raw, err := packet.Serialize()
	assert.Nil(t, err)
	assert.NotNil(t, raw)
	assert.Equal(t, 14, len(raw)) // 10 packet + 1 size header + 1 empty string termination
	assert.Equal(t, int32(10), packet.Size())
}

func TestSerializePacket(t *testing.T) {
	packet := NewPacket(serverDataAuth, "aaaa")

	raw, err := packet.Serialize()
	assert.NoError(t, packet.Validate())
	assert.Nil(t, err)
	assert.NotNil(t, raw)
	assert.Equal(t, int32(14), packet.Size())
	assert.Equal(t, 18, len(raw))
}

func TestDeserializePacket(t *testing.T) {
	p1 := NewPacket(serverDataAuth, "aaaa")

	raw, err := p1.Serialize()
	assert.NoError(t, p1.Validate())
	assert.Nil(t, err)
	assert.NotNil(t, raw)

	p2 := Packet{}
	err = p2.Deserialize(raw)
	assert.NoError(t, err)
	assert.Equal(t, "aaaa", p2.Body)
	assert.Equal(t, int32(14), p2.Size())
	assert.Equal(t, p1.ID, p2.ID)
	assert.Equal(t, p1.Type, p2.Type)
}
