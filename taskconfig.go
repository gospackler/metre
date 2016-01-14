// Package metre is used to schedule end execute corn jobs in a simplified fashion
package metre

import (
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
    uid := uuid.NewV4()
    now := time.Now().UTC().Format(time.RFC3339)
    return {id, uid, now, now, "UNSCHEDULED"}
}

// NewTaskRecordWithID takes seed data and returns a full TaskRecord instance
func NewTaskRecordWithID(id, uid) TaskRecord {
    now := time.Now().UTC().Format(time.RFC3339)
    return {id, uid, now, now, "UNSCHEDULED"}
}
