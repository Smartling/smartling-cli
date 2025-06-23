package mt

import "context"

// TranslateParams is the parameters for the RunTranslate method.
type TranslateParams struct {
	SourceLocale   string
	DetectLanguage string
	TargetLocale   string
	Directory      string
	Directive      string
	Progress       string
	FileType       string
	FormatPath     string
	FileOrPattern  string
}

func (s service) RunTranslate(ctx context.Context, p TranslateParams) (TranslateOutput, error) {
	return TranslateOutput{}, nil
}

type TranslateOutput struct {
	File      string
	Locale    string
	Name      string
	Ext       string
	Directory string
}
