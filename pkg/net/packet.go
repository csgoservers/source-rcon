package net

import (
	"errors"
	"fmt"
)

const (
	minimumPacketSize int32 = 10
	maximumPacketSize int32 = 4096
)

// PacketType indicates the purpose of the packet.
type PacketType int32

const (
	serverDataResponseValue PacketType = 0
)

// Packet is the payload that both requests and responses
// are sent as TCP packets.
type Packet struct {
	Size int32
	ID   int32
	Type PacketType
	Body string
}

// Validate validates that the packet size is the expected one.
// The minimum posible value for packet `Size` is 10 bytes
// and the maximum one is 4096.
func (p *Packet) Validate() error {
	if len(p.Body) == 0 {
		return errors.New("packet body is empty")
	}
	if int32(len(p.Body))+minimumPacketSize > maximumPacketSize {
		return fmt.Errorf("packet size is > %d", maximumPacketSize)
	}
	return nil
}
