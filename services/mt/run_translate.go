package mt

import (
	"context"

	sdk "github.com/Smartling/api-sdk-go/api/mt"
)

// TranslateParams is the parameters for the RunTranslate method.
type TranslateParams struct {
	SourceLocale     string
	DetectLanguage   bool
	TargetLocales    []string
	OutputDirectory  string
	Directives       map[string]string
	Progress         bool
	OverrideFileType string
	FileOrPattern    string
	ProjectID        string
	AccountUID       sdk.AccountUID
	URI              string
}

func (s service) RunTranslate(ctx context.Context, p TranslateParams) ([]TranslateOutput, error) {
	return nil, nil
}

// TranslateOutput is translate output
type TranslateOutput struct {
	File      string
	Locale    string
	Name      string
	Ext       string
	Directory string
}
