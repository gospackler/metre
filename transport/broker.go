package transport

import (
	"github.com/gospackler/metre/logging"
	zmq "github.com/pebbe/zmq4"
	"go.uber.org/zap"
)

// Create a new queue.
// Port Uri to Bind to.
// Dealer Uri listens to all the master instances.
// Router Uri listens to all the slave instances.
func StartBroker(du string, ru string) error {
	c, err := zmq.NewContext()
	if err != nil {
		return err
	}
	dSock, err := c.NewSocket(zmq.DEALER)
	if err != nil {
		return err
	}
	err = dSock.Bind(du)
	if err != nil {
		return err
	}
	logging.Logger.Debug("dealer socket connected",
		zap.String("uri", du),
	)
	rSock, err := c.NewSocket(zmq.ROUTER)
	if err != nil {
		return err
	}
	err = rSock.Bind(ru)
	if err != nil {
		return err
	}
	logging.Logger.Debug("router socket connected",
		zap.String("uri", ru),
	)
	poller := zmq.NewPoller()
	poller.Add(dSock, zmq.POLLIN)
	poller.Add(rSock, zmq.POLLIN)

	for {
		logging.Logger.Debug("polling for messages on the broker")
		sockets, err := poller.Poll(-1)
		if err != nil {
			return err
		}
		for _, socket := range sockets {
			logging.Logger.Debug("received a message on the socket")
			switch s := socket.Socket; s {
			case dSock:
				msg, err := s.RecvMessage(zmq.DONTWAIT)
				if err != nil {
					return err
				}
				logging.Logger.Debug("received a message in the dealer",
					zap.Strings("message", msg),
				)
				_, err = rSock.SendMessage(msg)
				if err != nil {
					return err
				}
			case rSock:
				msg, err := s.RecvMessage(zmq.DONTWAIT)
				if err != nil {
					return err
				}
				logging.Logger.Debug("received message in router",
					zap.Strings("message", msg),
				)

				_, err = dSock.SendMessage(msg)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
