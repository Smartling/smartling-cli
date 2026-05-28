package glimport

import (
	"context"
	"errors"
	"fmt"

	glossarycmd "github.com/Smartling/smartling-cli/cmd/glossary"
	"github.com/Smartling/smartling-cli/output"
	"github.com/Smartling/smartling-cli/output/static"
	srv "github.com/Smartling/smartling-cli/services/glossary"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	glossaryapi "github.com/Smartling/api-sdk-go/api/glossary"
)

func run(ctx context.Context,
	initializer glossarycmd.SrvInitializer,
	params srv.ImportParams,
	outputParams output.Params,
) error {
	rlog.Debugf("running glossary import with params: %v", params)
	glossarySrv, err := initializer.InitGlossarySrv(ctx)
	if err != nil {
		return clierror.UIError{
			Operation:   "init",
			Err:         err,
			Description: "unable to initialize Glossary service",
		}
	}

	importOutput, err := glossarySrv.RunImport(ctx, params)
	if err != nil {
		switch {
		case errors.Is(err, glossaryapi.ErrGlossaryNotFound):
			return clierror.UIError{
				Operation:   "find glossary",
				Err:         err,
				Description: fmt.Sprintf("no glossary found for %q", params.GlossaryUIDOrName),
			}
		case errors.Is(err, glossaryapi.ErrImportNotFound):
			return clierror.UIError{
				Operation:   "find import",
				Err:         err,
				Description: fmt.Sprintf("no glossary import found for %q", params.GlossaryUIDOrName),
			}
		}
		return err
	}

	outputFormat := static.GetOutputFormat[srv.ImportOutput](outputParams.Format)
	outputFormat.FormatAndRender(importOutput)

	return nil
}
