package metre

import (
	log "github.com/Sirupsen/logrus"

	"sync"
	"time"
)

type Task struct {
	sync.Mutex
	ScheduleLock   sync.Mutex
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

func (t *Task) IncrementMessageCount() {
	t.Lock()
	t.MessageCount++
	t.Unlock()
}

func (t *Task) IncrementScheduleCount() {
	t.Lock()
	t.ScheduleCount++
	t.Unlock()
}

func (t *Task) GetMessageCount() int {
	t.Lock()
	count := t.MessageCount
	t.Unlock()
	return count
}

func (t *Task) GetScheduleCount() int {
	t.Lock()
	count := t.ScheduleCount
	t.Unlock()
	return count
}

func (t *Task) SetScheduleDone() {
	t.Lock()
	t.ScheduleDone = true
	t.Unlock()
}

func (t *Task) GetScheduleDone() bool {
	t.Lock()
	schDone := t.ScheduleDone
	t.Unlock()
	return schDone
}

// Zeros the parameters that may change in a run()
func (t *Task) Zero() {
	t.Lock()
	t.MessageCount = 0
	t.ScheduleCount = 0
	t.StartTime = time.Now()
	t.ScheduleDone = false
	t.Unlock()
}
func (t *Task) checkComplete() bool {
	msgCount := t.GetMessageCount()
	schCount := t.GetScheduleCount()
	schDone := t.GetScheduleDone()
	return msgCount == schCount && schDone
}

func (t *Task) SendMessage(msg string) {
	t.MessageChannel <- msg
}

// FIXME: An option is to release the schedule lock if timeout is hit.
func (t *Task) TestTimeOut() {
	if t.TimeOut != 0 {
		time.Sleep(t.TimeOut)
		if !t.checkComplete() {
			t.Lock()
			t.SendMessage("Task hit timeout in " + t.TimeOut.String())
			t.Unlock()
		}
	}
}

func (t *Task) Track(trackMsg *trackMessage) {
	switch trackMsg.MessageType {
	case Status:
		t.IncrementMessageCount()
		if t.checkComplete() {
			t.SendMessage(trackMsg.TaskId + ": Complete")
			t.SendMessage(trackMsg.TaskId + ": Time taken " + time.Since(t.StartTime).String())
			t.ScheduleLock.Unlock()
		}
	case Debug:
		t.SendMessage("DEBUG:" + trackMsg.TaskId + ":" + trackMsg.Message)
	case Error:
		t.SendMessage("ERROR:" + trackMsg.TaskId + ":" + trackMsg.Message)
	default:
		log.Warn("Unknown message type")
	}
}
