package jobstrings

import (
	"context"

	"github.com/Smartling/smartling-cli/services/jobs/jobresolver"

	api "github.com/Smartling/api-sdk-go/api/job/string"
)

// AddParams define add string params.
type AddParams struct {
	ProjectID       string
	JobUIDOrName    string
	Hashcodes       []string
	TargetLocaleIDs []string
	MoveEnabled     bool
}

// Validate checks that AddParams are valid.
func (p AddParams) Validate() error {
	return validateMutate(p.ProjectID, p.JobUIDOrName, p.Hashcodes)
}

// RunAdd assigns strings to a translation job.
func (s service) RunAdd(ctx context.Context, params AddParams) (MutateOutput, error) {
	if err := params.Validate(); err != nil {
		return MutateOutput{}, err
	}
	jobUID, err := jobresolver.GetJobUID(ctx, s.job, params.ProjectID, params.JobUIDOrName)
	if err != nil {
		return MutateOutput{}, err
	}
	req := api.AddRequest{
		Hashcodes:       params.Hashcodes,
		TargetLocaleIDs: params.TargetLocaleIDs,
		MoveEnabled:     params.MoveEnabled,
	}
	res, err := s.jobString.Add(ctx, params.ProjectID, jobUID, req)
	if err != nil {
		return MutateOutput{}, err
	}
	return newMutateOutput("added", params.ProjectID, jobUID, params.Hashcodes, params.TargetLocaleIDs, res.SuccessCount, res.FailCount)
}
