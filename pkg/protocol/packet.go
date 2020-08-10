package protocol

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

	serverDataResponseValue packetType = 0
	serverDataExecCommand   packetType = 2
	serverDataAuthResponse  packetType = 2
	serverDataAuth          packetType = 3
)

// packetType indicates the purpose of the packet.
type packetType int32

// packet is the payload that both requests and responses
// are sent as TCP packets.
type packet struct {
	ID   int32
	Type packetType
	Body string
}

// newPacket creates a default packet
func newPacket(t packetType, body string) packet {
	return packet{
		ID:   rand.Int31(),
		Type: t,
		Body: body,
	}
}

// validate validates that the packet size is the expected one.
// The minimum posible value for packet `Size` is 10 bytes
// and the maximum one is 4096.
func (p *packet) validate() error {
	if p.size() > maximumPacketSize {
		return fmt.Errorf("packet size is > %d", maximumPacketSize)
	}
	return nil
}

// serialize transforms the packet to a byte array to be sent to the
// server.
func (p *packet) serialize() ([]byte, error) {
	raw := new(bytes.Buffer)
	binary.Write(raw, binary.LittleEndian, p.size())
	binary.Write(raw, binary.LittleEndian, p.ID)
	binary.Write(raw, binary.LittleEndian, p.Type)

	raw.WriteString(p.Body)
	raw.WriteByte(0) // body payload must be null terminated.
	raw.WriteByte(0)

	return raw.Bytes(), nil
}

func (p *packet) bodySize() int32 {
	return int32(len(p.Body))
}

// size computes the total package size
func (p *packet) size() int32 {
	return (p.bodySize() + minimumPacketSize)
}

// deserialize transforms the payload into a readable Packet.
func (p *packet) deserialize(payload []byte) error {
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
