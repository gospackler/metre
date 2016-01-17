# metre

Golang cron framework

# DO NOT USE
### This repository is under construction
#### Please check back in a month or so :)

## Example

```Go
// Package main demo's metre
package main

import (
  "github.com/johnhof/metre"
)


func main () {

  // Create
  m := metre.New("127.0.0.1:5555", "127.0.0.1:6379") // (zmq, redis)

  // Add tasks
  m.add(metre.Task{
    ID: "F", // used to prevent collision across tasks
    Interval: "* * * * * *", // Cron Schedule
    Schedule: func(t metre.TaskConfig, s metre.Scheduler, c metre.cache, q metre.queue)  {
      s.Schedule("TestId") // only schedules if "TestID" is not being processed ("F-TestId" not cached in a processing state)
    },
    Process: func(t metre.TaskConfig, s metre.Scheduler, c metre.Cache, q metre.Queue)  {
      log.Info("Processing: " + t.class + "-" + t.uid)
  })

  m.add(metre.Task{
    ID: "B", // used to prevent collision across tasks
    Interval: "@every minute", // Cron Schedule
    Schedule: func(s metre.Scheduler, c metre.Cache, q metre.Queue)  {
      s.Schedule("TestId") // only schedules if "TestID" is not being processed ("F-TestId" not cached in a processing state)
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

## TODO

- Scheduler seems to only queue every other iteration
- process function map

## Authors

- [John M Hofrichter](www.github.com/johnhof)
