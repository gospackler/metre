// Package metre is used to schedule end execute corn jobs in a simplified fashion
package metre

import (
    "strings"

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
func New(queueUri string, cacheUri string) Metre {
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
    q, _ := NewQueue(queueUri)
    c, _ := NewCache(cacheUri)
    s := NewScheduler(q, c)
    m := make(map[string]Task)
    return  Metre{cron, q, c, s, m}
}

// Add adds a cron job task to schedule and process
func (m *Metre) Add(t Task) {
    if _, exists := m.TaskMap[t.ID]; exists {
        panic("Attempted to add two tasks with the same ID [" + t.ID + "]")
    }

    m.TaskMap[t.ID] = t
    m.Cron.AddFunc(t.Interval, func () {
        log.Debug("Scheduling: " + t.ID)
        // task := NewTaskRecord(t.ID)
        // t.Schedule(task, m.Scheduler, m.Cache, m.Queue)
    })
}

func (m *Metre) StartMaster() {
    m.Cron.Start()
}

func (m *Metre) StartSlave() {
    for {
        log.Debug("Waiting for message")
        msg := m.Queue.Pop()
        log.Debug("Received")
        log.Debug(msg)
    }
}

//
// func (s Scheduler) Schedule(tsk Task) {
//     key := "cron-" + tsk.Type + ":" + tsk.UID
//
//     old, _ := s.Cache.Get(key)
//     schedule := false
//     var oldTsk Task
//
//     if old == "" {
//         schedule = true
//     } else {
//         oldTsk, _ = ParseTask(old)
//         if oldTsk.CanReschedule() {
//             schedule = true
//         }
//     }
//
//     if schedule {
//         log.Info("Scheduling: " + key)
//         tsk.SetScheduled()
//         str, _ := tsk.ToString()
//         s.Queue.Push(str)
//         s.Cache.Set(key, str)
//     } else {
//         log.Warn("Did not schedule: " + key)
//     }
//
// }
//
// func (s Scheduler) Process() {
//     m := s.Queue.Pop()
//     log.Info(m)
// }
