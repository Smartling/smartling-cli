package threadpool

import (
	"sync"
)

// ThreadPool is a thread pool implementation.
type ThreadPool struct {
	available chan struct{}
	size      uint32
	group     sync.WaitGroup
}

// NewThreadPool creates a new ThreadPool with the specified size.
func NewThreadPool(size uint32) *ThreadPool {
	available := make(chan struct{}, size)
	for i := uint32(0); i < size; i++ {
		available <- struct{}{}
	}

	return &ThreadPool{
		group:     sync.WaitGroup{},
		available: available,
		size:      size,
	}
}

// Do submits a task to the thread pool for execution.
func (pool *ThreadPool) Do(task func()) {
	<-pool.available

	pool.group.Add(1)

	go func() {
		defer func() {
			pool.group.Done()
			pool.available <- struct{}{}
		}()

		task()
	}()
}

// Wait waits for tasks in the thread pool to complete.
func (pool *ThreadPool) Wait() {
	pool.group.Wait()
}
