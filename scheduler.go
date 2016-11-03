package metre

import (
	"fmt"
)

type Scheduler struct {
	Queue   Queue
	TaskMap map[string]*Task
}

func NewScheduler(q Queue, t map[string]*Task) Scheduler {
	return Scheduler{q, t}
}

// Schedule schedules on the queue.
func (s Scheduler) Schedule(t TaskRecord) (string, error) {
	key := buildTaskKey(t)
	key, err := schedule(key, t, s.Queue)
	if err == nil {
		task := s.TaskMap[t.ID]
		// FIXME : The operation is better to be done from task itself.
		task.IncrementScheduleCount()
	}
	return key, err
}

// Push the task to the queue
func schedule(k string, t TaskRecord, q Queue) (string, error) {
	str, _ := t.ToString()
	_, qErr := q.Push(str)
	if qErr != nil {
		return k, fmt.Errorf("push queue returned error: %v", qErr)
	}

	t.SetScheduled()

	return k, nil
}
