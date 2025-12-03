package jobs

import (
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	srv "github.com/Smartling/smartling-cli/services/jobs"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
)

// SrvInitializer defines files service initializer
type SrvInitializer interface {
	InitJobSrv() (srv.Service, error)
}

// NewSrvInitializer returns new SrvInitializer implementation
func NewSrvInitializer() SrvInitializer {
	return srvInitializer{}
}

type srvInitializer struct{}

// InitJobSrv initializes `job` service with the client and configuration.
func (i srvInitializer) InitJobSrv() (srv.Service, error) {
	client, err := rootcmd.Client()
	if err != nil {
		return nil, err
	}
	jobApi := jobapi.NewJob(client.Client)
	jobSrv := srv.NewService(jobApi)
	return jobSrv, nil
}
