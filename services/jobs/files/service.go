package jobsfiles

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	globfiles "github.com/Smartling/smartling-cli/services/helpers/glob_files"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	api "github.com/Smartling/api-sdk-go/api/job/file"
	sdkfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
	"github.com/gobwas/glob"
)

// Service defines behavior for managing a translation job's files.
type Service interface {
	RunList(ctx context.Context, params ListParams) (ListOutput, error)
	RunAdd(ctx context.Context, params AddParams) (MutateOutput, error)
	RunRemove(ctx context.Context, params RemoveParams) (MutateOutput, error)
}

// NewService creates a new implementation of the Service. jobFile performs the
// add/remove; job resolves a UID from a UID-or-name; listFiles expands --file
// glob patterns against the project's files.
func NewService(jobFile api.JobFile, job jobapi.Job, listFiles globfiles.ListFilesFn) Service {
	return service{
		jobFile:   jobFile,
		job:       job,
		listFiles: listFiles,
	}
}

type service struct {
	jobFile   api.JobFile
	job       jobapi.Job
	listFiles globfiles.ListFilesFn
}

// FileResult is the per-file outcome of an add/remove operation.
type FileResult struct {
	FileURI      string `json:"fileUri"`
	SuccessCount int    `json:"successCount"`
	FailCount    int    `json:"failCount"`
	Error        string `json:"error,omitempty"`
}

// MutateOutput is the aggregated result of an add/remove files operation.
type MutateOutput struct {
	Action            string       `json:"action"`
	ProjectUID        string       `json:"projectUid"`
	TranslationJobUID string       `json:"translationJobUid"`
	TargetLocaleIDs   []string     `json:"targetLocaleIds,omitempty"`
	Files             []FileResult `json:"files"`
	Unmatched         []string     `json:"unmatched,omitempty"`
	SuccessCount      int          `json:"successCount"`
	FailCount         int          `json:"failCount"`

	JSON []byte `json:"-"`
}

func newMutateOutput(action, projectUID, jobUID string, targetLocaleIDs []string, files []FileResult, unmatched []string, successCount, failCount int) (MutateOutput, error) {
	o := MutateOutput{
		Action:            action,
		ProjectUID:        projectUID,
		TranslationJobUID: jobUID,
		TargetLocaleIDs:   targetLocaleIDs,
		Files:             files,
		Unmatched:         unmatched,
		SuccessCount:      successCount,
		FailCount:         failCount,
	}
	var err error
	if o.JSON, err = json.Marshal(o); err != nil {
		return MutateOutput{}, err
	}
	return o, nil
}

// JSONBytes returns the JSON representation of the result.
func (o MutateOutput) JSONBytes() []byte { return o.JSON }

// FailedFileURIs returns the URIs whose per-file API call errored.
func (o MutateOutput) FailedFileURIs() []string {
	var failed []string
	for _, f := range o.Files {
		if f.Error != "" {
			failed = append(failed, f.FileURI)
		}
	}
	return failed
}

// SimpleLines returns a human-readable summary of the result.
func (o MutateOutput) SimpleLines() []string {
	lines := []string{fmt.Sprintf("Files %s — job %s: %d succeeded, %d failed across %d file(s)",
		o.Action, o.TranslationJobUID, o.SuccessCount, o.FailCount, len(o.Files))}
	for _, f := range o.Files {
		if f.Error != "" {
			lines = append(lines, fmt.Sprintf("  %s: error: %s", f.FileURI, f.Error))
			continue
		}
		lines = append(lines, fmt.Sprintf("  %s: %d succeeded, %d failed", f.FileURI, f.SuccessCount, f.FailCount))
	}
	for _, pattern := range o.Unmatched {
		lines = append(lines, fmt.Sprintf("No files matched: %s", pattern))
	}
	return lines
}

// TableData returns one row per affected file.
func (o MutateOutput) TableData() ([]string, [][]string) {
	headers := []string{"FILE URI", "SUCCEEDED", "FAILED", "ERROR"}
	rows := make([][]string, 0, len(o.Files))
	for _, f := range o.Files {
		rows = append(rows, []string{f.FileURI, strconv.Itoa(f.SuccessCount), strconv.Itoa(f.FailCount), f.Error})
	}
	return headers, rows
}

// resolveURIs expands the glob patterns against the project's files, returning
// matched URIs (deduped, in pattern order) and patterns that matched nothing.
// It lists the project's files once and matches with the same glob engine as
// `files list`.
func (s service) resolveURIs(ctx context.Context, projectID string, patterns []string) (uris, unmatched []string, err error) {
	files, err := s.listFiles(ctx, projectID, sdkfile.FilesListRequest{})
	if err != nil {
		return nil, nil, fmt.Errorf("unable to list project files: %w", err)
	}

	seen := map[string]bool{}
	for _, pattern := range patterns {
		g, err := glob.Compile(pattern, '/')
		if err != nil {
			return nil, nil, fmt.Errorf("invalid --file pattern %q: %w", pattern, err)
		}
		matched := false
		for _, f := range files {
			if !g.Match(f.FileURI) {
				continue
			}
			matched = true
			if !seen[f.FileURI] {
				seen[f.FileURI] = true
				uris = append(uris, f.FileURI)
			}
		}
		if !matched {
			unmatched = append(unmatched, pattern)
		}
	}
	return uris, unmatched, nil
}

func validateMutate(projectID, jobUIDOrName string, patterns []string) error {
	switch {
	case projectID == "":
		return errors.New("project ID is required")
	case jobUIDOrName == "":
		return errors.New("translation job UID or name is required")
	case len(patterns) == 0:
		return errors.New("at least one --file is required")
	}
	return nil
}
