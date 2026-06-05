package jobstrings

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Smartling/smartling-cli/services/jobs/jobresolver"

	api "github.com/Smartling/api-sdk-go/api/job/string"
)

// ListParams defines the list-strings params.
type ListParams struct {
	ProjectID      string
	JobUIDOrName   string
	TargetLocaleID string
	Limit          uint32
	Offset         uint32
}

// Validate checks that ListParams are valid.
func (p ListParams) Validate() error {
	return validateIDs(p.ProjectID, p.JobUIDOrName)
}

// Item is a single string row in a job.
type Item struct {
	TargetLocaleID string `json:"targetLocaleId"`
	Hashcode       string `json:"hashcode"`
}

// ListOutput is the result of listing a job's strings.
type ListOutput struct {
	TotalCount uint32 `json:"totalCount"`
	Items      []Item `json:"items"`

	JSON []byte `json:"-"`
}

func newListOutput(resp api.ListResponse) (ListOutput, error) {
	o := ListOutput{
		Items:      make([]Item, len(resp.Items)),
		TotalCount: resp.TotalCount,
	}
	for i, item := range resp.Items {
		o.Items[i] = Item{TargetLocaleID: item.TargetLocaleID, Hashcode: item.Hashcode}
	}
	var err error
	if o.JSON, err = json.Marshal(o); err != nil {
		return ListOutput{}, err
	}
	return o, nil
}

// JSONBytes returns the JSON representation of the list.
func (o ListOutput) JSONBytes() []byte { return o.JSON }

// SimpleLines returns a human-readable summary of the list.
func (o ListOutput) SimpleLines() []string {
	if len(o.Items) == 0 {
		return []string{"No strings found."}
	}
	lines := make([]string, 0, len(o.Items)+1)
	for _, it := range o.Items {
		lines = append(lines, fmt.Sprintf("%s  %s", it.TargetLocaleID, it.Hashcode))
	}
	lines = append(lines, fmt.Sprintf("Showing %d of %d string(s). Use --limit/--offset to page.", len(o.Items), o.TotalCount))
	return lines
}

// TableData returns the list as one row per string.
func (o ListOutput) TableData() ([]string, [][]string) {
	headers := []string{"TARGET LOCALE ID", "HASHCODE"}
	rows := make([][]string, 0, len(o.Items))
	for _, it := range o.Items {
		rows = append(rows, []string{it.TargetLocaleID, it.Hashcode})
	}
	return headers, rows
}

// RunList retrieves the strings assigned to a translation job.
func (s service) RunList(ctx context.Context, params ListParams) (ListOutput, error) {
	if err := params.Validate(); err != nil {
		return ListOutput{}, err
	}
	jobUID, err := jobresolver.GetJobUID(ctx, s.job, params.ProjectID, params.JobUIDOrName)
	if err != nil {
		return ListOutput{}, err
	}
	resp, err := s.jobString.List(ctx, params.ProjectID, jobUID, api.ListParams{
		TargetLocaleID: params.TargetLocaleID,
		Limit:          params.Limit,
		Offset:         params.Offset,
	})
	if err != nil {
		return ListOutput{}, err
	}
	return newListOutput(resp)
}
