package queue

import (
	"log"
	"sync"
)

// Worker - Worker that procresses tasks
type Worker struct {

	// Debug Settings
	wLog   *log.Logger // logging
	wDebug bool        // enable debugging

	// task Channals
	wReadyChan    chan chan Task
	wAssignedTask chan Task

	// worker synchronization
	wIsDone sync.WaitGroup

	// worker quit
	wQuit chan bool
}

// NewWorker - Creates a new worker
func NewWorker(readyPool chan chan Task, done sync.WaitGroup) *Worker {
	return &Worker{
		wReadyChan:    readyPool,
		wAssignedTask: make(chan Task),
		wIsDone:       done,
		wQuit:         make(chan bool),
	}
}

// Start - Begins processing worker's task
func (w *Worker) Start() {
	go func() {
		w.wIsDone.Add(1)
		for {
			w.wReadyChan <- w.wAssignedTask // check the Task queue in
			select {
			case Task := <-w.wAssignedTask: // see if anything has been assigned to the queue
				Task.Run()
			case <-w.wQuit:
				w.wIsDone.Done()
				return
			}
		}
	}()
}

// Stop - Stops the worker
func (w *Worker) Stop() {
	w.wQuit <- true
}
