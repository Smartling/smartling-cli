package main

import (
	"fmt"
	"sync"
)

type Progress struct {
	sync.Mutex

	Current int
	Total   int

	Renderer ProgressRenderer
}

func (progress *Progress) String() string {
	return fmt.Sprintf("%d/%d", progress.Current, progress.Total)
}

func (progress *Progress) Increment() {
	progress.Lock()
	defer progress.Unlock()

	progress.Current++
}

func (progress Progress) Flush() {
	progress.Renderer.Render(progress)
}
