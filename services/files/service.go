package files

import (
	"context"
	"time"

	"github.com/Smartling/smartling-cli/services/helpers/config"

	sdk "github.com/Smartling/api-sdk-go"
	sdkjobs "github.com/Smartling/api-sdk-go/api/jobs"
)

var pollingInterval = time.Second

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
	Batch      sdkjobs.Batch
	Config     config.Config
	FileConfig config.FileConfig
}

// NewService creates a new instance of the Service with the provided client, and configurations.
func NewService(client sdk.APIClient,
	batch sdkjobs.Batch,
	config config.Config,
	fileConfig config.FileConfig,
) Service {
	return &service{
		APIClient:  client,
		Batch:      batch,
		Config:     config,
		FileConfig: fileConfig,
	}
}
