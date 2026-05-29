package glossary

import (
	"context"

	api "github.com/Smartling/api-sdk-go/api/glossary"
)

type Service interface {
	RunImport(ctx context.Context, params ImportParams) (ImportOutput, error)
	RunExport(ctx context.Context, params ExportParams) (ExportOutput, error)
	RunCreate(ctx context.Context, params CreateParams) (CreateOutput, error)
}

func NewService(glossaryApi api.Glossary) Service {
	return service{glossaryApi: glossaryApi}
}

type service struct {
	glossaryApi api.Glossary
}
