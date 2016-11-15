package metre

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gospackler/metre/transport"
)

// Slave has a set of workers which listens for request's
// Each of these workers, work in parallel, thanks to ZMQ
type Slave struct {
	TaskMap     map[string]*Task
	workerCount int
	dealerUri   string
}

// FIXME : Use procesParallel in some way.
func NewSlave(dealerUri string, processParallel int) (*Slave, error) {
	if dealerUri == "" {
		dealerUri = LocalHost + ":" + DealerPort
	} else if strings.Index(dealerUri, ":") == 0 {
		dealerUri = LocalHost + ":" + DealerPort
	}
	slave := &Slave{
		TaskMap:     make(map[string]*Task),
		workerCount: processParallel,
		dealerUri:   dealerUri,
	}
	go slave.StartSlave()
	return slave, nil
}

func (s *Slave) AddTask(t *Task) {
	id := t.GetID()
	if _, exists := s.TaskMap[id]; exists {
		log.Warn("attempted to add two tasks with the same ID [" + t.ID + "]")
	}

	s.TaskMap[id] = t
}

func (s *Slave) GetResponse(m string) string {
	msg, err := ParseMessage(m)
	if err != nil {
		return CreateErrorMsg(err, msg.TaskId, msg.UID)
	}
	if s.TaskMap[msg.TaskId] != nil {
		resp, err := s.TaskMap[msg.TaskId].Process(msg)
		if err != nil {
			return CreateErrorMsg(err, msg.TaskId, msg.UID)
		}
		return CreateMsg(Status, msg.TaskId, msg.UID, resp)
	}
	return CreateMsg(Error, msg.TaskId, msg.UID, "Task Id not found in the message")
}

func (s *Slave) Listen(id int) {
	log.Debug("Start Slave ", id)
	respConn, err := transport.NewRespConn(s.dealerUri)
	if err != nil {
		log.Error("Error starting slave " + err.Error())
	}

	err = respConn.Listen(s, id)
	if err != nil {
		log.Error("Slave start error " + err.Error())
	}
	respConn.Close()
}

func (s *Slave) StartSlave() {
	for i := 0; i < s.workerCount; i++ {
		go s.Listen(i)
	}
}
