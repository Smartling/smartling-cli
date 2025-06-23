package mt

import (
	"context"

	sdkfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
)

// DetectParams is the parameters for the RunDetect method.
type DetectParams struct {
	FileType      string
	FormatPath    string
	FileOrPattern string
}

func (s service) RunDetect(ctx context.Context, p DetectParams) (DetectOutput, error) {

	s.translationControl

	s.translationControl.DetectionProgress()

	sdkfile.File{}
	return DetectOutput{}, nil
}

type DetectOutput struct {
	File       string
	Language   string
	Confidence string
}
