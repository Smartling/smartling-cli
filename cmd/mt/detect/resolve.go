package detect

import (
	"fmt"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/mt"

	"github.com/spf13/cobra"
)

func resolveParams(cmd *cobra.Command, fileConfig mtcmd.FileConfig, fileOrPattern string) (srv.DetectParams, error) {
	rlog.Debugf("resolving params")
	fileTypeParam := resolve.FallbackString(cmd.Flags().Lookup(fileTypeFlag), resolve.StringParam{
		FlagName: fileTypeFlag,
	})
	inputDirectoryParam := resolve.FallbackString(cmd.Flags().Lookup(inputDirectoryFlag), resolve.StringParam{
		FlagName: inputDirectoryFlag,
		Config:   fileConfig.MT.InputDirectory,
	})
	cnf, err := rootcmd.Config()
	if err != nil {
		return srv.DetectParams{}, fmt.Errorf("unable to read config: %w", err)
	}
	accountUID, err := resolve.FallbackAccount(cmd.Root().PersistentFlags().Lookup("account"), cnf.AccountID)
	if err != nil {
		return srv.DetectParams{}, err
	}
	return srv.DetectParams{
		FileType:       fileTypeParam,
		InputDirectory: inputDirectoryParam,
		FileOrPattern:  fileOrPattern,
		AccountUID:     accountUID,
	}, nil
}
