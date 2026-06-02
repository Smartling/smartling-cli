package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	"github.com/Smartling/smartling-cli/services/jobs/jobresolver"

	smerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
)

const defaultFilesPageLimit = 500

// FilesParams carries the jobs-files request.
type FilesParams struct {
	ProjectUID   string
	JobUIDOrName string
	Limit        uint32
	Offset       uint32
}

// Validate checks required fields.
func (p FilesParams) Validate() error {
	if p.ProjectUID == "" {
		return smerror.ErrEmptyParam("ProjectUID")
	}
	if p.JobUIDOrName == "" {
		return smerror.ErrEmptyParam("JobUIDOrName")
	}
	return nil
}

// JobFileItem is a single source file row.
type JobFileItem struct {
	FileURI   string   `json:"fileUri"`
	LocaleIDs []string `json:"localeIds"`
}

// FilesOutput is the result of listing a job's files.
type FilesOutput struct {
	Files      []JobFileItem
	TotalCount int
	JSON       []byte
}

// JSONBytes returns the raw JSON payload.
func (o FilesOutput) JSONBytes() []byte { return o.JSON }

// SimpleLines returns a human-readable list.
func (o FilesOutput) SimpleLines() []string {
	if len(o.Files) == 0 {
		return []string{"No files found."}
	}
	lines := make([]string, 0, len(o.Files))
	for _, f := range o.Files {
		lines = append(lines, fmt.Sprintf("%s  %s", f.FileURI, strings.Join(f.LocaleIDs, ",")))
	}
	return lines
}

// TableData returns the files as a table.
func (o FilesOutput) TableData() ([]string, [][]string) {
	headers := []string{"FILE URI", "LOCALES"}
	rows := make([][]string, 0, len(o.Files))
	for _, f := range o.Files {
		rows = append(rows, []string{f.FileURI, strings.Join(f.LocaleIDs, ",")})
	}
	return headers, rows
}

// RunFiles resolves the job by UID or name and lists its source files.
func (s service) RunFiles(ctx context.Context, params FilesParams) (FilesOutput, error) {
	if err := params.Validate(); err != nil {
		return FilesOutput{}, fmt.Errorf("invalid files params: %w", err)
	}
	rlog.Debugf("running jobs files with params: %+v", params)

	jobUID, err := jobresolver.GetJobUID(ctx, s.job, params.ProjectUID, params.JobUIDOrName)
	if err != nil {
		return FilesOutput{}, fmt.Errorf("resolve job UID: %w", err)
	}

	limit := params.Limit
	if limit == 0 {
		limit = defaultFilesPageLimit
	}

	page, err := s.job.ListFiles(ctx, params.ProjectUID, jobUID, limit, params.Offset)
	if err != nil {
		return FilesOutput{}, fmt.Errorf("list files for job %q: %w", jobUID, err)
	}

	files := make([]JobFileItem, len(page.Items))
	for i, file := range page.Items {
		files[i] = JobFileItem{FileURI: file.FileURI, LocaleIDs: file.LocaleIDs}
	}
	res := FilesOutput{Files: files, TotalCount: page.TotalCount}
	b, err := json.Marshal(files)
	if err != nil {
		rlog.Errorf("failed to marshal job files to JSON: %v", err)
		return res, nil
	}
	res.JSON = b
	return res, nil
}
