package jobstrings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	api "github.com/Smartling/api-sdk-go/api/job/string"
)

// Service defines behavior for managing a translation job's strings.
type Service interface {
	RunAdd(ctx context.Context, params AddParams) (MutateOutput, error)
	RunRemove(ctx context.Context, params RemoveParams) (MutateOutput, error)
	RunList(ctx context.Context, params ListParams) (ListOutput, error)
}

// NewService creates a new implementation of the Service. The job API resolves a
// job UID from a UID-or-name; jobString performs the add/remove/list.
func NewService(jobString api.JobString, job jobapi.Job) Service {
	return service{
		jobString: jobString,
		job:       job,
	}
}

type service struct {
	jobString api.JobString
	job       jobapi.Job
}

// MutateOutput is the result of an add/remove strings operation. SuccessCount
// and FailCount come from the API and reflect what actually happened - a
// nonexistent hashcode is ignored by the API and counted in neither.
type MutateOutput struct {
	Action            string   `json:"action"`
	ProjectUID        string   `json:"projectUid"`
	TranslationJobUID string   `json:"translationJobUid"`
	Hashcodes         []string `json:"hashcodes"`
	LocaleIDs         []string `json:"localeIds,omitempty"`
	SuccessCount      int      `json:"successCount"`
	FailCount         int      `json:"failCount"`

	JSON []byte `json:"-"`
}

func newMutateOutput(action, projectUID, jobUID string, hashcodes, localeIDs []string, successCount, failCount int) (MutateOutput, error) {
	o := MutateOutput{
		Action:            action,
		ProjectUID:        projectUID,
		TranslationJobUID: jobUID,
		Hashcodes:         hashcodes,
		LocaleIDs:         localeIDs,
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

// SimpleLines returns a human-readable summary of the result.
func (o MutateOutput) SimpleLines() []string {
	lines := []string{fmt.Sprintf("Strings %s — job %s: %d succeeded, %d failed",
		o.Action, o.TranslationJobUID, o.SuccessCount, o.FailCount)}
	if o.SuccessCount == 0 && o.FailCount == 0 {
		lines = append(lines, fmt.Sprintf(
			"No strings were affected. Verify the hashcodes exist in the project: %s",
			strings.Join(o.Hashcodes, ", ")))
	}
	return lines
}

// TableData returns the result as a single-row table.
func (o MutateOutput) TableData() ([]string, [][]string) {
	return []string{"ACTION", "TRANSLATION JOB UID", "SUCCEEDED", "FAILED"},
		[][]string{{o.Action, o.TranslationJobUID, strconv.Itoa(o.SuccessCount), strconv.Itoa(o.FailCount)}}
}

func validateIDs(projectID, jobUIDOrName string) error {
	switch {
	case projectID == "":
		return errors.New("project ID is required")
	case jobUIDOrName == "":
		return errors.New("translation job UID or name is required")
	}
	return nil
}

func validateMutate(projectID, jobUIDOrName string, hashcodes []string) error {
	if err := validateIDs(projectID, jobUIDOrName); err != nil {
		return err
	}
	if len(hashcodes) == 0 {
		return errors.New("at least one --hashcode is required")
	}
	return nil
}
