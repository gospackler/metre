package metre

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
)

const queueFlag zmq.Flag = 0
const ZMQ_BLOCKING = iota

type Queue struct {
	URI        string
	PushSocket *zmq.Socket // Socket to push messages to
	PullSocket *zmq.Socket // SOcket to pull mesasges form
}

// BindPush binds to the socket to push
func (q Queue) BindPush() error {
	return q.PushSocket.Bind(q.URI)
}

func (q Queue) BindPull() error {
	return q.PullSocket.Bind(q.URI)
}

// ConnectPull connects to the socket to listen for messages
func (q Queue) ConnectPull() error {
	return q.PullSocket.Connect(q.URI)
}

func (q Queue) ConnectPush() error {
	return q.PushSocket.Connect(q.URI)
}

// Push pushes and element onto the queue
func (q Queue) Push(msg string) (int, error) {
	result, err := q.PushSocket.Send(msg, ZMQ_BLOCKING)
	return result, err
}

// Pop pulls off the last element from the queue
func (q Queue) Pop() string {
	m, _ := q.PullSocket.Recv(queueFlag)
	return m
}

// New acts as a queue constructor
func NewQueue(uri string) (Queue, error) {
	u := "tcp://" + uri
	c, _ := zmq.NewContext()

	pullSoc, pullErr := c.NewSocket(zmq.PULL)
	if pullErr != nil {
		return Queue{}, fmt.Errorf("pull socket initialization failed: %v", pullErr)
	}
	pushSoc, pushErr := c.NewSocket(zmq.PUSH)
	if pushErr != nil {
		return Queue{}, fmt.Errorf("push socket initialization failed: %v", pushErr)
	}

	q := Queue{u, pushSoc, pullSoc}
	return q, nil
}
