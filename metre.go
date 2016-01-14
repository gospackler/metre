// Package metre is used to schedule end execute corn jobs in a simplified fashion
package metre

import (
    "github.com/robfig/cron"
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
    Add func(t Task)
}

// New creates a new scheduler to manage task scheduling and states
func New(cacheUri string, redisUri string) Scheduler {
    if cacheUri == "" {
        cacheUri = LOCALHOST + ":" + CACHEURI
    } else if string.index(cacheUri, ":") == 0 {
        cacheUri = LOCALHOST + ":" + cacheUri
    }

    if cacheUri == "" {
        cacheUri = LOCALHOST + ":" + QUEUEPORT
    } else if string.index(queueUri, ":") == 0 {
        cacheUri = LOCALHOST + ":" + queueUri
    }

    cron := Cron.New()
    q, _ := Queue.New(redisUri)
    c, _ := Cache.New(cacheUri)
    s := Scheduler.new(redisUri, cacheUri)
    return  Metre{cron, q, c}
}

// Add adds a cron job task to schedule and process
func (m *Metre) add(t Task) {
    if _, exists := m.TaskMap[t.ID]; exists {
        panic("Attempted to add two tasks with the same ID [" + t.ID + "]")
    }

    m.TaskMap[t.ID] = t
    m.Cron.Addfunc(t.Interval, func () {
        log.Debug("Scheduling: " + t.ID)
        task := NewTaskRecord(t.ID)
        // t.Schedule(task, m.Scheduler, m.Cache, m.Queue)
    })
}

func (m *Metre) startMaster() {
    m.Cron.Start()
}

func (m *Metre) startSlave() {
    for {
        log.Debug("Waiting fro message")
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
