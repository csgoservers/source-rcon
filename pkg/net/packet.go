package net

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
)

const (
	minimumPacketSize int32 = 10
	maximumPacketSize int32 = 4096 + minimumPacketSize

	serverDataResponseValue PacketType = 0
	serverDataExecCommand   PacketType = 2
	serverDataAuthResponse  PacketType = 2
	serverDataAuth          PacketType = 3
)

// PacketType indicates the purpose of the packet.
type PacketType int32

// Packet is the payload that both requests and responses
// are sent as TCP packets.
type Packet struct {
	ID   int32
	Type PacketType
	Body string
}

// NewPacket creates a default packet
func NewPacket(t PacketType, body string) *Packet {
	return &Packet{
		ID:   rand.Int31(),
		Type: t,
		Body: body,
	}
}

// Validate validates that the packet size is the expected one.
// The minimum posible value for packet `Size` is 10 bytes
// and the maximum one is 4096.
func (p *Packet) Validate() error {
	if p.Size() > maximumPacketSize {
		return fmt.Errorf("packet size is > %d", maximumPacketSize)
	}
	return nil
}

// String return the packet representation
func (p *Packet) String() string {
	return fmt.Sprintf("%02x %02x %02x %v", p.Size(), p.ID, p.Type, []byte(p.Body))
}

// Serialize transforms the packet to a byte array to be sent to the
// server.
func (p *Packet) Serialize() ([]byte, error) {
	raw := new(bytes.Buffer)
	binary.Write(raw, binary.LittleEndian, p.Size())
	binary.Write(raw, binary.LittleEndian, p.ID)
	binary.Write(raw, binary.LittleEndian, p.Type)

	raw.WriteString(p.Body)
	raw.WriteByte(0) // body payload must be null terminated.
	raw.WriteByte(0)

	return raw.Bytes(), nil
}

func (p *Packet) bodySize() int32 {
	return int32(len(p.Body))
}

// Size computes the total package size
func (p *Packet) Size() int32 {
	return (p.bodySize() + minimumPacketSize)
}

// Deserialize transforms the payload into a readable Packet.
func (p *Packet) Deserialize(payload []byte) error {
	data := bytes.NewBuffer(payload)

	var size int32
	binary.Read(data, binary.LittleEndian, &size)
	binary.Read(data, binary.LittleEndian, &p.ID)
	binary.Read(data, binary.LittleEndian, &p.Type)
	// no data to read, return ok
	if size == 0 {
		p.Body = ""
		return nil
	}
	// read body data
	bodyData := make([]byte, size-minimumPacketSize)
	_, err := io.ReadFull(data, bodyData)
	if err == io.EOF {
		p.Body = ""
		return nil
	}
	if err != nil {
		return err
	}
	p.Body = string(bodyData)
	return nil
}
