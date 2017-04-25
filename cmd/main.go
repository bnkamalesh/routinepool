package main

import (
	"sync/atomic"
	"time"

	"github.com/bnkamalesh/routinepool"
)

const (
	jobs     = 64
	poolsize = 4
)

var completed uint64

//Do waits for 1 second and prints "Done!" to stdout
func Do() {
	time.Sleep(time.Millisecond * 50)
	atomic.AddUint64(&completed, 1)
}

func main() {
	p := routinepool.New(poolsize, jobs)

	//Start the routine pool and wait for jobs/tasks
	p.Start()

	for i := 0; i < jobs; i++ {
		go p.Push(Do)
	}

	//Print number of active jobs
	println("Active jobs:", p.Active())

	// time.Sleep(time.Second * 2)

	//Print number of jobs in queue
	println("Pending jobs:", p.Pending())

	//Gracefully shutdown the routine pool.
	p.Stop()

	println("Completed:", completed, "out of", jobs)
}
