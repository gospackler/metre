package metre

import (
	"sync"
	"time"

	"github.com/gospackler/metre/logging"

	"go.uber.org/zap"
)

type Task struct {
	sync.Mutex
	MessageChannel chan string
	StartTime      time.Time
	ID             string // Type Type of task (user as class prefix in cache)
	Interval       string // Schedule String in cron notation
	Schedule       func(*Master) error
	Process        func(*MetreMessage) (string, error)
}

func (t Task) GetID() string {
	return t.ID
}

// Zeros the parameters that may change in a run()
func (t *Task) Zero() {
	t.StartTime = time.Now()
}

func (t *Task) SendMessage(msg string) {
	t.MessageChannel <- msg
}

func (t *Task) Evaluate(msg *MetreMessage) {
	switch msg.MessageType {
	case Status:
		t.SendMessage("STATUS:" + msg.TaskId + ":" + msg.Message)
	case Debug:
		t.SendMessage("DEBUG:" + msg.TaskId + ":" + msg.Message)
	case Error:
		t.SendMessage("ERROR:" + msg.TaskId + ":" + msg.Message)
	default:
		logging.Logger.Warn("unknown message type",
			zap.Any("message_type", msg.MessageType),
		)
	}
}
