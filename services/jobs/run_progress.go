package jobs

import (
	"context"
	"errors"
	"fmt"

	api "github.com/Smartling/api-sdk-go/api/mt"
	"github.com/Smartling/smartling-cli/services/jobs/jobresolver"
)

var (
	errEmptyProjectUID   = errors.New("project UID is required")
	errEmptyJobUIDOrName = errors.New("job UID or job name is required")
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

	jobUID, err := jobresolver.GetJobUID(ctx, s.job, params.ProjectUID, params.JobUIDOrName)
	if err != nil {
		return ProgressOutput{}, fmt.Errorf("resolve job UID: %w", err)
	}

	progress, err := s.job.Progress(ctx, params.ProjectUID, jobUID)
	if err != nil {
		return ProgressOutput{}, fmt.Errorf("get job progress for %q: %w", jobUID, err)
	}

	return ProgressOutput{
		TranslationJobUID: jobUID,
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
