package jobsfiles

import (
	"context"

	"github.com/Smartling/smartling-cli/services/jobs/jobresolver"

	api "github.com/Smartling/api-sdk-go/api/job/file"
)

// AddParams carries the add-files request from CLI to service.
type AddParams struct {
	ProjectID       string
	JobUIDOrName    string
	FilePatterns    []string
	TargetLocaleIDs []string
}

// Validate checks that AddParams are valid.
func (p AddParams) Validate() error {
	return validateMutate(p.ProjectID, p.JobUIDOrName, p.FilePatterns)
}

// RunAdd attaches files matching the given patterns to a translation job, one
// API call per resolved fileUri.
func (s service) RunAdd(ctx context.Context, params AddParams) (MutateOutput, error) {
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
		res, addErr := s.jobFile.Add(ctx, params.ProjectID, jobUID, api.AddRequest{
			FileURI:         uri,
			TargetLocaleIDs: params.TargetLocaleIDs,
		})
		if addErr != nil {
			files = append(files, FileResult{FileURI: uri, Error: addErr.Error()})
			continue
		}
		totalSuccess += res.SuccessCount
		totalFail += res.FailCount
		files = append(files, FileResult{FileURI: uri, SuccessCount: res.SuccessCount, FailCount: res.FailCount})
	}

	return newMutateOutput("added", params.ProjectID, jobUID, params.TargetLocaleIDs, files, unmatched, totalSuccess, totalFail)
}
