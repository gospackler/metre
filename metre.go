// Package metre is used to schedule end execute cron jobs in a simplified fashion
package metre

import (
	"errors"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/robfig/cron"
)

const LOCALHOST string = "127.0.0.1" // Default host for cache and queue
const QUEUEPORT string = "5555"      // Default port for queue
const TRACKQUEUEPORT string = "5556" // Default port for queue
const CACHEPORT string = "6379"      // Default port for cache

type Metre struct {
	Cron           cron.Cron
	Queue          Queue
	TrackQueue     Queue
	Cache          Cache
	Scheduler      Scheduler
	TaskMap        map[string]*Task
	MessageChannel chan string
}

// New creates a new scheduler to manage task scheduling and states
func New(queueUri string, trackQueueUri string, cacheUri string) (*Metre, error) {
	if cacheUri == "" {
		cacheUri = LOCALHOST + ":" + CACHEPORT
	} else if strings.Index(cacheUri, ":") == 0 {
		cacheUri = LOCALHOST + ":" + cacheUri
	}

	if queueUri == "" {
		queueUri = LOCALHOST + ":" + QUEUEPORT
	} else if strings.Index(queueUri, ":") == 0 {
		queueUri = LOCALHOST + ":" + queueUri
	}

	if trackQueueUri == "" {
		trackQueueUri = LOCALHOST + ":" + TRACKQUEUEPORT
	} else if strings.Index(trackQueueUri, ":") == 0 {
		trackQueueUri = LOCALHOST + ":" + trackQueueUri
	}

	cron := *cron.New()
	c, cErr := NewCache(cacheUri)
	if cErr != nil {
		return nil, cErr
	}
	q, qErr := NewQueue(queueUri)
	if qErr != nil {
		return nil, qErr
	}

	t, tErr := NewQueue(trackQueueUri)
	if tErr != nil {
		return nil, tErr
	}

	m := make(map[string]*Task)
	s := NewScheduler(q, c, m)
	msgChan := make(chan string)
	return &Metre{cron, q, t, c, s, m, msgChan}, nil
}

// Add adds a cron job task to schedule and process
func (m *Metre) Add(t *Task) {
	id := t.GetID()
	if _, exists := m.TaskMap[id]; exists {
		panic("attempted to add two tasks with the same ID [" + t.ID + "]")
	}

	m.TaskMap[id] = t
	t.MessageChannel = m.MessageChannel
	m.Cron.AddFunc(t.Interval, func() {
		go m.scheduleFromId(id)
	})
}

func (m *Metre) scheduleFromId(ID string) (string, error) {
	var t *Task
	var ok bool
	if t, ok = m.TaskMap[ID]; !ok {
		return "", errors.New("task [" + ID + "] not recognized")
	}

	tr := NewTaskRecord(ID)
	// Making sure the next run is not affected by previous runs.
	t.TaskLock.Lock()
	t.Zero()
	go t.TestTimeOut()
	t.SendMessage(t.ID + ": Scheduled")
	t.Schedule(tr, m.Scheduler, m.Cache, m.Queue)
	t.ScheduleDoneLock.Lock()
	t.ScheduleDone = true
	t.ScheduleDoneLock.Unlock()
	return buildTaskKey(tr), nil
}

// Schedule schedules a singular cron task
func (m *Metre) Schedule(ID string) (string, error) {
	e := m.Queue.BindPush()
	if e != nil {
		return "", nil
	}
	return m.scheduleFromId(ID)
}

// Scheduler processes a singular cron task
func (m *Metre) Process(ID string) (string, error) {
	t, ok := m.TaskMap[ID]
	if ok == false {
		return "", errors.New("task [" + ID + "] not recognized")
	}

	tr := NewTaskRecord(ID)
	t.Process(tr, m.Scheduler, m.Cache, m.Queue)
	return buildTaskKey(tr), nil
}

func (m *Metre) StartMaster() {
	e := m.Queue.BindPush()
	if e != nil {
		panic(e)
	}
	m.Cron.Start()
	go m.track()
}

// This function tracks if the schedules get completed.
func (m *Metre) track() {
	e := m.TrackQueue.ConnectPull()
	if e != nil {
		panic(e)
	}
	for {
		msg := m.TrackQueue.Pop()
		trackMsg, err := parseMessage(msg)
		if err != nil {
			log.Warn("Error parsing track message" + err.Error())
		}
		task := m.TaskMap[trackMsg.TaskId]
		// FIXME: Race can be caused here.
		task.Track(trackMsg)
	}
}

func (m *Metre) runAndSendComplete(tr TaskRecord) {
	tsk := m.TaskMap[tr.ID]
	tsk.Process(tr, m.Scheduler, m.Cache, m.Queue)
	// The content do not matter as this is used for counting messages.
	// TODO Could add the success or failure of the task in future.
	statusMsg := createMsg(Status, tr.ID, tr.UID, "")
	_, err := m.TrackQueue.Push(statusMsg)
	if err != nil {
		log.Warn("Error while pushing completed status for a process")
	}
}

func (m *Metre) StartSlave() {
	err := m.TrackQueue.BindPush()
	if err != nil {
		log.Warn("Track queue not working properly ", err.Error())
	}

	e := m.Queue.ConnectPull()
	if e != nil {
		panic(e)
	}
	for {
		msg := m.Queue.Pop()
		tr, _ := ParseTask(msg)
		if tr.ID == "" || tr.UID == "" {
			log.Warn("Failed to parse task from message: " + msg)
			continue
		}

		m.Cache.Delete(buildTaskKey(tr))
		go m.runAndSendComplete(tr)
	}
}
