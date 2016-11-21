package transport

import (
	log "github.com/Sirupsen/logrus"
	zmq "github.com/pebbe/zmq4"
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
	log.Debug("Dealer bound to " + du)
	err = dSock.Bind(du)
	if err != nil {
		return err
	}

	rSock, err := c.NewSocket(zmq.ROUTER)
	if err != nil {
		return err
	}
	log.Debug("Router bound to " + ru)
	err = rSock.Bind(ru)
	if err != nil {
		return err
	}
	poller := zmq.NewPoller()
	poller.Add(dSock, zmq.POLLIN)
	poller.Add(rSock, zmq.POLLIN)

	for {
		log.Debug("Polling for messsages on broker")
		sockets, err := poller.Poll(-1)
		if err != nil {
			return err
		}
		for _, socket := range sockets {
			log.Debug("Received something on the socket ... ")
			switch s := socket.Socket; s {
			case dSock:
				msg, err := s.RecvMessage(zmq.DONTWAIT)
				if err != nil {
					return err
				}
				log.Debug("Received in dealer ", msg)
				_, err = rSock.SendMessage(msg)
				if err != nil {
					return err
				}
			case rSock:
				msg, err := s.RecvMessage(zmq.DONTWAIT)
				if err != nil {
					return err
				}

				log.Debug("Received in router ", msg)
				_, err = dSock.SendMessage(msg)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
