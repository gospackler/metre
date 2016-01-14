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
	mFlagPtr := flag.Bool("master", false, "a boolean to run master")
	sFlagPtr := flag.Bool("slave", false, "a boolean to run slave")
	flag.Parse()

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

- Transfer POC from internal project.
- Everything else

## Authors

- [John M Hofrichter](www.github.com/johnhof)
