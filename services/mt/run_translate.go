package mt

import (
	"context"

	sdk "github.com/Smartling/api-sdk-go/api/mt"
)

// TranslateParams is the parameters for the RunTranslate method.
type TranslateParams struct {
	SourceLocale   string
	DetectLanguage string
	TargetLocale   []string
	Directory      string
	Directives     map[string]string
	Progress       bool
	FileType       string
	OutputFormat   string
	FileOrPattern  string
	ProjectID      string
	AccountUID     sdk.AccountUID
	URI            string
}

func (s service) RunTranslate(ctx context.Context, p TranslateParams) (TranslateOutput, error) {
	return TranslateOutput{}, nil
}

// TranslateOutput is translate output
type TranslateOutput struct {
	File      string
	Locale    string
	Name      string
	Ext       string
	Directory string
}
