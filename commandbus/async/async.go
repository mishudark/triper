package async

import (
	"github.com/golang/glog"
	"github.com/mishudark/triper"
)

// worker contains the basic info to manage commands.
type worker struct {
	workerPool     chan chan triper.Command
	JobChannel     chan triper.Command
	CommandHandler triper.CommandHandlerRegister
}

// Bus stores the command handler.
type Bus struct {
	CommandHandler triper.CommandHandlerRegister
	maxworkers     int
	workerPool     chan chan triper.Command
}

// start initialize a worker ready to receive jobs.
func (w *worker) start() {
	go func() {
		for {
			w.workerPool <- w.JobChannel

			job := <-w.JobChannel
			handler, err := w.CommandHandler.GetHandler(job)
			if err != nil {
				continue
			}

			if !job.IsValid() {
				continue
			}

			if err = handler.Handle(job); err != nil {
				glog.Error(err)
			}
		}
	}()
}

// newWorker initialize the values of worker and start it.
func newWorker(commandHandler triper.CommandHandlerRegister, workerPool chan chan triper.Command) {
	w := worker{
		workerPool:     workerPool,
		CommandHandler: commandHandler,
		JobChannel:     make(chan triper.Command),
	}

	w.start()
}

// HandleCommand ad a job to the queue.
func (b *Bus) HandleCommand(command triper.Command) (id string) {
	// generate an unique identifier to trace the command
	command.GenerateUUID()
	go func(c triper.Command) {
		workerJobQueue := <-b.workerPool
		workerJobQueue <- c
	}(command)

	return command.GetID()
}

// NewBus return a bus with command handler register.
func NewBus(register triper.CommandHandlerRegister, maxworkers int) *Bus {
	b := &Bus{
		CommandHandler: register,
		maxworkers:     maxworkers,
		workerPool:     make(chan chan triper.Command),
	}

	// start the bus
	b.start()
	return b
}

// start the bus
func (b *Bus) start() {
	for i := 0; i < b.maxworkers; i++ {
		newWorker(b.CommandHandler, b.workerPool)
	}
}
