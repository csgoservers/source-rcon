package net

import (
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
	log.Printf("Connection initialized for %s", r.options.url())
	// first check if is not authenticated and connection options
	// contains a password to try to authenticate the connection to
	// the server.
	if !r.authenticated && len(r.options.Password) > 0 {
		err := r.authenticate()
		if err != nil {
			return nil, err
		}
	}
	log.Println("Connection to server is authenticated")
	// we send the given command but first we create a packet to
	// be sent.
	packet := NewPacket()
	packet.Type = serverDataExecCommand
	packet.Body = cmd

	err := r.send(packet)
	if err != nil {
		return nil, err
	}
	log.Printf("Command %s sent to the server", cmd)
	return nil, nil
}

// Close closes the server connection
func (r *RemoteConnection) Close() {
	r.connection.Close()
	r.authenticated = false
	r.connection = nil
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
	return r.send(authPacket)
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
	_, err = r.connection.Write(content)
	return err
}
