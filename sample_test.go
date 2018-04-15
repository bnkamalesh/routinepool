package routinepool

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func ExamplePool() {
	var completed uint64
	const (
		jobs     = 32
		poolsize = 4
	)
	p := New(poolsize, jobs)

	//Start the routine pool and wait for jobs/tasks
	p.Start()

	for i := 0; i < jobs; i++ {
		go p.Push(func() {
			time.Sleep(time.Millisecond * 50)
			atomic.AddUint64(&completed, 1)
		})
	}

	//Gracefully shutdown the routine pool.
	p.Stop()

	fmt.Println("Active jobs:", p.Active(), "Pending jobs:", p.Pending(), "Completed:", completed, "out of", jobs)
	// Output: Active jobs: 0 Pending jobs: 0 Completed: 32 out of 32
}

func TestPool(t *testing.T) {
	var completed uint64
	const (
		jobs     = 32
		poolsize = 4
	)
	p := New(poolsize, jobs)

	//Start the routine pool and wait for jobs/tasks
	p.Start()
	for i := 0; i < jobs; i++ {
		// Pushing tasks to the pool
		go p.Push(func() {
			time.Sleep(time.Millisecond * 50)
			atomic.AddUint64(&completed, 1)
		})
	}
	//Gracefully shutdown the routine pool.
	p.Stop()
	err := p.Stop()
	if err != ErrShutdown {
		t.Log("Expected", ErrShutdown, "\n got", err)
		t.Fail()
	}

	if jobs != completed {
		t.Log("Only completed", completed, "out of", jobs)
		t.Fail()
	}

}
