package glossary

import (
	"context"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	srv "github.com/Smartling/smartling-cli/services/glossary"

	glossaryapi "github.com/Smartling/api-sdk-go/api/glossary"
)

// SrvInitializer defines glossary service initializer
type SrvInitializer interface {
	InitGlossarySrv(ctx context.Context) (srv.Service, error)
}

// NewSrvInitializer returns new SrvInitializer implementation
func NewSrvInitializer() SrvInitializer {
	return srvInitializer{}
}

type srvInitializer struct{}

// InitGlossarySrv initializes `glossary` service with the client.
func (i srvInitializer) InitGlossarySrv(ctx context.Context) (srv.Service, error) {
	client, err := rootcmd.Client(ctx)
	if err != nil {
		return nil, err
	}
	glossaryApi := glossaryapi.NewGlossary(client.Client)
	glossarySrv := srv.NewService(glossaryApi)
	return glossarySrv, nil
}
