package jobs

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/Smartling/api-sdk-go/api/job"
	api "github.com/Smartling/api-sdk-go/api/mt"
)

var (
	// ErrJobNotFound is returned when a job is not found.
	ErrJobNotFound       = errors.New("job not found")
	errEmptyProjectUID   = errors.New("project UID is required")
	errEmptyJobUIDOrName = errors.New("job UID or job name is required")
	jobUIDPattern        = regexp.MustCompile(`^[a-z0-9]{12}$`)
)

// ProgressParams is the parameters for the RunProgress method.
type ProgressParams struct {
	AccountUID   api.AccountUID
	ProjectUID   string
	JobUIDOrName string
}

// Validate validates params for RunProgress.
func (p ProgressParams) Validate() error {
	if err := p.AccountUID.Validate(); err != nil {
		return err
	}
	if p.ProjectUID == "" {
		return errEmptyProjectUID
	}
	if p.JobUIDOrName == "" {
		return errEmptyJobUIDOrName
	}
	return nil
}

// RunProgress resolves the job by UID or name and returns its translation progress.
func (s service) RunProgress(ctx context.Context, params ProgressParams) (ProgressOutput, error) {
	if err := params.Validate(); err != nil {
		return ProgressOutput{}, err
	}

	var translationJobUID string
	if params.JobUIDOrName != "" && jobUIDPattern.MatchString(params.JobUIDOrName) {
		jb, err := s.job.GetJob(ctx, params.ProjectUID, params.JobUIDOrName)
		switch {
		case err == nil:
			translationJobUID = jb.TranslationJobUID
		case errors.Is(err, job.ErrNotFound):
			// 12-char input wasn't a UID — could still be a job name, fall through
		default:
			return ProgressOutput{}, fmt.Errorf("get job by UID %q: %w", params.JobUIDOrName, err)
		}
	}

	if translationJobUID == "" {
		jobs, err := s.job.SearchByName(ctx, params.ProjectUID, params.JobUIDOrName)
		if err != nil {
			return ProgressOutput{}, fmt.Errorf("search jobs by name %q: %w", params.JobUIDOrName, err)
		}
		if len(jobs) == 0 {
			return ProgressOutput{}, ErrJobNotFound
		}
		j, found := job.FindFirstJobByName(jobs, params.JobUIDOrName)
		if !found {
			return ProgressOutput{}, ErrJobNotFound
		}
		translationJobUID = j.TranslationJobUID
	}

	progress, err := s.job.Progress(ctx, params.ProjectUID, translationJobUID)
	if err != nil {
		return ProgressOutput{}, fmt.Errorf("get job progress for %q: %w", translationJobUID, err)
	}

	return ProgressOutput{
		TranslationJobUID: translationJobUID,
		TotalWordCount:    progress.TotalWordCount,
		PercentComplete:   progress.PercentComplete,
		JSON:              progress.JSON,
	}, nil
}

// ProgressOutput represents the result of a job progress
type ProgressOutput struct {
	TranslationJobUID string
	TotalWordCount    uint32
	PercentComplete   float64
	JSON              []byte
}
