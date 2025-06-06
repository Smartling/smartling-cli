package progress

import (
	"fmt"
	"os"
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

func (progress *Progress) Flush() {
	err := progress.Renderer.Render(progress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to Render: "+err.Error())
	}
}
