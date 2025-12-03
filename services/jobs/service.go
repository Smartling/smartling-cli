package jobs

import (
	"context"

	api "github.com/Smartling/api-sdk-go/api/job"
)

// Service defines behavior for interacting with Smartling MT.
type Service interface {
	RunProgress(ctx context.Context, p ProgressParams) (ProgressOutput, error)
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
