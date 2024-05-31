package main

import (
	"io"
	"os"
	"strings"

	smartling "github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

func readFilesFromStdin() ([]smartling.File, error) {
	lines, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, hierr.Errorf(
			err,
			"unable to read stdin",
		)
	}

	var files []smartling.File

	for _, line := range strings.Split(string(lines), "\n") {
		if line == "" {
			continue
		}

		files = append(files, smartling.File{
			FileURI: line,
		})
	}

	return files, nil
}
