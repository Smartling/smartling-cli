package locales

import (
	"context"

	"github.com/Smartling/smartling-cli/services/jobs/jobresolver"
)

// AddParams carries the add-locale request from CLI to service.
type AddParams struct {
	ProjectID      string
	JobUIDOrName   string
	TargetLocaleID string
}

// Validate checks that AddParams carry the required fields.
func (p AddParams) Validate() error {
	return validateParams(p.ProjectID, p.JobUIDOrName, p.TargetLocaleID)
}

// RunAdd assigns a target locale to a translation job.
func (s service) RunAdd(ctx context.Context, params AddParams) (Output, error) {
	if err := params.Validate(); err != nil {
		return Output{}, err
	}
	jobUID, err := jobresolver.GetJobUID(ctx, s.job, params.ProjectID, params.JobUIDOrName)
	if err != nil {
		return Output{}, err
	}
	if err := s.locale.Add(ctx, params.ProjectID, jobUID, params.TargetLocaleID); err != nil {
		return Output{}, err
	}
	return newOutput("added", params.ProjectID, jobUID, params.TargetLocaleID)
}
