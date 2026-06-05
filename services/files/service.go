package files

import (
	"context"
	"time"

	"github.com/Smartling/smartling-cli/services/helpers/config"

	sdk "github.com/Smartling/api-sdk-go"
	batchapi "github.com/Smartling/api-sdk-go/api/batches"
	jobapi "github.com/Smartling/api-sdk-go/api/job"
	jobfile "github.com/Smartling/api-sdk-go/api/job/file"
)

const defaultJobNameTemplate = "CLI uploads"

var (
	pollingInterval = 5 * time.Second
	pollingDuration = 5 * time.Minute
)

// ListJobFilesFn returns a page of a translation job's source files. It lets the
// pull service list job files without depending on the full JobFile interface.
type ListJobFilesFn func(ctx context.Context, projectID, jobUID string, limit, offset uint32) (jobfile.ListResponse, error)

// Service defines behaviors to interact with Smartling files.
type Service interface {
	RunDelete(ctx context.Context, uri string) error
	RunImport(ctx context.Context, params ImportParams) error
	RunList(ctx context.Context, formatType string, short bool, uri string) error
	RunPull(ctx context.Context, params PullParams) error
	RunPush(ctx context.Context, params PushParams) error
	RunRename(ctx context.Context, oldURI, newURI string) error
	RunStatus(ctx context.Context, params StatusParams) error
}

// service provides methods to interact with Smartling files.
type service struct {
	APIClient    sdk.APIClient
	BatchApi     batchapi.Batch
	JobApi       jobapi.Job
	ListJobFiles ListJobFilesFn
	Config       config.Config
	FileConfig   config.FileConfig
}

// NewService creates a new instance of the Service with the provided client, and configurations.
func NewService(client sdk.APIClient,
	batchApi batchapi.Batch,
	jobApi jobapi.Job,
	listJobFiles ListJobFilesFn,
	config config.Config,
	fileConfig config.FileConfig,
) Service {
	return &service{
		APIClient:    client,
		BatchApi:     batchApi,
		JobApi:       jobApi,
		ListJobFiles: listJobFiles,
		Config:       config,
		FileConfig:   fileConfig,
	}
}
