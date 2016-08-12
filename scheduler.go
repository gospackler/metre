package metre

import (
	"errors"
	"fmt"

	"github.com/satori/go.uuid"
)

type Scheduler struct {
	Queue   Queue
	Cache   Cache
	TaskMap map[string]*Task
}

func NewScheduler(q Queue, c Cache, t map[string]*Task) Scheduler {
	return Scheduler{q, c, t}
}

// Schedule schedules a task in the cache and queue if no task is actively waiting to be processed
func (s Scheduler) Schedule(t TaskRecord) (string, error) {
	key := buildTaskKey(t)
	old, _ := s.Cache.Get(key)
	sched := false
	var oldTsk TaskRecord
	var err error

	if old == "" {
		sched = true
	} else {
		oldTsk, _ = ParseTask(old)
		if oldTsk.CanReschedule() {
			sched = true
		} else {
			err = errors.New("A Task with the submitted ID and UID [" + oldTsk.ID + ", " + oldTsk.UID + "] is being processed")
		}
	}

	if sched {
		key, err = schedule(key, t, s.Queue, s.Cache)
		if err == nil {
			task := s.TaskMap[t.ID]
			// FIXME : The operation is better to be done from task itself.
			task.ScheduleCount++
		}
	}
	return key, err
}

// ForceSchedule schedules a task in the cache and queue regardless of tasks actively waiting to be processed
func (s Scheduler) ForceSchedule(t TaskRecord) (string, error) {
	key := buildTaskKey(t)
	old, _ := s.Cache.Get(key)
	var oldTsk TaskRecord

	// affix an additional UID if there was a collision
	if old != "" {
		oldTsk, _ = ParseTask(old)
		if oldTsk.UID == t.UID {
			uid := uuid.NewV4().String()
			t.UID = t.UID + "-" + uid
		}
	}

	return schedule(buildTaskKey(t), t, s.Queue, s.Cache)
}

// SetExpire set teh expiration for a task
func (s Scheduler) SetExpire(t TaskRecord, time int) {
	s.Cache.Expire(buildTaskKey(t), time)
}

// scheduler performs a transaction cache and queue
// The order is: push the task, change the state, update state in the cache.
func schedule(k string, t TaskRecord, q Queue, c Cache) (string, error) {
	str, _ := t.ToString()
	_, qErr := q.Push(str)
	if qErr != nil {
		return k, fmt.Errorf("push queue returned error: %v", qErr)
	}

	t.SetScheduled()

	_, cErr := c.Set(k, str)
	if cErr != nil {
		return k, fmt.Errorf("set cache returned error: %v", cErr)
	}

	return k, nil
}
