package jobstrings

import (
	"context"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	srv "github.com/Smartling/smartling-cli/services/jobs/strings"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	stringapi "github.com/Smartling/api-sdk-go/api/job/string"
)

// SrvInitializer defines job strings service initializer
type SrvInitializer interface {
	InitJobStringsSrv(ctx context.Context) (srv.Service, error)
}

// NewSrvInitializer returns new SrvInitializer implementation
func NewSrvInitializer() SrvInitializer {
	return srvInitializer{}
}

type srvInitializer struct{}

// InitJobStringsSrv initializes the job strings service with the client.
func (i srvInitializer) InitJobStringsSrv(ctx context.Context) (srv.Service, error) {
	client, err := rootcmd.Client(ctx)
	if err != nil {
		return nil, err
	}
	stringsSrv := srv.NewService(stringapi.NewJobString(client.Client), jobapi.NewJob(client.Client))
	return stringsSrv, nil
}
