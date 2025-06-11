//go:build linux || darwin || freebsd || netbsd || openbsd || dragonfly

package progress

import (
	"fmt"
	"os"
)

// ProgressRenderer is a renderer for progress that outputs to standard error.
type ProgressRenderer struct{}

// Render outputs the progress to standard error.
func (renderer ProgressRenderer) Render(progress *Progress) error {
	_, err := fmt.Fprintf(os.Stderr, "%s\r", progress.String())

	return err
}
