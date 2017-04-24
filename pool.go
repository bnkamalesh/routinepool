package routinepool

import "errors"

var errBlocked = errors.New("Pool stopped, cannot push any further tasks")

//Pool is the struct which handles the worker pool
type Pool struct {
	//size is the pool size, i.e. no.of go routines
	size uint64

	//channelSize is size of the channel buffer to which the tasks/jobs are sent
	channelSize uint8

	//workerPool is the channel buffer to which all tasks/jobs are sent
	workerPool chan poolFn

	//block if set to true, will block any further jobs being pushed/send to the pool
	block bool

	//active is the channel used to keep count of active tasks
	active chan *struct{}

	//quit is used to exit all the routes
	quit []chan *struct{}

	//done is used to identify all done/idle routines
	done []chan *struct{}
}

type poolFn func()

type worker struct {
	fn poolFn
}

type workerPool chan []worker

//Push pushes a task to the worker pool
func (p *Pool) Push(work poolFn) error {
	if p.block {
		return errBlocked
	}

	p.workerPool <- work
	return nil
}

//Start starts the worker Q
func (p *Pool) Start() {
	//all individual quit channels are saved in this
	p.quit = make([]chan *struct{}, p.size)
	p.done = make([]chan *struct{}, p.size)

	for i := uint64(0); i < p.size; i++ {
		//an individual quit channel is given to each routine
		q := make(chan *struct{})
		d := make(chan *struct{})
		p.quit[i] = q
		p.done[i] = d

		go func(quit, done chan *struct{}) {
			var shutdown = false
			for {
				select {
				case <-quit:
					shutdown = true

				case work := <-p.workerPool:
					p.active <- nil
					work()
					<-p.active
				}

				if shutdown && len(p.workerPool) == 0 {
					done <- nil
					return
				}
			}
		}(q, d)
	}
}

//Stop stops the pool and exits all the go routines immediately
func (p *Pool) Stop() {
	//block accepting any furhter tasks
	p.block = true

	//Graceful shutdown of pool, makes sure if all pending tasks are completed
	//sending kill signal to invidual routines
	for _, q := range p.quit {
		q <- nil
	}

	//Wait for all gouroutines to send done signal
	for _, d := range p.done {
		<-d
	}
}

//Active returns the number of active jobs
func (p *Pool) Active() int {
	return len(p.active)
}

//Pending returns the number of jobs in queue
func (p *Pool) Pending() int {
	return len(p.workerPool)
}

//New returns a Pool object pointer with all the default values set
func New(pSize uint64, csize uint8) *Pool {
	p := &Pool{
		size:        pSize,
		channelSize: csize,
	}

	//pSize = 0 will fall back to 1
	if p.size == 0 {
		p.size = 1
	}

	if p.channelSize == 0 {
		p.channelSize = 100
	}

	p.workerPool = make(chan poolFn, p.channelSize)

	p.active = make(chan *struct{}, p.size)

	return p
}
