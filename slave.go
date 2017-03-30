package metre

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/gospackler/metre/logging"
	"github.com/gospackler/metre/transport"
)

var metreSlaveZap = zap.String("metre.type", "slave")

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
		logging.Logger.Warn("attempted to add two tasks with the same ID",
			zap.String("id", t.ID),
		)
	}

	s.TaskMap[id] = t
}

func (s *Slave) logGetResponseInfo(message string, payload *MetreMessage) {
	logging.Logger.Info("metre slave response get",
		metreSlaveZap,
		zap.String("message", message),
		zap.String("task_id", payload.TaskId),
		zap.String("uid", payload.UID),
		zap.String("job_guid", payload.JobGUID),
	)
}

func (s *Slave) logGetResponseErr(message string, payload *MetreMessage) {
	logging.Logger.Error("metre slave response get",
		metreSlaveZap,
		zap.String("message", message),
		zap.String("task_id", payload.TaskId),
		zap.String("uid", payload.UID),
		zap.String("job_guid", payload.JobGUID),
	)
}

func (s *Slave) GetResponse(m string) (ret string) {
	msg, err := ParseMessage(m)
	if err != nil {
		return CreateErrorMsg(err, msg.TaskId, msg.UID, msg.JobGUID)
	}
	s.logGetResponseInfo("Deserialized job", msg)
	if s.TaskMap[msg.TaskId] != nil {
		// Recover from jobs that potentially panic.
		defer func() {
			if r := recover(); r != nil {
				s.logGetResponseErr("Job panicked", msg)
				ret = CreateErrorMsg(fmt.Errorf("Panic while processing task with error: %s", r), msg.TaskId, msg.UID, msg.JobGUID)
			}
		}()

		s.logGetResponseInfo("About to process job", msg)
		resp, err := s.TaskMap[msg.TaskId].Process(msg)
		if err != nil {
			s.logGetResponseErr("Job returned an error while processing", msg)
			return CreateErrorMsg(err, msg.TaskId, msg.UID, msg.JobGUID)
		}
		s.logGetResponseInfo("Process job completed", msg)
		return SerializeMsg(CreateMsg(Status, msg.TaskId, msg.UID, msg.JobGUID, resp))
	}
	s.logGetResponseErr("Task not found for job", msg)
	return SerializeMsg(CreateMsg(Error, msg.TaskId, msg.UID, msg.JobGUID, "Task Id not found in the message"))
}

func (s *Slave) Listen(id int) {
	logging.Logger.Debug("starting slave",
		zap.Int("id", id),
	)
	respConn, err := transport.NewRespConn(s.dealerUri)
	if err != nil {
		logging.Logger.Error("error connection to the dealer",
			zap.Error(err),
		)
	}

	err = respConn.Listen(s, id)
	if err != nil {
		logging.Logger.Error("error starting slave",
			zap.Error(err),
		)
	}
	respConn.Close()
}

func (s *Slave) StartSlave() {
	for i := 0; i < s.workerCount; i++ {
		go s.Listen(i)
	}
}
