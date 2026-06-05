package jobsfiles

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Smartling/smartling-cli/services/jobs/jobresolver"

	smerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
)

// DefaultListPageLimit is the default page size for listing a job's files.
const DefaultListPageLimit = 500

// ListParams carries the jobs-files-list request.
type ListParams struct {
	ProjectID    string
	JobUIDOrName string
	Limit        uint32
	Offset       uint32
}

// Validate checks required fields.
func (p ListParams) Validate() error {
	if p.ProjectID == "" {
		return smerror.ErrEmptyParam("ProjectID")
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

// ListOutput is the result of listing a job's files.
type ListOutput struct {
	Files      []JobFileItem
	TotalCount int
	Offset     uint32
	JSON       []byte `json:"-"`
}

// truncated reports whether more files exist beyond the returned page.
func (o ListOutput) truncated() bool {
	return int(o.Offset)+len(o.Files) < o.TotalCount
}

// truncationNote describes the visible-vs-total page when truncated.
func (o ListOutput) truncationNote() string {
	return fmt.Sprintf(
		"Showing %d of %d files. Use --offset %d to see more.",
		len(o.Files), o.TotalCount, o.Offset+uint32(len(o.Files)),
	)
}

// JSONBytes returns the raw JSON payload.
func (o ListOutput) JSONBytes() []byte { return o.JSON }

// SimpleLines returns a human-readable list.
func (o ListOutput) SimpleLines() []string {
	if len(o.Files) == 0 {
		return []string{"No files found."}
	}
	lines := make([]string, 0, len(o.Files)+1)
	for _, f := range o.Files {
		lines = append(lines, fmt.Sprintf("%s  %s", f.FileURI, strings.Join(f.LocaleIDs, ",")))
	}
	if o.truncated() {
		lines = append(lines, o.truncationNote())
	}
	return lines
}

// TableData returns the files as a table.
func (o ListOutput) TableData() ([]string, [][]string) {
	headers := []string{"FILE URI", "LOCALES"}
	rows := make([][]string, 0, len(o.Files)+1)
	for _, f := range o.Files {
		rows = append(rows, []string{f.FileURI, strings.Join(f.LocaleIDs, ",")})
	}
	if o.truncated() {
		rows = append(rows, []string{o.truncationNote(), ""})
	}
	return headers, rows
}

// RunList resolves the job by UID or name and lists its source files.
func (s service) RunList(ctx context.Context, params ListParams) (ListOutput, error) {
	if err := params.Validate(); err != nil {
		return ListOutput{}, fmt.Errorf("invalid list params: %w", err)
	}

	jobUID, err := jobresolver.GetJobUID(ctx, s.job, params.ProjectID, params.JobUIDOrName)
	if err != nil {
		return ListOutput{}, fmt.Errorf("resolve job UID: %w", err)
	}

	limit := params.Limit
	if limit == 0 {
		limit = DefaultListPageLimit
	}

	page, err := s.jobFile.List(ctx, params.ProjectID, jobUID, limit, params.Offset)
	if err != nil {
		return ListOutput{}, fmt.Errorf("list files for job %q: %w", jobUID, err)
	}

	files := make([]JobFileItem, len(page.Items))
	for i, file := range page.Items {
		files[i] = JobFileItem{FileURI: file.FileURI, LocaleIDs: file.LocaleIDs}
	}
	res := ListOutput{Files: files, TotalCount: page.TotalCount, Offset: params.Offset}
	b, err := json.Marshal(filesJSON{
		Files:      files,
		TotalCount: page.TotalCount,
		Offset:     params.Offset,
	})
	if err != nil {
		return ListOutput{}, fmt.Errorf("marshal job files to JSON: %w", err)
	}
	res.JSON = b
	return res, nil
}

// filesJSON is the JSON shape for jobs-files-list output, carrying pagination
// metadata so consumers can detect truncated pages.
type filesJSON struct {
	Files      []JobFileItem `json:"files"`
	TotalCount int           `json:"totalCount"`
	Offset     uint32        `json:"offset"`
}
