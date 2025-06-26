package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/reconquest/hierr-go"
)

func getGitBranch() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", hierr.Errorf(
			err,
			"unable to get current working directory",
		)
	}

	for {
		if dir == "/" {
			return "", hierr.Errorf(
				err,
				"no git repository can be found containing current directory",
			)
		}

		_, err := os.Stat(filepath.Join(dir, ".git"))
		if err != nil {
			if !os.IsNotExist(err) {
				return "", hierr.Errorf(
					err,
					`unable to get stats for "%s"`,
					dir,
				)
			}

			dir = filepath.Dir(dir)

			continue
		} else {
			break
		}
	}

	head, err := os.ReadFile(filepath.Join(dir, ".git", "HEAD"))
	if err != nil {
		return "", hierr.Errorf(
			err,
			"unable to read git HEAD",
		)
	}

	return filepath.Base(strings.TrimSpace(string(head))), nil
}
