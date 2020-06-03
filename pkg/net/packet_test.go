package net

import (
	"testing"
)

func TestPacketBodyNotValid(t *testing.T) {
	packet := Packet{}
	err := packet.Validate()
	if err == nil {
		t.Error(err)
	}
}

func TestPacketSizeGreaterThanAllowed(t *testing.T) {
	packet := Packet{}
	b := make([]byte, 4087)
	for i := range b {
		b[i] = 'a'
	}
	packet.Body = string(b)

	err := packet.Validate()
	if err == nil {
		t.Error(err)
	}
}
