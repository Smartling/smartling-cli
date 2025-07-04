package detect

import (
	"fmt"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	srv "github.com/Smartling/smartling-cli/services/mt"

	api "github.com/Smartling/api-sdk-go/api/mt"
	"github.com/spf13/cobra"
)

func resolveParams(cmd *cobra.Command, fileConfig mtcmd.FileConfig, fileOrPattern string) (srv.DetectParams, error) {
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
	var accountIDConfig *string
	if cnf.AccountID != "" {
		accountIDConfig = &cnf.AccountID
	}
	accountUIDParam := resolve.FallbackString(cmd.Root().PersistentFlags().Lookup("account"), resolve.StringParam{
		FlagName: "account",
		Config:   accountIDConfig,
	})
	return srv.DetectParams{
		FileType:       fileTypeParam,
		InputDirectory: inputDirectoryParam,
		FileOrPattern:  fileOrPattern,
		AccountUID:     api.AccountUID(accountUIDParam),
	}, nil
}
