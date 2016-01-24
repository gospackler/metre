# metre

Golang cron framework. Scheduling through ØMQ push-pull and state managed through redis

## Requirements

- [ØMQ](http://zeromq.org/)
- [Redis](http://redis.io/)

## Installation

`go get github.com/johnhof/metre`

## Documentation

- API References
  - [gowalker.org](https://godoc.org/github.com/johnhof/metre)
  - [godoc.com](https://gowalker.org/github.com/johnhof/metre)

## Example

```Go
// Package main demo's metre
package main

import (
  "github.com/johnhof/metre"
)


func main () {

  // Create
  m, _ := metre.New("127.0.0.1:5555", "127.0.0.1:6379") // (zmq, redis)

  // Add tasks
  m.add(metre.Task{
    ID: "F", // used to prevent collision across tasks
    Interval: "* * * * * *", // Cron expression
    Schedule: func(t metre.TaskConfig, s metre.Scheduler, c metre.cache, q metre.queue)  {
      t.UID = "TestID" // overwrite the generated UID with a static namespace
      s.Schedule(t) // only schedules if "TestID" is not being processed ("F-TestId" not cached in a processing state)
      s.SetExpire(t, 15000) // set expire in milliseconds
    },
    Process: func(t metre.TaskConfig, s metre.Scheduler, c metre.Cache, q metre.Queue)  {
      log.Info("Processing: " + t.class + "-" + t.uid)
  })

  m.add(metre.Task{
    ID: "B", // used to prevent collision across tasks
    Interval: "@every minute", // Cron Schedule
    Schedule: func(s metre.Scheduler, c metre.Cache, q metre.Queue)  {
      t.UID = "TestID" // overwrite the generated UID with a static namespace
      s.ForceSchedule(t) // schedule regardless of task state
    },
    Process: func(t metre.TaskInstance, s metre.Scheduler, c metre.Cache, q metre.Queue)  {
      log.Info("Processing: " + t.class + "-" + t.uid)
    },
  })

  go m.RunSlave() // run slave in background
  m.RunMaster() // run master

  // this approach means we can run master and slave
  // on separate servers from the same code base
  // allowing us to scale them separately
}
```

## Authors

- [John M Hofrichter](www.github.com/johnhof)
