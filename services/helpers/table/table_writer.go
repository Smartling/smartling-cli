package table

import (
	"io"
	"text/tabwriter"

	"github.com/reconquest/hierr-go"
)

func NewTableWriter(target io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(target, 2, 4, 2, ' ', 0)
}

func Render(writer *tabwriter.Writer) error {
	err := writer.Flush()
	if err != nil {
		return hierr.Errorf(
			err,
			"unable to flush table to stdout",
		)
	}

	return nil
}
