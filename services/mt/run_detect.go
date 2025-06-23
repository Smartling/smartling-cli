package mt

import (
	"context"

	globfiles "github.com/Smartling/smartling-cli/services/helpers/glob_files"

	sdk "github.com/Smartling/api-sdk-go/api/mt"
)

// DetectParams is the parameters for the RunDetect method.
type DetectParams struct {
	FileType      string
	FormatPath    string
	FileOrPattern string
	ProjectID     string
	AccountUID    sdk.AccountUID
	URI           string
}

func (s service) RunDetect(ctx context.Context, p DetectParams, listAllFilesFn globfiles.ListFilesFn) ([]DetectOutput, error) {
	files, err := globfiles.Remote(listAllFilesFn, p.ProjectID, p.URI)
	if err != nil {
		return nil, err
	}

	var res []DetectOutput
	for _, file := range files {
		detectedLang, err := s.translationControl.DetectFileLanguage(p.AccountUID, sdk.FileUID(file.FileURI))
		if err != nil {
			return nil, err
		}
		res = append(res, DetectOutput{
			File:       file.FileURI,
			Language:   detectedLang.Code,
			Confidence: "",
		})
	}

	return res, nil
}

type DetectOutput struct {
	File       string
	Language   string
	Confidence string
}
