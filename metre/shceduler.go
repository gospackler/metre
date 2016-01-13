// Package metre is used to schedule end execute corn jobs in a simplified fashion
package metre

import (
    "github.com/robfig/cron"
)

const LOCALHOST string = "127.0.0.1"
const QUEUEPORT string = "5555"
const CACHEPORT string = "6379"

type Metre struct {
    cron cron.Cron
    queue Queue
    cache Cache
    add func(t Task)
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

    q, _ := Queue.New(redisUri)
    c, _ := Cache.New(cacheUri)
    return  Metre{q,c}
}

func (m *Metre) add(t Task) {

}

func (m *Metre) startMaster() {

}

func (m *Metre) startSlave() {

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
