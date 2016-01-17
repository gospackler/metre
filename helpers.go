package metre

// buildTaskKey builds a redis unique key for the task
func buildTaskKey(t TaskRecord) string {
    return "cron-" + t.ID + ":" + t.UID
}
