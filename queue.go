// Package metre is used to schedule end execute corn jobs in a simplified fashion
package metre

import (
    zmq "github.com/pebbe/zmq4"
    log "github.com/Sirupsen/logrus"
)

type Queue struct {
    PushSocket *zmq.Socket // Socket to push messages to
    PullSocket *zmq.Socket // SOcket to pull mesasges form
}

// Push pushes and element onto the queue
func (q Queue) Push(msg string) {
    _, err := q.PushSocket.Send(msg, 0)
    if (err != nil) {
        log.Warn(err)
    }
}

// Pop pulls off the last element from the queue
func (q Queue) Pop() string {
    log.Debug("before")
    m, _ := q.PullSocket.Recv(0)
    log.Debug("after")
    return m
}

// New acts as a queue constructor
func NewQueue(uri string) (Queue, error) {
    u := "tcp://" + uri
    c, _ := zmq.NewContext()
    pushSoc, _ := c.NewSocket(zmq.PUSH)
    pullSoc, _ := c.NewSocket(zmq.PULL)

    pushSoc.Connect(u)
    pullSoc.Connect(u)

    q := Queue{pushSoc, pullSoc}
    return q, nil
}
