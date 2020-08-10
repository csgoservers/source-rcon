package protocol

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	authFailedID    int32  = -1
	emptyPacketBody string = ""
)

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

// RemoteConnection to a Source server
type RemoteConnection struct {
	options       *Options
	connection    net.Conn
	authenticated bool
	mutex         *sync.Mutex
}

// New creates a new RCON for the given host and port.
func New(opts *Options) *RemoteConnection {
	return &RemoteConnection{
		options: opts,
		mutex:   &sync.Mutex{},
	}
}

// ExecCommand sends a command to the server. When this command is executed it will check
// if the connection is authenticated or needs to be authenticated. Note that this method
// makes exclusive use of the underlying connection to make it safe to use in a
// multi routine scenario.
func (r *RemoteConnection) ExecCommand(cmd string) ([]byte, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	// initialize connection to the server if needed
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
	packet := newPacket(serverDataExecCommand, cmd)
	err := r.send(packet)
	if err != nil {
		return nil, err
	}
	// In order to support multiple packet responses for large packets
	// we use a workaround described in the official protocol documentation.
	// https://developer.valvesoftware.com/wiki/Source_RCON_Protocol#Multiple-packet_Responses
	mirrorPacket := newPacket(serverDataResponseValue, emptyPacketBody)
	err = r.send(mirrorPacket)
	if err != nil {
		return nil, err
	}

	// receive response from server. We iterate over the received
	// data to search for the mirror packet (empty). If found, then
	// we break the loop and return the result.
	var raw bytes.Buffer
	for {
		packet, err = r.receive()
		if err != nil {
			return nil, err
		}
		// received mirror packet, return raw bytes.
		if mirrorPacket.ID == packet.ID {
			return raw.Bytes(), nil
		}
		// here we append packet bodies because mirror
		// packet is not received at this moment.
		line := bytes.TrimSpace([]byte(packet.Body))
		raw.Write(line)
	}
}

// Close closes the server connection
func (r *RemoteConnection) Close() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.connection == nil {
		return nil
	}
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
	host := net.JoinHostPort(r.options.Host, strconv.Itoa(r.options.Port))
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
	authPacket := newPacket(serverDataAuth, r.options.Password)
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
func (r *RemoteConnection) send(packet *packet) error {
	err := packet.validate()
	if err != nil {
		return err
	}
	content, err := packet.serialize()
	if err != nil {
		return err
	}
	_, err = r.connection.Write(content)
	return err
}

// receive the responses from the server or an error. Returned
// packet contains the data to be able to correlate the ID with
// the originally sent packet.
func (r *RemoteConnection) receive() (*packet, error) {
	reader := bufio.NewReader(r.connection)
	for {
		chunk := make([]byte, maximumPacketSize)
		num, err := reader.Read(chunk)
		if err != nil {
			return nil, err
		}
		packet := packet{}
		packet.deserialize(chunk[:num])
		return &packet, nil
	}
}
