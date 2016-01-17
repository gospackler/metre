// Package metre is used to schedule end execute corn jobs in a simplified fashion
package metre

import (
    "errors"
    "github.com/satori/go.uuid"
)

type Scheduler struct {
    Queue Queue
    Cache Cache
}

func NewScheduler(q Queue, c Cache) Scheduler {
    return Scheduler{q, c}
}


// Schedule schedules a task in the cache and queue if no task is actively waiting to be processed
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
        _, cErr := s.Cache.Set(key, str)
        if cErr != nil {
            return key, err
        }
        _, qErr := s.Queue.Push(str)
        if qErr != nil {
            return key, err
        }

        return key, nil
    } else {
        return key, err
    }
}

// ForceSchedule schedules a task in the cache and queue regardless of tasks actively waiting to be processed
func (s Scheduler) ForceSchedule(t TaskRecord) (string, error) {
    key := buildTaskKey(t)
    old, _ := s.Cache.Get(key)
    var oldTsk TaskRecord
    var err error

    // affix an additional UID if there was a collision
    if old != "" {
        oldTsk, _ = ParseTask(old)
        if oldTsk.UID == t.UID {
            uid := uuid.NewV4().String()
            t.UID  = t.UID + "-" +  uid
        }
    }

    key = buildTaskKey(t)
    t.SetScheduled()
    str, _ := t.ToString()
    _, cErr := s.Cache.Set(key, str)
    if cErr != nil {
        return key, err
    }
    _, qErr := s.Queue.Push(str)
    if qErr != nil {
        return key, err
    }

    return key, nil
}


// buildTaskKey builds a redis unique key for the task
func buildTaskKey(t TaskRecord) string {
    return "cron-" + t.ID + ":" + t.UID
}
