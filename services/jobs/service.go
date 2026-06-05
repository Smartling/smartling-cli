package jobs

import (
	"context"

	api "github.com/Smartling/api-sdk-go/api/job"
)

// Service defines behavior for interacting with Smartling jobs.
type Service interface {
	RunProgress(ctx context.Context, p ProgressParams) (ProgressOutput, error)
	RunList(ctx context.Context, p ListParams) (ListOutput, error)
	RunView(ctx context.Context, p ViewParams) (ViewOutput, error)
}

// NewService creates a new implementation of the Service
func NewService(job api.Job) Service {
	return service{
		job: job,
	}
}

type service struct {
	job api.Job
}
