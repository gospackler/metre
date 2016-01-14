// Package metre is used to schedule end execute corn jobs in a simplified fashion
package metre

type Scheduler struct {
    Queue Queue
    Cache Cache
}

func NewScheduler(q Queue, c Cache) Scheduler {
    return Scheduler{q, c}
}
