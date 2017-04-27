## Routinepool

[GoDoc](https://godoc.org/github.com/bnkamalesh/routinepool)

Routinepool is a simple Go library to start a worker pool. The [sample](https://github.com/bnkamalesh/routinepool/tree/master/cmd) included with the repo illustrates everything what this worker pool can do. e.g. See number of active jobs, pending jobs, stop worker pool etc.

### Methods available

1. `New()` - returns a new Pool pointer with all the configurations set, ready to be used
2. `p.Start()` - starts the workerpool and waits for tasks to be pushed
3. `p.Push(<task>)` - pushes a work/task to the queue
4. `p.Active()` - returns the number of tasks which are actively running
5. `p.Pending()` - returns the number of tasks which are pending in the queue
6. `p.Stop()` - gracefully shutsdown the workerpool, i.e. waits for the queue to be empty and shutdown