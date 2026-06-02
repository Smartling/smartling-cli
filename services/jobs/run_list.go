package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Smartling/smartling-cli/services/helpers"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	smerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
	"github.com/Smartling/api-sdk-go/helpers/uid"
)

// ListParams carries the jobs-list request from CLI to service.
type ListParams struct {
	AccountUID         uid.AccountUID
	ProjectUID         string
	Account            bool
	JobName            string
	JobNumber          string
	JobStatus          []string
	TranslationJobUIDs []string
	ProjectIDs         []string
	FileURIs           []string
	Hashcodes          []string
	WithPriority       bool
	Limit              uint32
	Offset             uint32
	SortBy             string
	SortDirection      string
}

// Validate checks the fields required for the chosen scope.
func (p ListParams) Validate() error {
	if p.Account {
		return p.AccountUID.Validate()
	}
	if p.ProjectUID == "" {
		return smerror.ErrEmptyParam("ProjectUID")
	}
	return nil
}

// searchScope reports whether the file/hashcode search endpoint should be used.
func (p ListParams) searchScope() bool {
	return len(p.FileURIs) > 0 || len(p.Hashcodes) > 0
}

// JobListItem is a single row in a jobs list.
type JobListItem struct {
	TranslationJobUID string   `json:"translationJobUid"`
	JobName           string   `json:"jobName"`
	JobNumber         string   `json:"jobNumber"`
	JobStatus         string   `json:"jobStatus"`
	DueDate           string   `json:"dueDate"`
	TargetLocaleIDs   []string `json:"targetLocaleIds"`
	ProjectID         string   `json:"projectId,omitempty"`
	Priority          int      `json:"priority,omitempty"`
}

// ListOutput is the result of a jobs list.
type ListOutput struct {
	Jobs    []JobListItem
	Account bool
	JSON    []byte
}

// JSONBytes returns the raw JSON payload of the list.
func (o ListOutput) JSONBytes() []byte { return o.JSON }

// SimpleLines returns a human-readable summary.
func (o ListOutput) SimpleLines() []string {
	if len(o.Jobs) == 0 {
		return []string{"No jobs found."}
	}
	lines := make([]string, 0, len(o.Jobs))
	for _, j := range o.Jobs {
		lines = append(lines, fmt.Sprintf("%s  %s  %s", j.TranslationJobUID, j.JobName, j.JobStatus))
	}
	return lines
}

// TableData returns the list as a table.
func (o ListOutput) TableData() ([]string, [][]string) {
	if o.Account {
		headers := []string{"TRANSLATION JOB UID", "NAME", "STATUS", "DUE DATE", "PROJECT ID", "PRIORITY"}
		rows := make([][]string, 0, len(o.Jobs))
		for _, j := range o.Jobs {
			rows = append(rows, []string{j.TranslationJobUID, j.JobName, j.JobStatus, j.DueDate, j.ProjectID, fmt.Sprintf("%d", j.Priority)})
		}
		return headers, rows
	}
	headers := []string{"TRANSLATION JOB UID", "NAME", "NUMBER", "STATUS", "DUE DATE", "LOCALES"}
	rows := make([][]string, 0, len(o.Jobs))
	for _, j := range o.Jobs {
		rows = append(rows, []string{j.TranslationJobUID, j.JobName, j.JobNumber, j.JobStatus, j.DueDate, strings.Join(j.TargetLocaleIDs, ",")})
	}
	return headers, rows
}

// RunList lists jobs by project, account, or file/hashcode search.
func (s service) RunList(ctx context.Context, params ListParams) (ListOutput, error) {
	if err := params.Validate(); err != nil {
		return ListOutput{}, fmt.Errorf("invalid list params: %w", err)
	}
	rlog.Debugf("running jobs list with params: %+v", params)

	var (
		resp jobapi.ListJobsResponse
		err  error
	)
	switch {
	case params.searchScope():
		resp, err = s.job.SearchJobs(ctx, params.ProjectUID, jobapi.SearchJobsRequest{
			FileURIs:           params.FileURIs,
			Hashcodes:          params.Hashcodes,
			TranslationJobUIDs: params.TranslationJobUIDs,
		})
	case params.Account:
		resp, err = s.job.ListAccountJobs(ctx, string(params.AccountUID), jobapi.ListAccountJobsParams{
			JobName:       params.JobName,
			ProjectIDs:    params.ProjectIDs,
			JobStatus:     params.JobStatus,
			WithPriority:  params.WithPriority,
			Limit:         params.Limit,
			Offset:        params.Offset,
			SortBy:        params.SortBy,
			SortDirection: params.SortDirection,
		})
	default:
		resp, err = s.job.ListProjectJobs(ctx, params.ProjectUID, jobapi.ListProjectJobsParams{
			JobName:            params.JobName,
			JobNumber:          params.JobNumber,
			TranslationJobUIDs: params.TranslationJobUIDs,
			JobStatus:          params.JobStatus,
			Limit:              params.Limit,
			Offset:             params.Offset,
			SortBy:             params.SortBy,
			SortDirection:      params.SortDirection,
		})
	}
	if err != nil {
		return ListOutput{}, fmt.Errorf("failed to list jobs: %w", err)
	}

	return toListOutput(resp, params.Account), nil
}

func toListOutput(resp jobapi.ListJobsResponse, account bool) ListOutput {
	items := make([]JobListItem, 0, len(resp.Items))
	for _, j := range resp.Items {
		items = append(items, JobListItem{
			TranslationJobUID: j.TranslationJobUID,
			JobName:           j.JobName,
			JobNumber:         j.JobNumber,
			JobStatus:         j.JobStatus,
			DueDate:           helpers.TimeToString(j.Dates.Due, time.RFC3339),
			TargetLocaleIDs:   j.TargetLocaleIDs,
			ProjectID:         j.ProjectID,
			Priority:          j.Priority,
		})
	}
	out := ListOutput{Jobs: items, Account: account}
	b, err := json.Marshal(items)
	if err != nil {
		rlog.Errorf("failed to marshal jobs list to JSON: %v", err)
		return out
	}
	out.JSON = b
	return out
}
