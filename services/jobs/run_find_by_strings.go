package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Smartling/smartling-cli/services/helpers"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	smerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
)

// MaxFindByStringsRecords is the API limit on hashcodes×locales the
// find-jobs-by-strings endpoint accepts per request.
const MaxFindByStringsRecords = 20000

// FindByStringsParams carries the find-jobs-by-strings request.
type FindByStringsParams struct {
	ProjectUID string
	Hashcodes  []string
	LocaleIDs  []string
}

// Validate checks required fields and the API record limit. The endpoint counts
// records as hashcodes×locales (hashcodes alone when no locales are given).
func (p FindByStringsParams) Validate() error {
	if p.ProjectUID == "" {
		return smerror.ErrEmptyParam("ProjectUID")
	}
	if len(p.Hashcodes) == 0 {
		return smerror.ErrEmptyParam("Hashcodes")
	}
	locales := len(p.LocaleIDs)
	if locales == 0 {
		locales = 1
	}
	if len(p.Hashcodes) > MaxFindByStringsRecords/locales {
		return fmt.Errorf("too many records: %d hashcodes × %d locales exceeds the limit of %d",
			len(p.Hashcodes), locales, MaxFindByStringsRecords)
	}
	return nil
}

// FindByStringsMatch is a single hashcode+locale match against one job.
type FindByStringsMatch struct {
	Hashcode          string `json:"hashcode"`
	LocaleID          string `json:"localeId"`
	TranslationJobUID string `json:"translationJobUid"`
	JobName           string `json:"jobName"`
	DueDate           string `json:"dueDate"`
}

// FindByStringsOutput is the result of find-jobs-by-strings, flattened to one
// row per hashcode+locale+job match. TotalCount is the number of matching jobs
// reported by the API (distinct from the number of flattened rows).
type FindByStringsOutput struct {
	Matches    []FindByStringsMatch
	TotalCount int
	JSON       []byte `json:"-"`
}

// JSONBytes returns the raw JSON payload.
func (o FindByStringsOutput) JSONBytes() []byte { return o.JSON }

// SimpleLines returns a human-readable summary.
func (o FindByStringsOutput) SimpleLines() []string {
	if len(o.Matches) == 0 {
		return []string{"No jobs found for the given strings."}
	}
	lines := make([]string, 0, len(o.Matches)+1)
	for _, m := range o.Matches {
		lines = append(lines, fmt.Sprintf("%s  %s  %s  %s  %s",
			m.Hashcode, m.LocaleID, m.TranslationJobUID, m.JobName, m.DueDate))
	}
	lines = append(lines, fmt.Sprintf("%d job(s) matched.", o.TotalCount))
	return lines
}

// TableData returns the matches as a table.
func (o FindByStringsOutput) TableData() ([]string, [][]string) {
	headers := []string{"HASHCODE", "LOCALE", "TRANSLATION JOB UID", "JOB NAME", "DUE DATE"}
	rows := make([][]string, 0, len(o.Matches))
	for _, m := range o.Matches {
		rows = append(rows, []string{m.Hashcode, m.LocaleID, m.TranslationJobUID, m.JobName, m.DueDate})
	}
	return headers, rows
}

// RunFindByStrings finds jobs containing the given strings (by hashcode) in the
// given locales.
func (s service) RunFindByStrings(ctx context.Context, params FindByStringsParams) (FindByStringsOutput, error) {
	if err := params.Validate(); err != nil {
		return FindByStringsOutput{}, fmt.Errorf("invalid find-by-strings params: %w", err)
	}
	rlog.Debugf("running jobs find-by-strings with params: %+v", params)

	reqParams := jobapi.FindJobsByStringsRequest{
		Hashcodes: params.Hashcodes,
		LocaleIDs: params.LocaleIDs,
	}
	resp, err := s.job.FindJobsByStrings(ctx, params.ProjectUID, reqParams)
	if err != nil {
		return FindByStringsOutput{}, fmt.Errorf("failed to find jobs by strings: %w", err)
	}

	return toFindByStringsOutput(resp)
}

// findByStringsJSON is the JSON shape for find-by-strings output.
type findByStringsJSON struct {
	Matches    []FindByStringsMatch `json:"matches"`
	TotalCount int                  `json:"totalCount"`
}

func toFindByStringsOutput(resp jobapi.FindJobsByStringsResponse) (FindByStringsOutput, error) {
	matches := make([]FindByStringsMatch, 0, len(resp.Items))
	for _, job := range resp.Items {
		dueDate := helpers.TimeToString(job.DueDate, time.RFC3339)
		for _, byLocale := range job.HashcodesByLocale {
			for _, hashcode := range byLocale.Hashcodes {
				matches = append(matches, FindByStringsMatch{
					Hashcode:          hashcode,
					LocaleID:          byLocale.LocaleID,
					TranslationJobUID: job.TranslationJobUID,
					JobName:           job.JobName,
					DueDate:           dueDate,
				})
			}
		}
	}

	b, err := json.Marshal(findByStringsJSON{Matches: matches, TotalCount: resp.TotalCount})
	if err != nil {
		return FindByStringsOutput{}, fmt.Errorf("marshal find-by-strings to JSON: %w", err)
	}
	return FindByStringsOutput{Matches: matches, TotalCount: resp.TotalCount, JSON: b}, nil
}
