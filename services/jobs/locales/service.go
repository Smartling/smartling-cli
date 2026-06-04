package locales

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	api "github.com/Smartling/api-sdk-go/api/job/locale"
)

// Service defines behavior for managing a translation job's target locales.
type Service interface {
	RunAdd(ctx context.Context, params AddParams) (Output, error)
	RunRemove(ctx context.Context, params RemoveParams) (Output, error)
}

// NewService creates a new implementation of the Service. The job API is used to
// resolve a job UID from a UID-or-name; locale performs the add/remove.
func NewService(jobLocale api.JobLocale, job jobapi.Job) Service {
	return service{
		locale: jobLocale,
		job:    job,
	}
}

type service struct {
	job    jobapi.Job
	locale api.JobLocale
}

// Output is the result of an add/remove locale operation.
type Output struct {
	Action            string `json:"action"`
	ProjectUID        string `json:"projectUid"`
	TranslationJobUID string `json:"translationJobUid"`
	TargetLocaleID    string `json:"targetLocaleId"`

	JSON []byte `json:"-"`
}

func newOutput(action, projectUID, jobUID, localeID string) (Output, error) {
	o := Output{
		Action:            action,
		ProjectUID:        projectUID,
		TranslationJobUID: jobUID,
		TargetLocaleID:    localeID,
	}
	var err error
	o.JSON, err = json.Marshal(o)
	if err != nil {
		return Output{}, err
	}
	return o, err
}

// JSONBytes returns the JSON representation of the result.
func (o Output) JSONBytes() []byte { return o.JSON }

// SimpleLines returns a human-readable summary of the result.
func (o Output) SimpleLines() []string {
	return []string{fmt.Sprintf("Locale %s %s job %s", o.TargetLocaleID, o.Action, o.TranslationJobUID)}
}

// TableData returns the result as a single-row table.
func (o Output) TableData() ([]string, [][]string) {
	return []string{"ACTION", "PROJECT UID", "TRANSLATION JOB UID", "TARGET LOCALE ID"},
		[][]string{{o.Action, o.ProjectUID, o.TranslationJobUID, o.TargetLocaleID}}
}

func validateParams(projectID, translationJobUID, targetLocaleID string) error {
	switch {
	case projectID == "":
		return errors.New("project ID is required")
	case translationJobUID == "":
		return errors.New("translation job UID or name is required")
	case targetLocaleID == "":
		return errors.New("target locale ID is required")
	}
	return nil
}
