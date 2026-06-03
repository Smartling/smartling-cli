package glossary

import (
	"context"

	api "github.com/Smartling/api-sdk-go/api/glossary"
)

// Service is the glossary business-logic interface.
type Service interface {
	RunImport(ctx context.Context, params ImportParams) (ImportOutput, error)
	RunExport(ctx context.Context, params ExportParams) (ExportOutput, error)
	RunCreate(ctx context.Context, params CreateParams) (CreateOutput, error)
	RunList(ctx context.Context, params ListParams) (ListOutput, error)
}

// NewService builds a glossary Service.
func NewService(glossaryApi api.Glossary) Service {
	return service{glossaryApi: glossaryApi}
}

type service struct {
	glossaryApi api.Glossary
}
