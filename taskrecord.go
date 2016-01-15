// Package metre is used to schedule end execute corn jobs in a simplified fashion
package metre

import (
    "encoding/json"
    "time"
    "github.com/satori/go.uuid"
)

type TaskRecord struct {
    ID string
    UID string
    UpdatedAt string
    CreatedAt string
    State string
}

// NewTaskRecord takes seed data and returns a full TaskRecord instance
func NewTaskRecord(id string) TaskRecord {
    uid := uuid.NewV4().String()
    now := time.Now().UTC().Format(time.RFC3339)
    return TaskRecord{id, uid, now, now, "UNSCHEDULED"}
}

// NewTaskRecordWithID takes seed data and returns a full TaskRecord instance
func NewTaskRecordWithID(id string, uid string) TaskRecord {
    now := time.Now().UTC().Format(time.RFC3339)
    return TaskRecord{id, uid, now, now, "UNSCHEDULED"}
}


// ParseTask returns a task instance from a serialized json string
func ParseTask(s string) (TaskRecord, error) {
    // var dat map[string]interface{}
    tsk := TaskRecord{}
    err := json.Unmarshal([]byte(s), &tsk)
    return tsk, err
}

// ToString returns a serialized json string representation of the task
func (t TaskRecord) ToString() (string, error) {
    b, e := json.Marshal(t)
    return string(b), e
}

// TriggerUpdate sets the UpdatedAt property to teh current time
func (t *TaskRecord) TriggerUpdate() {
    t.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
}

// CanReschedule returns whether or not the task can be rescheduled
func (t *TaskRecord) CanReschedule() bool {
    return (t.State == "UNSCHEDULED")
}

// SetProcessing set the state of the task to "PROCESSING"
func (t *TaskRecord) SetProcessing() {
    t.TriggerUpdate()
    t.State = "PROCESSING"
}

// SetProcessing set the state of the task to "SCHEDULED"
func (t *TaskRecord) SetScheduled() {
    t.TriggerUpdate()
    t.State = "SCHEDULED"
}

// SetProcessing set the state of the task to "UNSCHEDULED"
func (t *TaskRecord) SetUnscheduled() {
    t.TriggerUpdate()
    t.State = "UNSCHEDULED"
}
