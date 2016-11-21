package metre

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"

	log "github.com/Sirupsen/logrus"
)

var test = &Task{
	ID:       "Test",
	Interval: "0 * * * * *",
	Schedule: func(m *Master) error {
		for i := 0; i < 1000; i++ {
			uid := fmt.Sprintf("%d", i)
			log.Info("Scheduling test ")
			resp, err := m.Schedule("Test", uid)
			if err != nil {
				log.Error(err.Error())
			}
			log.Info("Response :" + resp)
		}
		return nil
	},
	Process: func(msg *MetreMessage) (string, error) {
		log.Info("Processing Test  " + msg.UID)
		//	return msg.TaskId + msg.UID, nil
		return "", errors.New("Test error ->" + msg.UID)
	},
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func printMsgs(msgChan chan string) {
	go func() {
		for {
			msg := <-msgChan
			log.Info(msg)
		}
	}()
}

func TestMasterSlave(t *testing.T) {
	var wg sync.WaitGroup
	dealerUri := "tcp://127.0.0.1:5555"
	routerUri := "tcp://127.0.0.1:5556"
	master, err := NewMaster(routerUri, 5)
	if err != nil {
		t.Error(err.Error())
	}

	StartBroker(dealerUri, routerUri)

	slave, err := NewSlave(dealerUri, 6)
	if err != nil {
		t.Error(err.Error())
	}

	met, err := New(master, slave)
	if err != nil {
		t.Error(err.Error())
	}

	// Add the task to both master and slave.
	met.Add(test)
	printMsgs(met.MessageChannel)
	master.Start()

	wg.Add(1)
	wg.Wait()
}
