package mt

import (
	"fmt"
	"io"
	"os"

	"github.com/Smartling/smartling-cli/services/helpers/format"
	"github.com/Smartling/smartling-cli/services/helpers/table"
	"github.com/Smartling/smartling-cli/services/mt"
)

const DefaultDetectTemplate = `{{.File}}\t{{.Language}}\n`

// RenderDetect render detect output
func RenderDetect(detectOutputs []mt.DetectOutput, outputFormat, outputTemplate string) error {
	if outputTemplate == "" {
		outputTemplate = DefaultDetectTemplate
	}

	format, err := format.Compile(outputTemplate)
	if err != nil {
		return err
	}

	tableWriter := table.NewTableWriter(os.Stdout)

	for _, file := range detectOutputs {
		row, err := format.Execute(file)
		if err != nil {
			return err
		}

		_, err = io.WriteString(tableWriter, row)
		if err != nil {
			return fmt.Errorf("unable to write row to output table: %w", err)
		}
	}

	return table.Render(tableWriter)
}
