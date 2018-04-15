[![](https://travis-ci.org/bnkamalesh/routinepool.svg?branch=master)](https://travis-ci.org/bnkamalesh/routinepool)
[![](https://cover.run/go/github.com/bnkamalesh/routinepool.svg?tag=golang-1.10)](https://cover.run/go/github.com/bnkamalesh/routinepool)
[![](https://goreportcard.com/badge/github.com/bnkamalesh/routinepool)](https://goreportcard.com/report/github.com/bnkamalesh/routinepool)
[![](https://api.codeclimate.com/v1/badges/d943c47434acfe470a88/maintainability)](https://codeclimate.com/github/bnkamalesh/routinepool/maintainability)
[![](https://godoc.org/github.com/nathany/looper?status.svg)](http://godoc.org/github.com/bnkamalesh/routinepool)

## Routinepool

Routinepool is a simple Go library to start a worker pool. The example provided 
in the test illustrates everything what this worker pool can do. 
e.g. See number of active jobs, pending jobs, stop worker pool etc.

1. Start worker pool
2. Push tasks/work to the pool
3. See active number of tasks
4. See pending tasks

### Methods available

1. `New()` - returns a new Pool pointer with all the configurations set, ready to be used
2. `p.Start()` - starts the workerpool and waits for tasks to be pushed
3. `p.Push(<task>)` - pushes a work/task to the queue
4. `p.Active()` - returns the number of tasks which are actively running
5. `p.Pending()` - returns the number of tasks which are pending in the queue
6. `p.Stop()` - gracefully shutsdown the workerpool, i.e. waits for the queue to be empty and shutdown