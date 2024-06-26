// Package internal contains internal logic of the application.
package internal

import "log"

// TODO: Think about this

// Task is a function that should be executed in Worker
type Task = func() error

// CreateRunTask creates Task for running tests
func CreateRunTask(run *Run) Task {
	return func() error {
		_, err := run.RunTests()
		return err
	}
}

// CreateRejudgeTask creates Task for rejudging tests
// FIXME: Implement me
func CreateRejudgeTask(runInfo *RunInfo) Task {
	panic("implement me")
}

// TestingQueue is a queue for running tasks
type TestingQueue struct {
	queue   chan Task
	workers chan Worker
}

// NewTestingQueue creates new pointer to TestingQueue
func NewTestingQueue(workers chan Worker) *TestingQueue {
	return &TestingQueue{
		// FIXME: Think about buffer size
		queue:   make(chan Task, 1000),
		workers: workers,
	}
}

// PushTask pushes task to the queue
func (tq *TestingQueue) PushTask(task Task) {
	tq.queue <- task
}

// Update updates the queue
func (tq *TestingQueue) Update() {
	for {
		select {
		case task := <-tq.queue:
			go func() {
				// Deadlock?
				worker := <-tq.workers
				defer func() { tq.workers <- worker }()
				err := worker.RunTask(task)
				if err != nil {
					log.Print(err)
				}
			}()
		}
	}
}

// CreateSimpleTestingQueue creates simple TestingQueue
func CreateSimpleTestingQueue() *TestingQueue {
	workers := make(chan Worker, 10)
	for i := 0; i < 10; i++ {
		worker := NewLocalWorker()
		workers <- worker
	}
	return NewTestingQueue(workers)
}

var MyTestingQueue = CreateSimpleTestingQueue()
