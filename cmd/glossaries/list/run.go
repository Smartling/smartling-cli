package gllist

import (
	"context"

	glossariescmd "github.com/Smartling/smartling-cli/cmd/glossaries"
	"github.com/Smartling/smartling-cli/output"
	"github.com/Smartling/smartling-cli/output/static"
	srv "github.com/Smartling/smartling-cli/services/glossary"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
)

func run(ctx context.Context,
	initializer glossariescmd.SrvInitializer,
	params srv.ListParams,
	outputParams output.Params,
) error {
	rlog.Debugf("running glossary list with params: %v", params)
	glossarySrv, err := initializer.InitGlossarySrv(ctx)
	if err != nil {
		return clierror.UIError{
			Operation:   "init",
			Err:         err,
			Description: "unable to initialize Glossary service",
		}
	}

	listOutput, err := glossarySrv.RunList(ctx, params)
	if err != nil {
		return err
	}

	outputFormat := static.GetOutputFormat[srv.ListOutput](outputParams.Format)
	outputFormat.FormatAndRender(listOutput)

	return nil
}
