package reader

import (
	"io"
	"os"
	"strings"

	sdk "github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

// ReadFilesFromStdin reads file URIs from stdin, one per line.
// Returns a slice of sdk.File with file URIs, and an error if any.
func ReadFilesFromStdin() ([]sdk.File, error) {
	lines, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, hierr.Errorf(
			err,
			"unable to read stdin",
		)
	}

	var files []sdk.File

	for _, line := range strings.Split(string(lines), "\n") {
		if line == "" {
			continue
		}

		files = append(files, sdk.File{
			FileURI: line,
		})
	}

	return files, nil
}
