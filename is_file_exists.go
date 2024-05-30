package main

import "os"

func isFileExists(path string) bool {
	// we don't care about any other errors there, just return false if stat
	// failed for whatever reason
	_, err := os.Stat(path)
	return err == nil
}
