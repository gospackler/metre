package metre

// task needs to be atomic. How can we achieve it.
// task can exist in both schedule and process, how do we handle it.

// track Messages and Schedules are possible for it.
// Track messages update task and schedules update
// Kill task if it hits timeout.
// Only one task can be scheduled at a time.

import (
	log "github.com/Sirupsen/logrus"

	"sync"
	"time"
)

type Task struct {
	// This is metre's message channel.
	TaskLock          sync.Mutex
	MessageChannel    chan string
	StartTime         time.Time
	TimeOut           time.Duration
	MessageCountLock  sync.Mutex
	MessageCount      int
	ScheduleCountLock sync.Mutex
	ScheduleCount     int
	ScheduleDoneLock  sync.Mutex
	ScheduleDone      bool
	ID                string // Type Type of task (user as class prefix in cache)
	Interval          string // Schedule String in cron notation
	Schedule          func(t TaskRecord, s Scheduler, c Cache, q Queue)
	Process           func(t TaskRecord, s Scheduler, c Cache, q Queue)
}

func (t Task) GetID() string {
	return t.ID
}

func (t Task) GetInterval() string {
	return t.Interval
}

// Zeros the parameters that may change in a run()
func (t *Task) Zero() {
	t.MessageCountLock.Lock()
	t.MessageCount = 0
	t.MessageCountLock.Unlock()
	t.ScheduleCountLock.Lock()
	t.ScheduleCount = 0
	t.ScheduleCountLock.Unlock()
	t.StartTime = time.Now()
	t.ScheduleDoneLock.Lock()
	t.ScheduleDone = false
	t.ScheduleDoneLock.Unlock()
}
func (t *Task) checkComplete() bool {
	t.MessageCountLock.Lock()
	t.ScheduleCountLock.Lock()
	t.ScheduleDoneLock.Lock()
	status := t.MessageCount == t.ScheduleCount && t.ScheduleDone
	t.MessageCountLock.Unlock()
	t.ScheduleCountLock.Unlock()
	t.ScheduleDoneLock.Unlock()
	return status
}

func (t *Task) SendMessage(msg string) {
	t.MessageChannel <- msg
}

func (t *Task) TestTimeOut() {
	if t.TimeOut != 0 {
		time.Sleep(t.TimeOut)
		if !t.checkComplete() {
			t.SendMessage("Task hit timeout in " + t.TimeOut.String())
		}
	}
}

func (t *Task) Track(trackMsg *trackMessage) {
	switch trackMsg.MessageType {
	case Status:
		t.MessageCountLock.Lock()
		t.MessageCount++
		t.MessageCountLock.Unlock()
		if t.checkComplete() {
			t.SendMessage(trackMsg.TaskId + ": Complete")
			t.SendMessage(trackMsg.TaskId + ": Time taken " + time.Since(t.StartTime).String())
			t.TaskLock.Unlock()
		}
	case Debug:
		t.SendMessage("DEBUG:" + trackMsg.TaskId + ":" + trackMsg.Message)
	case Error:
		t.SendMessage("ERROR:" + trackMsg.TaskId + ":" + trackMsg.Message)
	default:
		log.Warn("Unknown message type")
	}
}
