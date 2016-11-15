package transport

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
)

type ReqConn struct {
	Conn *Connection
}

func NewReqConn(uri string) (*ReqConn, error) {
	conn, err := NewConnection()
	if err != nil {
		return nil, err
	}

	err = conn.Connect(uri, zmq.REQ)
	if err != nil {
		return nil, err
	}

	fmt.Println("Req connected to " + uri)
	return &ReqConn{
		Conn: conn,
	}, nil
}

// This is the client for the request.
// Like a web client, a Send is always followed by a recieve.
func (r *ReqConn) MakeReq(msg string) (string, error) {
	fmt.Println("Making request with message", msg)
	_, err := r.Conn.Sock.Send(msg, zmq.DONTWAIT)
	if err != nil {
		return "", err
	}

	fmt.Println("Sent message waiting for resp")
	reply, err := r.Conn.Sock.Recv(0)
	if err != nil {
		return "", err
	}
	return reply, nil
}

func (r *ReqConn) Close() {
	r.Conn.Close()
}
