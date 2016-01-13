// Package metre is used to schedule end execute corn jobs in a simplified fashion
package metre

type Task struct {
    Type string // Type Type of task (user as class prefix in cache)
    ScheduleStr string // Schedule String in cron notation
    Schedule func(s scheduler.Scheduler, c scheduler.Cache, q scheduler.Queue)
    Process func(s scheduler.Scheduler, c scheduler.Cache, q scheduler.Queue)
}
