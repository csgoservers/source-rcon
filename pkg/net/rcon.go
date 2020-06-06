package net

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

const defaultRCONPort int = 27015

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
	packet := NewPacket()
	packet.Type = serverDataExecCommand
	packet.Body = cmd

	err := r.send(packet)
	if err != nil {
		return nil, err
	}
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
func (r *RemoteConnection) authenticate() error {
	authPacket := NewPacket()
	authPacket.Type = serverDataAuth
	authPacket.Body = r.options.Password
	err := r.send(authPacket)
	return err
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

// receive the responses from the server or an error.
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
