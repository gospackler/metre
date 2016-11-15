package metre

import (
	"testing"

	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

var test = &Task{
	ID:       "Test",
	Interval: "0 * * * * *",
	Schedule: func(m *Master) error {

		for i := 0; i < 10; i++ {
			uid := fmt.Sprintf("%d", i)
			log.Info("Scheduling test ")
			time.Sleep(time.Second)

			// Create a utility function to do it.
			//endChan := make(chan int)
			msg := CreateMsg(Request, "Test", uid, "")
			req := NewScheduleInput(msg)
			m.SchInpChan <- req
			select {
			case resp := <-req.RespChan:
				log.Info("Received resp" + resp)
			case err := <-req.ErrorChan:
				log.Info("Error received" + err.Error())
				break
			}
			log.Info("Outside the infinite loop")
			defer req.Close()
		}
		return nil
	},
	Process: func(msg *MetreMessage) (string, error) {
		log.Info("Processing Test  " + msg.UID)
		return msg.TaskId + msg.UID, nil
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

func TestLife(t *testing.T) {
	/*
		var wg sync.WaitGroup
		met, err := New("127.0.0.1:5555", "127.0.0.1:5556", 1)
		if err != nil {
			t.Errorf("Metre creation error" + err.Error())
		}

		printMsgs(met.MessageChannel)

		met.Add(test)
		met.StartMaster()
		go met.StartSlave()

		met.Schedule(test.ID)
		wg.Add(1)
		wg.Wait()
	*/
}

func TestMasterSlave(t *testing.T) {
	var wg sync.WaitGroup
	dealerUri := "tcp://127.0.0.1:5555"
	routerUri := "tcp://127.0.0.1:5556"
	master, err := NewMaster(routerUri, 2)
	if err != nil {
		t.Error(err.Error())
	}

	StartBroker(dealerUri, routerUri)

	slave, err := NewSlave(dealerUri, 2)
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
