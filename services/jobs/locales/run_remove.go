package locales

import (
	"context"

	"github.com/Smartling/smartling-cli/services/jobs/jobresolver"
)

// RemoveParams carries the remove-locale request from CLI to service.
type RemoveParams struct {
	ProjectID      string
	JobUIDOrName   string
	TargetLocaleID string
}

// Validate checks that RemoveParams carry the required fields.
func (p RemoveParams) Validate() error {
	return validateParams(p.ProjectID, p.JobUIDOrName, p.TargetLocaleID)
}

// RunRemove unassigns a target locale from a translation job.
func (s service) RunRemove(ctx context.Context, params RemoveParams) (Output, error) {
	if err := params.Validate(); err != nil {
		return Output{}, err
	}
	jobUID, err := jobresolver.GetJobUID(ctx, s.job, params.ProjectID, params.JobUIDOrName)
	if err != nil {
		return Output{}, err
	}
	if err := s.locale.Remove(ctx, params.ProjectID, jobUID, params.TargetLocaleID); err != nil {
		return Output{}, err
	}
	return newOutput("removed", params.ProjectID, jobUID, params.TargetLocaleID)
}
