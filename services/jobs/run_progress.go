package jobs

import (
	"context"
	"regexp"

	"github.com/Smartling/api-sdk-go/api/job"
	api "github.com/Smartling/api-sdk-go/api/mt"
)

// ProgressParams is the parameters for the RunProgress method.
type ProgressParams struct {
	AccountUID  api.AccountUID
	ProjectUID  string
	JobIDOrName string
}

func (p ProgressParams) Validate() error {
	if err := p.AccountUID.Validate(); err != nil {
		return err
	}
	return nil
}

func (s service) RunProgress(ctx context.Context, params ProgressParams) (ProgressOutput, error) {
	if err := params.Validate(); err != nil {
		return ProgressOutput{}, err
	}

	pattern := `^[a-z0-9]{12}$`
	var translationJobUID string
	if re := regexp.MustCompile(pattern); params.JobIDOrName != "" && re.MatchString(params.JobIDOrName) {
		if jb, err := s.job.Get(params.ProjectUID, params.JobIDOrName); err == nil {
			translationJobUID = jb.TranslationJobUID
		}
	}
	if translationJobUID == "" {
		jobs, err := s.job.GetAllByName(params.ProjectUID, params.JobIDOrName)
		if err != nil {
			return ProgressOutput{}, err
		}
		if len(jobs) == 0 {
			return ProgressOutput{}, nil
		}
		job := job.FindFirstJobByName(jobs, params.JobIDOrName)
		translationJobUID = job.TranslationJobUID
	}
	if translationJobUID == "" {
		return ProgressOutput{}, nil
	}

	progress, err := s.job.Progress(params.ProjectUID, translationJobUID)
	if err != nil {
		return ProgressOutput{}, err
	}

	return ProgressOutput{
		TranslationJobUID: translationJobUID,
		TotalWordCount:    progress.TotalWordCount,
		PercentComplete:   progress.PercentComplete,
		Json:              progress.Json,
	}, nil
}

// ProgressOutput represents the result of a job progress
type ProgressOutput struct {
	TranslationJobUID string
	TotalWordCount    uint32
	PercentComplete   uint32
	Json              []byte
}

// DetectUpdates defines updates
type DetectUpdates struct {
	ID       uint32
	Language *string
	Upload   *bool
	Detect   *string
}
