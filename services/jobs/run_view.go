package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Smartling/smartling-cli/services/helpers"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	"github.com/Smartling/smartling-cli/services/jobs/jobresolver"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	smerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
)

// ViewParams carries the jobs-view request.
type ViewParams struct {
	ProjectUID   string
	JobUIDOrName string
}

// Validate checks required fields.
func (p ViewParams) Validate() error {
	if p.ProjectUID == "" {
		return smerror.ErrEmptyParam("ProjectUID")
	}
	if p.JobUIDOrName == "" {
		return smerror.ErrEmptyParam("JobUIDOrName")
	}
	return nil
}

// ViewOutput is the full detail of a single job.
type ViewOutput struct {
	jobapi.GetJobResponse
	JSON []byte
}

// JSONBytes returns the raw JSON payload.
func (o ViewOutput) JSONBytes() []byte { return o.JSON }

// SimpleLines returns a human-readable detail block.
func (o ViewOutput) SimpleLines() []string {
	return []string{
		fmt.Sprintf("Translation job UID: %s", o.TranslationJobUID),
		fmt.Sprintf("Name:                %s", o.JobName),
		fmt.Sprintf("Number:              %s", o.JobNumber),
		fmt.Sprintf("Status:              %s", o.JobStatus),
		fmt.Sprintf("Description:         %s", o.Description),
		fmt.Sprintf("Reference number:    %s", o.ReferenceNumber),
		fmt.Sprintf("Due date:            %s", helpers.TimeToString(o.Dates.Due, time.RFC3339)),
		fmt.Sprintf("Modified date:       %s", helpers.TimeToString(o.Dates.Modified, time.RFC3339)),
		fmt.Sprintf("Created date:        %s", helpers.TimeToString(o.Dates.Created, time.RFC3339)),
		fmt.Sprintf("Priority:            %d", o.Priority),
		fmt.Sprintf("Target locales:      %s", strings.Join(o.TargetLocaleIDs, ", ")),
		fmt.Sprintf("Source files:        %d", len(o.SourceFiles)),
		fmt.Sprintf("Source issues:       %d", o.Issues.SourceIssuesCount),
		fmt.Sprintf("Translation issues:  %d", o.Issues.TranslationIssuesCount),
	}
}

// TableData returns the detail as a two-column field/value table.
func (o ViewOutput) TableData() ([]string, [][]string) {
	headers := []string{"FIELD", "VALUE"}
	rows := [][]string{
		{"TRANSLATION JOB UID", o.TranslationJobUID},
		{"NAME", o.JobName},
		{"NUMBER", o.JobNumber},
		{"STATUS", o.JobStatus},
		{"DESCRIPTION", o.Description},
		{"REFERENCE NUMBER", o.ReferenceNumber},
		{"DUE DATE", helpers.TimeToString(o.Dates.Due, time.RFC3339)},
		{"MODIFIED DATE", helpers.TimeToString(o.Dates.Modified, time.RFC3339)},
		{"CREATED DATE", helpers.TimeToString(o.Dates.Created, time.RFC3339)},
		{"PRIORITY", fmt.Sprintf("%d", o.Priority)},
		{"TARGET LOCALES", strings.Join(o.TargetLocaleIDs, ", ")},
	}
	return headers, rows
}

// RunView resolves the job by UID or name and returns its full detail.
func (s service) RunView(ctx context.Context, params ViewParams) (ViewOutput, error) {
	if err := params.Validate(); err != nil {
		return ViewOutput{}, fmt.Errorf("invalid view params: %w", err)
	}
	rlog.Debugf("running jobs view with params: %+v", params)

	jobUID, err := jobresolver.GetJobUID(ctx, s.job, params.ProjectUID, params.JobUIDOrName)
	if err != nil {
		return ViewOutput{}, fmt.Errorf("resolve job UID: %w", err)
	}

	detail, err := s.job.GetJob(ctx, params.ProjectUID, jobUID)
	if err != nil {
		return ViewOutput{}, fmt.Errorf("get job %q: %w", jobUID, err)
	}

	out := ViewOutput{GetJobResponse: detail}
	b, err := json.Marshal(detail)
	if err != nil {
		rlog.Errorf("failed to marshal job detail to JSON: %v", err)
		return out, nil
	}
	out.JSON = b
	return out, nil
}
