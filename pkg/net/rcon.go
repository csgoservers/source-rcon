package net

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

const authFailedID int32 = -1

var (
	errorPacketIDNotMatch = errors.New("packet ID doesn't match for authentication")
	errorAuthFailed       = errors.New("authentication failed")
)

// Options for the remote connection to be stablished
type Options struct {
	Host     string
	Password string
	Port     int
	Timeout  time.Time
}

func (o *Options) url() string {
	return fmt.Sprintf("%s:%d", o.Host, o.Port)
}

// RemoteConnection to a Source server
type RemoteConnection struct {
	options       *Options
	connection    net.Conn
	authenticated bool
}

// New creates a new RCON for the given host and port.
func New(opts *Options) *RemoteConnection {
	return &RemoteConnection{options: opts}
}

// ExecCommand sends a command to the server. When this command is executed it check
// if the connection is authenticated or needs to be authenticated.
func (r *RemoteConnection) ExecCommand(cmd string) ([]byte, error) {
	if r.connection == nil {
		err := r.initialize()
		if err != nil {
			return nil, err
		}
	}
	// first check if is not authenticated and connection options
	// contains a password to try to authenticate the connection to
	// the server.
	if !r.authenticated && len(r.options.Password) > 0 {
		err := r.authenticate()
		if err != nil {
			return nil, err
		}
	}
	// we send the given command but first we create a packet to
	// be sent.
	packet := NewPacket(serverDataExecCommand, cmd)
	err := r.send(packet)
	if err != nil {
		return nil, err
	}
	// receive response from server
	result, err := r.receive()
	if err != nil {
		return nil, err
	}
	log.Printf("Received packet: %s", result.String())
	return nil, nil
}

// Close closes the server connection
func (r *RemoteConnection) Close() error {
	err := r.connection.Close()
	if err != nil {
		return err
	}
	r.authenticated = false
	r.connection = nil
	return nil
}

// initialize is the method to setup connection with the server
func (r *RemoteConnection) initialize() error {
	host := r.options.url()
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}
	conn.SetDeadline(r.options.Timeout)
	r.connection = conn
	return nil
}

// authenticate is used to authenticate the connection with the server.
// When the server receives an auth request, it will respond with an empty
// response followed immediately by a server data auth response indicating
// if the authentication succeeded or failed. Note that the status code
// is returned in the packet id field. If the ID is equals to -1 then the
// authentication failed.
func (r *RemoteConnection) authenticate() error {
	authPacket := NewPacket(serverDataAuth, r.options.Password)
	err := r.send(authPacket)

	// here we expect an empty response form the server. ID from
	// auth request packet mus be the same that this first one.
	result, err := r.receive()
	if err != nil {
		return err
	}
	if authPacket.ID != result.ID {
		return errorPacketIDNotMatch
	}
	// this second packet receives the actual result of the authentication
	// process. ID must be the same as the original auth packet. If packet
	// ID is -1 then authentication failed.
	result, err = r.receive()
	if err != nil {
		return err
	}
	if result.ID == authFailedID {
		return errorAuthFailed
	}
	if authPacket.ID != result.ID {
		return errorPacketIDNotMatch
	}
	r.authenticated = true
	return nil
}

// send the given packge after validate it. Note that
// this packet is serialized into a byte array and
// sent using the given connection.
func (r *RemoteConnection) send(packet *Packet) error {
	err := packet.Validate()
	if err != nil {
		return err
	}
	content, err := packet.Serialize()
	if err != nil {
		return err
	}
	log.Printf("Sent packet: %v", packet.String())
	_, err = r.connection.Write(content)
	return err
}

// receive the responses from the server or an error. Returned
// packet contains the data to be able to correlate the ID with
// the originally sent packet.
func (r *RemoteConnection) receive() (*Packet, error) {
	reader := bufio.NewReader(r.connection)
	for {
		chunk := make([]byte, maximumPacketSize)
		num, err := reader.Read(chunk)
		if err != nil {
			return nil, err
		}
		packet := Packet{}
		packet.Deserialize(chunk[:num])
		return &packet, nil
	}
}
