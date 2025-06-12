package projects

import (
	"fmt"
	"io"
	"os"

	"github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/format"
	"github.com/Smartling/smartling-cli/services/helpers/table"

	sdk "github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

// LocalesParams is the parameters for the RunLocales method.
type LocalesParams struct {
	Format string
	Short  bool
	Source bool
}

// RunLocales retrieves and outputs the locales.
func (s service) RunLocales(params LocalesParams) error {
	formatType := params.Format
	if formatType == "" {
		formatType = format.DefaultProjectsLocalesFormat
	}

	format, err := format.Compile(formatType)
	if err != nil {
		return err
	}

	details, err := s.Client.GetProjectDetails(s.Config.ProjectID)
	if err != nil {
		if _, ok := err.(sdk.NotFoundError); ok {
			return clierror.ProjectNotFoundError{}
		}

		return hierr.Errorf(
			err,
			`unable to get project "%s" details`,
			s.Config.ProjectID,
		)
	}

	tableWriter := table.NewTableWriter(os.Stdout)

	if params.Source {
		if params.Short {
			fmt.Fprintf(tableWriter, "%s\n", details.SourceLocaleID)
		} else {
			fmt.Fprintf(
				tableWriter,
				"%s\t%s\n",
				details.SourceLocaleID,
				details.SourceLocaleDescription,
			)
		}
	} else {
		for _, locale := range details.TargetLocales {
			if params.Short {
				fmt.Fprintf(tableWriter, "%s\n", locale.LocaleID)
			} else {
				row, err := format.Execute(locale)
				if err != nil {
					return err
				}

				_, err = io.WriteString(tableWriter, row)
				if err != nil {
					return hierr.Errorf(
						err,
						"unable to write row to output table",
					)
				}
			}
		}
	}

	err = table.Render(tableWriter)
	if err != nil {
		return err
	}

	return nil
}
