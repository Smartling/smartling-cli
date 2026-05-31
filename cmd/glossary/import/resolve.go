package glimport

import (
	"path/filepath"
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
	rlog.Debugf("resolving import params")

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

	cfg := fileConfig.Glossary.Import

	archiveMode := resolve.FallbackBool(cmd.Flags().Lookup(archiveModeFlag), resolve.BoolParam{FlagName: archiveModeFlag, Config: &cfg.ArchiveMode})
	mediaType := resolve.FallbackString(cmd.Flags().Lookup(mediaTypeFlag), resolve.StringParam{FlagName: mediaTypeFlag, Config: &cfg.MediaType})
	if mediaType == "" {
		mediaType = mediaTypeFromPath(inFile)
	}

	return srv.ImportParams{
		AccountUID:        accountUID,
		GlossaryUIDOrName: glossaryUIDOrName,
		ArchiveMode:       archiveMode,
		ImportFile: srv.ImportFile{
			Path:      inFile,
			Name:      filepath.Base(inFile),
			MediaType: mediaType,
		},
	}, nil
}

// mediaTypeFromPath maps a file extension to the importFileMediaType enum; "" means unknown.
func mediaTypeFromPath(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".csv":
		return "text/csv"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".tbx", ".xml":
		return "text/xml"
	default:
		return ""
	}
}
