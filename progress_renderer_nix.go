//go:build linux || darwin || freebsd || netbsd || openbsd || dragonfly

package main

import (
	"fmt"
	"os"
)

type ProgressRenderer struct{}

func (renderer ProgressRenderer) Render(progress *Progress) error {
	_, err := fmt.Fprintf(os.Stderr, "%s\r", progress.String())

	return err
}
