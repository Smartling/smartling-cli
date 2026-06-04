package locales

import (
	"context"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	srv "github.com/Smartling/smartling-cli/services/jobs/locales"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	localeapi "github.com/Smartling/api-sdk-go/api/job/locale"
)

// SrvInitializer defines job locales service initializer
type SrvInitializer interface {
	InitJobLocalesSrv(ctx context.Context) (srv.Service, error)
}

// NewSrvInitializer returns new SrvInitializer implementation
func NewSrvInitializer() SrvInitializer {
	return srvInitializer{}
}

type srvInitializer struct{}

// InitJobLocalesSrv initializes the job locales service with the client.
func (i srvInitializer) InitJobLocalesSrv(ctx context.Context) (srv.Service, error) {
	client, err := rootcmd.Client(ctx)
	if err != nil {
		return nil, err
	}
	localeSrv := srv.NewService(localeapi.NewJobLocale(client.Client), jobapi.NewJob(client.Client))
	return localeSrv, nil
}
