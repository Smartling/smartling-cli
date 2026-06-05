package jobsfiles

import (
	"context"

	"github.com/Smartling/smartling-cli/services/jobs/jobresolver"

	api "github.com/Smartling/api-sdk-go/api/job/file"
)

// RemoveParams carries the remove-files request from CLI to service.
type RemoveParams struct {
	ProjectID    string
	JobUIDOrName string
	FilePatterns []string
}

// Validate checks that RemoveParams carry the required fields.
func (p RemoveParams) Validate() error {
	return validateMutate(p.ProjectID, p.JobUIDOrName, p.FilePatterns)
}

// RunRemove detaches files matching the given patterns from a translation job,
// one API call per resolved fileUri.
func (s service) RunRemove(ctx context.Context, params RemoveParams) (MutateOutput, error) {
	if err := params.Validate(); err != nil {
		return MutateOutput{}, err
	}
	jobUID, err := jobresolver.GetJobUID(ctx, s.job, params.ProjectID, params.JobUIDOrName)
	if err != nil {
		return MutateOutput{}, err
	}
	uris, unmatched, err := s.resolveURIs(ctx, params.ProjectID, params.FilePatterns)
	if err != nil {
		return MutateOutput{}, err
	}

	var (
		files                   []FileResult
		totalSuccess, totalFail int
	)
	for _, uri := range uris {
		res, removeErr := s.jobFile.Remove(ctx, params.ProjectID, jobUID, api.RemoveRequest{FileURI: uri})
		if removeErr != nil {
			files = append(files, FileResult{FileURI: uri, Error: removeErr.Error()})
			continue
		}
		totalSuccess += res.SuccessCount
		totalFail += res.FailCount
		files = append(files, FileResult{FileURI: uri, SuccessCount: res.SuccessCount, FailCount: res.FailCount})
	}

	return newMutateOutput("removed", params.ProjectID, jobUID, nil, files, unmatched, totalSuccess, totalFail)
}
