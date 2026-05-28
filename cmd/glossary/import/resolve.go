package glimport

import (
	"strings"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	glossarycmd "github.com/Smartling/smartling-cli/cmd/glossary"
	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	srv "github.com/Smartling/smartling-cli/services/glossary"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

func resolveParams(cmd *cobra.Command, fileConfig glossarycmd.FileConfig, glossaryUIDOrName, inFile string) (srv.ImportParams, error) {
	rlog.Debugf("resolving params")

	cnf, err := rootcmd.Config()
	if err != nil {
		return srv.ImportParams{}, clierror.UIError{
			Operation:   "config",
			Err:         err,
			Description: "failed to read config",
		}
	}

	accountUID, err := resolve.FallbackAccount(cmd.Root().PersistentFlags().Lookup("account"), cnf.AccountID)
	if err != nil {
		return srv.ImportParams{}, err
	}

	cfg := fileConfig.Glossary.Export

	fileType := resolve.FallbackString(cmd.Flags().Lookup(fileTypeFlag), resolve.StringParam{FlagName: fileTypeFlag, Config: &cfg.FileType})
	fileType = strings.ToLower(fileType)
	tbxVersion := resolve.FallbackString(cmd.Flags().Lookup(tbxVersionFlag), resolve.StringParam{FlagName: tbxVersionFlag, Config: &cfg.TbxVersion})
	tbxVersion = strings.ToLower(tbxVersion)
	params := srv.ImportParams{
		AccountUID:        accountUID,
		GlossaryUIDOrName: glossaryUIDOrName,
		ArchiveMode:       resolve.FallbackBool(cmd.Flags().Lookup(archiveModeFlag), resolve.BoolParam{FlagName: archiveModeFlag, Config: &cfg.ArchiveMode}),
		ImportFile: srv.ImportFile{
			Path:      "",
			Name:      "",
			MediaType: "",
		},
	}

	if params.Filter.Created.Date, err = resolve.FallbackDate(cmd, createdDateFlag, cfg.Filter.Created.Date); err != nil {
		return srv.ImportParams{}, err
	}
	if params.Filter.LastModified.Date, err = resolve.FallbackDate(cmd, lastModifiedDateFlag, cfg.Filter.LastModified.Date); err != nil {
		return srv.ImportParams{}, err
	}

	return params, nil
}
