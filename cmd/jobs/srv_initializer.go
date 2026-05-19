package jobs

import (
	"context"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	srv "github.com/Smartling/smartling-cli/services/jobs"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
)

// SrvInitializer defines jobs service initializer
type SrvInitializer interface {
	InitJobSrv(ctx context.Context) (srv.Service, error)
}

// NewSrvInitializer returns new SrvInitializer implementation
func NewSrvInitializer() SrvInitializer {
	return srvInitializer{}
}

type srvInitializer struct{}

// InitJobSrv initializes `job` service with the client and configuration.
func (i srvInitializer) InitJobSrv(ctx context.Context) (srv.Service, error) {
	client, err := rootcmd.Client(ctx)
	if err != nil {
		return nil, err
	}
	jobApi := jobapi.NewJob(client.Client)
	jobSrv := srv.NewService(jobApi)
	return jobSrv, nil
}
