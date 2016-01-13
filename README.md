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
  "sync"
  "flag"
  "github.com/johnhof/metre"
)


func main () {
  var wg sync.WaitGroup
	mFlagPtr := flag.Bool("master", false, "a boolean to run master")
	sFlagPtr := flag.Bool("slave", false, "a boolean to run slave")
	flag.Parse()

  // Create
  m := metre.New("127.0.0.1:5555", "127.0.0.1:6379") // (zmq, redis)

  // Add tasks
  m.add(metre.Task{
    class: "F", // used to prevent collision across tasks
    interval: "* * * * * *", // Cron schedule
    schedule: func(s metre.Scheduler, c metre.cache, q metre.queue)  {
      s.schedule("TestId") // only schedules if "TestID" is not being processed ("F-TestId" not cached in a processing state)
    },
    process: func(s metre.Scheduler, c metre.cache, q metre.queue)  {
      log.Info("Processing: test")
    },
  })

  m.add(metre.Task{
    class: "B", // used to prevent collision across tasks
    interval: "@every minute", // Cron schedule
    schedule: func(s metre.Scheduler, c metre.Cache, q metre.Queue)  {
      s.schedule("TestId") // only schedules if "TestID" is not being processed ("F-TestId" not cached in a processing state)
    },
    process: func(t metre.TaskInstance, s metre.Scheduler, c metre.Cache, q metre.Qßueue)  {
      log.Info("Processing: " + t.class + "-" + t.uid)
    },
  })

  // Execute master
  if *mFlagPtr {
    log.Info("Executing Master...")
    wg.Add(1)
    go m.runMaster()
  }

  // Execute slave
  if *sFlagPtr {
    log.Info("Executing Slave...")
    wg.Add(1)
    go m.runSlave()
  }

  // Keep the server alive
	wg.Wait();
}
```

## TODO

- Transfer POC from internal project.
- Everything else

## Authors

- [John M Hofrichter](www.github.com/johnhof)
