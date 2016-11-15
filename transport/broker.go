package transport

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
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
	fmt.Println("Dealer bound to " + du)
	err = dSock.Bind(du)
	if err != nil {
		return err
	}

	rSock, err := c.NewSocket(zmq.ROUTER)
	if err != nil {
		return err
	}
	fmt.Println("Router bound to " + ru)
	err = rSock.Bind(ru)
	if err != nil {
		return err
	}
	/*
		go func() {
			for {
				fmt.Println("Waiting for messages on router")
				msg, err := rSock.RecvMessage(0)
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println(msg)
				fmt.Println("Received in router ", msg)
				_, err = dSock.SendMessage(msg, zmq.DONTWAIT)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}()

		go func() {
			for {
				fmt.Println("Waiting for messages on dealer")
				msg, err := dSock.RecvMessage(0)
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println("Received in dealer ", msg)
				_, err = rSock.SendMessage(msg, zmq.DONTWAIT)
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println(msg)
			}
		}()
	*/
	poller := zmq.NewPoller()
	poller.Add(dSock, zmq.POLLIN)
	poller.Add(rSock, zmq.POLLIN)

	for {
		fmt.Println("Polling for messsages on broker")
		sockets, err := poller.Poll(-1)
		if err != nil {
			return err
		}
		for _, socket := range sockets {
			fmt.Println("Received something on the socket ... ")
			switch s := socket.Socket; s {
			case dSock:
				// FIXME - Should we wait ?
				msg, err := s.RecvMessage(zmq.DONTWAIT)
				if err != nil {
					return err
				}
				fmt.Println("Received in dealer ", msg)
				_, err = rSock.SendMessage(msg)
				if err != nil {
					return err
				}
			case rSock:
				msg, err := s.RecvMessage(zmq.DONTWAIT)
				if err != nil {
					return err
				}

				fmt.Println("Received in router ", msg)
				_, err = dSock.SendMessage(msg)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
