package metre

import (
	log "github.com/Sirupsen/logrus"

	"time"
)

type Task struct {
	MessageChannel chan string
	StartTime      time.Time
	TimeOut        time.Duration
	MessageCount   int
	ScheduleCount  int
	ScheduleDone   bool
	ID             string // Type Type of task (user as class prefix in cache)
	Interval       string // Schedule String in cron notation
	Schedule       func(t TaskRecord, s Scheduler, c Cache, q Queue)
	Process        func(t TaskRecord, s Scheduler, c Cache, q Queue)
}

func (t Task) GetID() string {
	return t.ID
}

func (t Task) GetInterval() string {
	return t.Interval
}

// Zeros the parameters that may change in a run()
func (t *Task) Zero() {
	t.MessageCount = 0
	t.ScheduleCount = 0
	t.StartTime = time.Now()
	t.ScheduleDone = false
}
func (t *Task) checkComplete() bool {
	return t.MessageCount == t.ScheduleCount && t.ScheduleDone
}

func (t *Task) sendMessage(msg string) {
	t.MessageChannel <- msg
}

func (t *Task) TestTimeOut() {
	if t.TimeOut != 0 {
		time.Sleep(t.TimeOut)
		if !t.checkComplete() {
			t.sendMessage("Task hit timeout in " + t.TimeOut.String())
		}
	}
}

func (t *Task) Track(trackMsg *trackMessage) {
	switch trackMsg.MessageType {
	case Status:
		t.MessageCount++
		if t.checkComplete() {
			t.sendMessage(trackMsg.TaskId + ":Complete")
			t.sendMessage(trackMsg.TaskId + ": Time taken " + time.Since(t.StartTime).String())
		}
	case Debug:
		t.sendMessage("DEBUG:" + trackMsg.TaskId + ":" + trackMsg.Message)
	case Error:
		t.sendMessage("ERROR:" + trackMsg.TaskId + ":" + trackMsg.Message)
	default:
		log.Warn("Unknown message type")
	}
}
