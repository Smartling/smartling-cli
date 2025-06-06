package threadpool

import (
	"sync"
)

type ThreadPool struct {
	available chan struct{}
	size      int
	group     sync.WaitGroup
}

func NewThreadPool(size int) *ThreadPool {
	available := make(chan struct{}, size)
	for i := 0; i < size; i++ {
		available <- struct{}{}
	}

	return &ThreadPool{
		group:     sync.WaitGroup{},
		available: available,
		size:      size,
	}
}

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

func (pool *ThreadPool) Wait() {
	pool.group.Wait()
}
