package socket

import (
	"encoding/json"
	"net"
	"os"
	"path"

	"github.com/icon-project/centralized-relay/relayer"
)

// temdir is a temporary directory for the socket

var (
	unixSocketPath = path.Join(os.TempDir(), "relayer.sock")
	network        = "unix"
)

type Message struct {
	Event Event
	Data  any
}

type dbServer struct {
	listener net.Listener
	rly      map[string]*relayer.Chain
}

func NewSocket(rly map[string]*relayer.Chain) (*dbServer, error) {
	l, err := net.Listen(network, unixSocketPath)
	if err != nil {
		return nil, err
	}
	return &dbServer{listener: l, rly: rly}, nil
}

// Listen to socket
func (s *dbServer) Listen(errChan chan error) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			errChan <- err
			return
		}
		go s.server(conn)
	}
}

// Send sends message to socket
func (s *dbServer) server(c net.Conn) {
	for {
		buf := make([]byte, 1024)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}
		msg, err := s.Parse(buf[:nr])
		if err != nil {
			return
		}
	}
}

// Parse message from socket
func (s *dbServer) Parse(data []byte) (*Message, error) {
	msg := new(Message)
	err := json.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// Send message to socket
func (s *dbServer) Send(conn net.Conn, msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// Generate Message for the client to write to socket
func (s *dbServer) GenerateMessage(event Event, data any) []byte {
	return &Message{
		Event: event,
		Data:  data,
	}
}

func (s *dbServer) Close() error {
	return s.listener.Close()
}

func (s *dbServer) IsClosed() bool {
	return false
}
