package jobfiles

import (
	"context"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	srv "github.com/Smartling/smartling-cli/services/jobs/files"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	filejob "github.com/Smartling/api-sdk-go/api/job/file"
)

// SrvInitializer defines job files service initializer
type SrvInitializer interface {
	InitJobFilesSrv(ctx context.Context) (srv.Service, error)
}

// NewSrvInitializer returns new SrvInitializer implementation
func NewSrvInitializer() SrvInitializer {
	return srvInitializer{}
}

type srvInitializer struct{}

// InitJobFilesSrv initializes the job files service. It wires the JobFile API
// (add/remove), the Job API (UID-or-name resolution) and ListAllFiles (used to
// expand --file glob patterns against the project's files).
func (i srvInitializer) InitJobFilesSrv(ctx context.Context) (srv.Service, error) {
	client, err := rootcmd.Client(ctx)
	if err != nil {
		return nil, err
	}
	filesSrv := srv.NewService(
		filejob.NewJobFile(client.Client),
		jobapi.NewJob(client.Client),
		client.ListAllFiles,
	)
	return filesSrv, nil
}
