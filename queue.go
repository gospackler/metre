// Package metre is used to schedule end execute corn jobs in a simplified fashion
package metre

import (
    zmq "github.com/pebbe/zmq4"
    log "github.com/Sirupsen/logrus"
)

const queueFlag zmq.Flag = 0

type Queue struct {
    PushSocket *zmq.Socket // Socket to push messages to
    PullSocket *zmq.Socket // SOcket to pull mesasges form
}

// Push pushes and element onto the queue
func (q Queue) Push(msg string) (int, error) {
    result, err := q.PushSocket.Send(msg, queueFlag)
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
    pushSoc, _ := c.NewSocket(zmq.PUSH)
    pullSoc, _ := c.NewSocket(zmq.PULL)

    pushSoc.Bind(u)
    pullSoc.Bind(u)

    q := Queue{pushSoc, pullSoc}
    return q, nil
}
