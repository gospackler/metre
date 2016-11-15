package metre

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gospackler/metre/transport"
)

type Slave struct {
	RespConn *transport.RespConn
	TaskMap  map[string]*Task
}

// FIXME : Use procesParallel in some way.
func NewSlave(dealerUri string, processParallel int) (*Slave, error) {
	if dealerUri == "" {
		dealerUri = LocalHost + ":" + DealerPort
	} else if strings.Index(dealerUri, ":") == 0 {
		dealerUri = LocalHost + ":" + DealerPort
	}
	respConn, err := transport.NewRespConn(dealerUri)
	if err != nil {
		return nil, err
	}

	slave := &Slave{
		RespConn: respConn,
		TaskMap:  make(map[string]*Task),
	}
	slave.StartSlave()
	return slave, nil
}

func (s *Slave) AddTask(t *Task) {
	id := t.GetID()
	if _, exists := s.TaskMap[id]; exists {
		log.Warn("attempted to add two tasks with the same ID [" + t.ID + "]")
	}

	s.TaskMap[id] = t
}

func (s *Slave) Run(m string) string {
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

func (s *Slave) StartSlave() {
	go func() {
		err := s.RespConn.Listen(s)
		if err != nil {
			log.Error("Slave start error " + err.Error())
		}
		defer s.RespConn.Close()
	}()
}
