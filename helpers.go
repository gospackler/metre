// Package metre is used to schedule end execute corn jobs in a simplified fashion
package metre

// buildTaskKey builds a redis unique key for the task
func buildTaskKey(t TaskRecord) string {
    return "cron-" + t.ID + ":" + t.UID
}
