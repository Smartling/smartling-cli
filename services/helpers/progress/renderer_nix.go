//go:build linux || darwin || freebsd || netbsd || openbsd || dragonfly

package progress

import (
	"fmt"
	"os"
)

// Renderer is a renderer for progress that outputs to standard error.
type Renderer struct{}

// Render outputs the progress to standard error.
func (r Renderer) Render(progress *Progress) error {
	_, err := fmt.Fprintf(os.Stderr, "%s\r", progress.String())

	return err
}
