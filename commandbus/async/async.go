package async

import (
	"github.com/golang/glog"
	"github.com/mishudark/triper"
)

// Worker contains the basic info to manage commands.
type Worker struct {
	WorkerPool     chan chan triper.Command
	JobChannel     chan triper.Command
	CommandHandler triper.CommandHandlerRegister
}

// Bus stores the command handler.
type Bus struct {
	CommandHandler triper.CommandHandlerRegister
	maxWorkers     int
	workerPool     chan chan triper.Command
}

// Start initialize a worker ready to receive jobs.
func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel

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

// NewWorker initialize the values of worker and start it.
func NewWorker(commandHandler triper.CommandHandlerRegister, workerPool chan chan triper.Command) {
	w := Worker{
		WorkerPool:     workerPool,
		CommandHandler: commandHandler,
		JobChannel:     make(chan triper.Command),
	}

	w.Start()
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
func NewBus(register triper.CommandHandlerRegister, maxWorkers int) *Bus {
	b := &Bus{
		CommandHandler: register,
		maxWorkers:     maxWorkers,
		workerPool:     make(chan chan triper.Command),
	}

	// start the bus
	b.Start()
	return b
}

// Start the bus
func (b *Bus) Start() {
	for i := 0; i < b.maxWorkers; i++ {
		NewWorker(b.CommandHandler, b.workerPool)
	}
}
