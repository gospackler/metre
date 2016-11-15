package transport

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
)

type Process interface {
	Run(string) string
}

type RespConn struct {
	Conn *Connection
}

func NewRespConn(uri string) (*RespConn, error) {
	conn, err := NewConnection()
	if err != nil {
		return nil, err
	}

	err = conn.Connect(uri, zmq.REP)
	if err != nil {
		return nil, err
	}

	fmt.Println("Response Server conected to " + uri)

	return &RespConn{
		Conn: conn,
	}, nil
}

// Call this function from a goroutine
func (r *RespConn) Listen(process Process) error {
	// FIXME : Probably stream the errors and log it in if the server
	// continues to crash with listen errors.
	for {
		fmt.Println("Listener waiting for message .. ")
		//  Wait for next request from client
		req, err := r.Conn.Sock.Recv(0)
		if err != nil {
			return err
		}
		resp := process.Run(req)
		fmt.Println("Received response from Run ", resp)
		// Send reply back to client
		_, err = r.Conn.Sock.Send(resp, zmq.DONTWAIT)
		if err != nil {
			return err
		}
	}
}

func (r *RespConn) Close() {
	r.Conn.Close()
}
