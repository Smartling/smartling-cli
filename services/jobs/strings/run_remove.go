package jobstrings

import (
	"context"

	"github.com/Smartling/smartling-cli/services/jobs/jobresolver"

	api "github.com/Smartling/api-sdk-go/api/job/string"
)

// RemoveParams defines the remove-strings params.
type RemoveParams struct {
	ProjectID    string
	JobUIDOrName string
	Hashcodes    []string
	LocaleIDs    []string
}

// Validate checks that RemoveParams are valid.
func (p RemoveParams) Validate() error {
	return validateMutate(p.ProjectID, p.JobUIDOrName, p.Hashcodes)
}

// RunRemove unassigns strings from a translation job.
func (s service) RunRemove(ctx context.Context, params RemoveParams) (MutateOutput, error) {
	if err := params.Validate(); err != nil {
		return MutateOutput{}, err
	}
	jobUID, err := jobresolver.GetJobUID(ctx, s.job, params.ProjectID, params.JobUIDOrName)
	if err != nil {
		return MutateOutput{}, err
	}
	req := api.RemoveRequest{
		Hashcodes: params.Hashcodes,
		LocaleIDs: params.LocaleIDs,
	}
	res, err := s.jobString.Remove(ctx, params.ProjectID, jobUID, req)
	if err != nil {
		return MutateOutput{}, err
	}
	return newMutateOutput("removed", params.ProjectID, jobUID, params.Hashcodes, nil, params.LocaleIDs, res.SuccessCount, res.FailCount)
}
