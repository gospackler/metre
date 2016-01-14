// Package metre is used to schedule end execute corn jobs in a simplified fashion
package metre

type Task struct {
    ID string // Type Type of task (user as class prefix in cache)
    Interval string // Schedule String in cron notation
    Schedule func(t TaskRecord, s Scheduler, c Cache, q Queue)
    Process func(t TaskRecord, s Scheduler, c Cache, q Queue)
}
