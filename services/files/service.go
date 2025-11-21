package files

import (
	"context"
	"time"

	"github.com/Smartling/smartling-cli/services/helpers/config"

	sdk "github.com/Smartling/api-sdk-go"
	batchapi "github.com/Smartling/api-sdk-go/api/batches"
	jobapi "github.com/Smartling/api-sdk-go/api/job"
)

const defaultJobNameTemplate = "CLI uploads"

var (
	pollingInterval = 5 * time.Second
	pollingDuration = 5 * time.Minute
)

// Service defines behaviors to interact with Smartling files.
type Service interface {
	RunDelete(uri string) error
	RunImport(params ImportParams) error
	RunList(formatType string, short bool, uri string) error
	RunPull(params PullParams) error
	RunPush(ctx context.Context, params PushParams) error
	RunRename(oldURI, newURI string) error
	RunStatus(params StatusParams) error
}

// service provides methods to interact with Smartling files.
type service struct {
	APIClient  sdk.APIClient
	BatchApi   batchapi.Batch
	JobApi     jobapi.Job
	Config     config.Config
	FileConfig config.FileConfig
}

// NewService creates a new instance of the Service with the provided client, and configurations.
func NewService(client sdk.APIClient,
	batchApi batchapi.Batch,
	jobApi jobapi.Job,
	config config.Config,
	fileConfig config.FileConfig,
) Service {
	return &service{
		APIClient:  client,
		BatchApi:   batchApi,
		JobApi:     jobApi,
		Config:     config,
		FileConfig: fileConfig,
	}
}
