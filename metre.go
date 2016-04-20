// Package metre is used to schedule end execute cron jobs in a simplified fashion
package metre

import (
    "strings"
    "errors"

    "github.com/robfig/cron"
    log "github.com/Sirupsen/logrus"
)

const LOCALHOST string = "127.0.0.1" // Default host for cache and queue
const QUEUEPORT string = "5555" // Default port for queue
const CACHEPORT string = "6379" // Default port for cache

type Metre struct {
    Cron cron.Cron
    Queue Queue
    Cache Cache
    Scheduler Scheduler
    TaskMap map[string]Task
    // Add func(t Task)
}

// New creates a new scheduler to manage task scheduling and states
func New(queueUri string, cacheUri string) (Metre, error) {
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

    cron := *cron.New()
    c, cErr := NewCache(cacheUri)
    if cErr != nil {
        return Metre{}, cErr
    }
    q, qErr := NewQueue(queueUri)
    if qErr != nil {
        return Metre{}, qErr
    }

    s := NewScheduler(q, c)
    m := make(map[string]Task)
    return  Metre{cron, q, c, s, m}, nil
}

// Add adds a cron job task to schedule and process
func (m *Metre) Add(t Task) {
    if _, exists := m.TaskMap[t.ID]; exists {
        panic("attempted to add two tasks with the same ID [" + t.ID + "]")
    }

    m.TaskMap[t.ID] = t
    m.Cron.AddFunc(t.Interval, func () {
        t.Schedule(NewTaskRecord(t.ID), m.Scheduler, m.Cache, m.Queue)
    })
}

// Schedule schedules a singular cron task
func (m *Metre) Schedule(ID string) (string, error) {
    e := m.Queue.BindPush()
    if e != nil{
        return "", nil
    }
    t, ok := m.TaskMap[ID]
    if ok == false {
        return "", errors.New("task [" + ID + "] not recognized")
    }

    tr := NewTaskRecord(t.ID)
    t.Schedule(tr, m.Scheduler, m.Cache, m.Queue)
    return buildTaskKey(tr), nil
}

// Scheduler processes a singular cron task
func (m *Metre) Process(ID string) (string, error) {
    t, ok := m.TaskMap[ID]
    if ok == false {
        return "", errors.New("task [" + ID + "] not recognized")
    }

    tr := NewTaskRecord(t.ID)
    t.Process(tr, m.Scheduler, m.Cache, m.Queue)
    return buildTaskKey(tr), nil
}

func (m *Metre) StartMaster() {
    e := m.Queue.BindPush()
    if e != nil {
        panic(e)
    }
    m.Cron.Start()
}

func (m *Metre) StartSlave() {
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
        tsk := m.TaskMap[tr.ID]
        go tsk.Process(tr, m.Scheduler, m.Cache, m.Queue)
    }
}
