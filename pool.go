package routinepool

//Pool is the struct which handles the worker pool
type Pool struct {
	//size is the pool size, i.e. no.of go routines
	size uint64

	//channelSize is size of the channel buffer to which the works/jobs are sent
	channelSize uint64

	//active is the no.of work/tasks that are active
	active uint64

	//workerPool is the channel buffer to which all works/jobs are sent
	workerPool chan poolFn

	//counterChan is the channel used to keep count of active works
	counterChan chan *struct{}

	//quit is used to exit all the routes
	quit chan struct{}
}

type poolFn func()

type worker struct {
	fn poolFn
}

type workerPool chan []worker

//Push pushes a task to the worker pool
func (p *Pool) Push(work poolFn) {
	p.workerPool <- work
}

//Start starts the worker Q
func (p *Pool) Start() {
	p.quit = make(chan struct{})

	println("Starting routine pool")

	for i := 0; i < int(p.size); i++ {
		go func(pl *Pool, quit chan struct{}) {
			for {
				select {
				//Listening for quit signal
				case <-quit:
					return

				// Pulling work from the channel buffer
				case work := <-pl.workerPool:
					p.counterChan <- nil
					work()
					<-p.counterChan
				}
			}
		}(p, p.quit)
	}
}

//Stop stops the pool and exits all the go routines immediately
func (p *Pool) Stop() {
	close(p.quit)
	println("Stopped routine pool")
}

//GetActive returns the number of active jobs
func (p *Pool) GetActive() int {
	return len(p.counterChan)
}

//GetPending returns the number of jobs in queue
func (p *Pool) GetPending() int {
	return len(p.workerPool)
}

//New returns a Pool object pointer with all the default values set
func New(pSize uint64, csize uint64) *Pool {
	p := &Pool{
		size:        pSize,
		channelSize: csize,
	}

	//pSize = 0 will fall back to 1
	if p.size == 0 {
		p.size = 1
	}

	if p.channelSize == 0 {
		p.channelSize = p.size + 50
	}

	p.workerPool = make(chan poolFn, p.channelSize)
	p.counterChan = make(chan *struct{}, p.channelSize)

	return p
}
