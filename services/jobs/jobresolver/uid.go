package jobresolver

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
)

var (
	ErrEmptyJob   = fmt.Errorf("job UID or name must be specified")
	jobUIDPattern = regexp.MustCompile(`^[a-z0-9]{12}$`)
)

func GetJobUID(ctx context.Context, api jobapi.Job, projectUID, jobUIDOrName string) (string, error) {
	if jobUIDOrName == "" {
		return "", ErrEmptyJob
	}
	if jobUIDPattern.MatchString(jobUIDOrName) {
		jb, err := api.GetJob(ctx, projectUID, jobUIDOrName)
		switch {
		case err == nil:
			return jb.TranslationJobUID, nil
		case errors.Is(err, jobapi.ErrNotFound):
			// 12-char input wasn't a UID — could still be a job name, fall through
		default:
			return "", fmt.Errorf("get job by UID %q: %w", jobUIDOrName, err)
		}
	}

	jobs, err := api.SearchByName(ctx, projectUID, jobUIDOrName)
	if err != nil {
		return "", fmt.Errorf("search jobs by name %q: %w", jobUIDOrName, err)
	}
	if len(jobs) == 0 {
		return "", jobapi.ErrNotFound
	}
	j, found := jobapi.FindFirstJobByName(jobs, jobUIDOrName)
	if !found {
		return "", jobapi.ErrNotFound
	}
	if strings.TrimSpace(j.TranslationJobUID) == "" {
		return "", jobapi.ErrNotFound
	}
	return j.TranslationJobUID, nil
}
