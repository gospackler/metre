// Package metre is used to schedule end execute corn jobs in a simplified fashion
package metre

import (
    "errors"
)

type Scheduler struct {
    Queue Queue
    Cache Cache
}

func NewScheduler(q Queue, c Cache) Scheduler {
    return Scheduler{q, c}
}


// Schedule schedules a task in the cahce and queue if no task is actively waiting to be processed
func (s Scheduler) Schedule(t TaskRecord) (string, error) {
    key := buildTaskKey(t)
    old, _ := s.Cache.Get(key)
    schedule := false
    var oldTsk TaskRecord
    var err error

    if old == "" {
        schedule = true
    } else {
        oldTsk, _ = ParseTask(old)
        if oldTsk.CanReschedule() {
            schedule = true
        } else {
            err = errors.New("A Task with the submitted ID and UID [" + oldTsk.ID + ", " + oldTsk.UID + "] is being processed")
        }
    }

    if schedule {
        t.SetScheduled()
        str, _ := t.ToString()
        cRsult, cErr := s.Cache.Set(key, str)
        if cErr != nil {
            return key, err
        }
        qResult, qErr := s.Queue.Push(str)
        if qErr != nil {
            return key, err
        }
        return key, nil
    } else {
        return key, err
    }
}


// buildTaskKey builds a redis unique key for the task
func buildTaskKey(t TaskRecord) string {
    return "cron-" + t.ID + ":" + t.UID
}
