package mt

import (
	"github.com/Smartling/smartling-cli/services/helpers/pointer"
	"github.com/Smartling/smartling-cli/services/mt"
)

type Static struct {
	model        Model
	OutputFormat OutputFormat
	dataProvider TableDataProvider
}

func (s *Static) Init(dataProvider TableDataProvider, files []string, outputFormat, outputTemplate string) {
	s.model.OutputFormat = outputFormat
	s.model.OutputTemplate = outputTemplate
	s.OutputFormat = GetOutputFormat(s.model.OutputFormat, s.model.OutputTemplate)
	s.dataProvider = dataProvider

	s.model.Headers = dataProvider.Headers()
	s.model.RowByHeader = dataProvider.RowByHeaderName()

	rows := dataProvider.ToTableRows(files)
	//dataProvider.SetRows(rows)

	s.model.Data = rows
}

func (s *Static) Run() error {
	return nil
}

func (s *Static) Update(updates chan any) {
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
		}
	}
}

func (s *Static) End() {
	s.OutputFormat.FormatAndRender(s.model.Headers, s.model.Data)
}
