package glcreate

import (
	"context"
	"errors"
	"fmt"

	glossariescmd "github.com/Smartling/smartling-cli/cmd/glossaries"
	"github.com/Smartling/smartling-cli/output"
	"github.com/Smartling/smartling-cli/output/static"
	srv "github.com/Smartling/smartling-cli/services/glossary"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	glossaryapi "github.com/Smartling/api-sdk-go/api/glossary"
)

func run(ctx context.Context, initializer glossariescmd.SrvInitializer, params srv.CreateParams, outputParams output.Params) error {
	rlog.Debugf("running glossary create with params: %v", params)
	glossarySrv, err := initializer.InitGlossarySrv(ctx)
	if err != nil {
		return clierror.UIError{
			Operation:   "init",
			Err:         err,
			Description: "unable to initialize Glossary service",
		}
	}

	createOutput, err := glossarySrv.RunCreate(ctx, params)
	if err != nil {
		if errors.Is(err, glossaryapi.ErrGlossaryNotFound) {
			return clierror.UIError{
				Operation:   "create glossary",
				Err:         err,
				Description: fmt.Sprintf("failed to create glossary %q", params.GlossaryName),
			}
		}
		return err
	}

	outputFormat := static.GetOutputFormat[srv.CreateOutput](outputParams.Format)
	outputFormat.FormatAndRender(createOutput)

	return nil
}
