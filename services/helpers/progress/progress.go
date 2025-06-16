package progress

import (
	"fmt"
	"os"
	"sync"
)

// Progress is progress tracking structure
type Progress struct {
	sync.Mutex

	Current int
	Total   int

	Renderer Renderer
}

// String returns string representation of the progress
func (progress *Progress) String() string {
	return fmt.Sprintf("%d/%d", progress.Current, progress.Total)
}

// Increment increments the current progress by one.
func (progress *Progress) Increment() {
	progress.Lock()
	defer progress.Unlock()

	progress.Current++
}

// Flush renders the current progress using the specified renderer.
func (progress *Progress) Flush() {
	err := progress.Renderer.Render(progress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to Render: %s", err.Error())
	}
}
