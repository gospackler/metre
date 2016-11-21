package transport

import (
	zmq "github.com/pebbe/zmq4"
)

type Connection struct {
	Ctx  *zmq.Context
	Sock *zmq.Socket
}

// Push pushes and element onto the queue
func NewConnection() (*Connection, error) {
	c, err := zmq.NewContext()
	if err != nil {
		return nil, err
	}

	return &Connection{
		Ctx: c,
	}, nil
}

func (c *Connection) Connect(uri string, sockType zmq.Type) (err error) {
	c.Sock, err = c.Ctx.NewSocket(sockType)
	if err != nil {
		return err
	}
	err = c.Sock.Connect(uri)
	return
}

func (c *Connection) Close() {
	defer c.Sock.Close()
}
