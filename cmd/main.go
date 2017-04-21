package main

import (
	"time"

	"github.com/bnkamalesh/routinepool"
)

const (
	jobs     = 32
	poolsize = 3
)

//Do waits for 1 second and prints "Done!" to stdout
func Do() {
	time.Sleep(time.Second * 1)
	println("Done!")

}

func main() {
	p := routinepool.New(poolsize, jobs)

	//Start the routine pool and wait for jobs/tasks
	p.Start()

	for i := 0; i < jobs; i++ {
		p.Push(Do)
	}

	time.Sleep(time.Second * 3)

	//Print number of active jobs
	println("Active jobs:", p.Active())

	time.Sleep(time.Second * 1)

	//Print number of jobs in queue
	println("Pending jobs:", p.Pending())

	//Stop the routine pool
	p.Stop()
}
