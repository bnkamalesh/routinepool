package routinepool

import "errors"

var errBlocked = errors.New("Pool stopped, cannot push any further tasks")

//Pool is the struct which handles the worker pool
type Pool struct {
	//size is the pool size, i.e. no.of go routines
	size uint64

	//channelSize is size of the channel buffer to which the works/jobs are sent
	channelSize uint8

	//workerPool is the channel buffer to which all works/jobs are sent
	workerPool chan poolFn

	//active is the channel used to keep count of active works
	active chan *struct{}

	//block if set to true, will block any further jobs being pushed/send to the pool
	block bool

	//quit is used to exit all the routes
	quit chan struct{}
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
	p.quit = make(chan struct{})

	for i := uint64(0); i < p.size; i++ {
		go func(pl *Pool, quit chan struct{}) {
			for {
				select {
				//Listening for quit signal
				case <-quit:
					return

				// Pulling work from the channel buffer
				case work := <-pl.workerPool:
					p.active <- nil
					work()
					if pl.block == false {
						<-p.active
					}
				}
			}
		}(p, p.quit)
	}
}

//Stop stops the pool and exits all the go routines immediately
func (p *Pool) Stop() {
	//block accepting any furhter tasks
	p.block = true

	//Graceful shutdown of pool, makes sure if all pending tasks are completed
	for {
		select {
		case <-p.active:
			if len(p.workerPool) == 0 && len(p.active) == 0 {
				close(p.quit)
				return
			}
		}
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
