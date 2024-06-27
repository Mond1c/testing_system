// Package internal contains internal logic of the application.
package pkg

// Worker is an interface that should be implemented by all workers.
type Worker interface {
	RunTask(task Task) error
}

// LocalWorker is a worker that runs tasks locally.
type LocalWorker struct {
}

// NewLocalWorker creates new LocalWorker.
func NewLocalWorker() *LocalWorker {
	return &LocalWorker{}
}

// RunTask runs task locally.
func (lw *LocalWorker) RunTask(task Task) error {
	return task.Local()
}

// TODO: Add remote worker and create package with remote worker

type RemoteWorker struct {
	address string
}

func NewRemoteWorker(address string) *RemoteWorker {
	return &RemoteWorker{address: address}
}

func (rw *RemoteWorker) RunTask(task Task) error {
	return task.Remote(rw.address)
}
