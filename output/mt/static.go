package mt

import (
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/pointer"
	"github.com/Smartling/smartling-cli/services/mt"
)

// Static defines static output mode
type Static struct {
	model        Model
	OutputFormat OutputFormat
	dataProvider TableDataProvider
}

// Init inits Static
func (s *Static) Init(dataProvider TableDataProvider, files []string, targetLocalesQnt uint8, outputFormat, outputTemplate string) {
	s.OutputFormat = GetOutputFormat(outputFormat, outputTemplate)
	s.dataProvider = dataProvider

	s.model.Headers = dataProvider.Headers()
	s.model.RowByHeader = dataProvider.RowByHeaderName()

	rows := dataProvider.ToTableRows(files, targetLocalesQnt)

	s.model.Data = rows
}

// Run runs Static
func (s *Static) Run() error {
	return nil
}

// Update handle updates
func (s *Static) Update(updates chan any) error {
	for update := range updates {
		switch update := update.(type) {
		case mt.TranslateUpdates:
			rowByHeader := s.model.RowByHeader
			if row, found := rowByHeader["locale"]; found {
				s.model.Data[update.ID][row] = pointer.PNew(update.Locale)
			}
			if row, found := rowByHeader["upload"]; found {
				s.model.Data[update.ID][row] = done
			}
			if row, found := rowByHeader["translate"]; found && update.Translate != nil {
				s.model.Data[update.ID][row] = pointer.PNew(update.Translate)
			}
			if row, found := rowByHeader["translated_file"]; found && update.TranslatedFile != nil {
				s.model.Data[update.ID][row] = pointer.PNew(update.TranslatedFile)
			}
			if row, found := rowByHeader["download"]; found {
				s.model.Data[update.ID][row] = done
			}
		case mt.DetectUpdates:
			rowByHeader := s.model.RowByHeader
			if row, found := rowByHeader["language"]; found {
				s.model.Data[update.ID][row] = pointer.PNew(update.Language)
			}
			if row, found := rowByHeader["upload"]; found {
				s.model.Data[update.ID][row] = done
			}
			if row, found := rowByHeader["detect"]; found {
				s.model.Data[update.ID][row] = pointer.PNew(update.Detect)
			}
		case clierror.UIError:
			return update
		case error:
			return update
		}
	}
	return nil
}

// End ends static output
func (s *Static) End() {
	s.OutputFormat.FormatAndRender(s.model.Headers, s.model.Data)
}
