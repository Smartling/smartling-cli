package mt

import (
	"fmt"
	"slices"
	"strings"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	globfiles "github.com/Smartling/smartling-cli/services/helpers/glob_files"
	srv "github.com/Smartling/smartling-cli/services/mt"

	api "github.com/Smartling/api-sdk-go/api/mt"
	"github.com/spf13/cobra"
)

var (
	output         string
	noProgress     bool
	allowedOutputs = []string{
		"table",
		"json",
		"simple",
	}
	joinedAllowedOutputs = strings.Join(allowedOutputs, ", ")
)

// NewMTCmd ...
func NewMTCmd() *cobra.Command {
	mtCmd := &cobra.Command{
		Use:   "mt",
		Short: "mt...",
		Long:  `mt...`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !slices.Contains(allowedOutputs, output) {
				return fmt.Errorf("invalid output: %s (allowed: %s)", output, joinedAllowedOutputs)
			}
			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {

		},
	}

	mtCmd.PersistentFlags().StringVar(&output, "output", "simple", "Output format: "+joinedAllowedOutputs)
	mtCmd.PersistentFlags().BoolVar(&noProgress, "no-progress", false, "Disable progress indicators")

	return mtCmd
}

// SrvInitializer defines files service initializer
type SrvInitializer interface {
	InitMTSrv() (srv.Service, globfiles.ListFilesFn, error)
}

// NewSrvInitializer returns new SrvInitializer implementation
func NewSrvInitializer() SrvInitializer {
	return srvInitializer{}
}

type srvInitializer struct{}

// InitMTSrv initializes `mt` service with the client and configuration.
func (i srvInitializer) InitMTSrv() (srv.Service, globfiles.ListFilesFn, error) {
	client, err := rootcmd.Client()
	if err != nil {
		return nil, nil, err
	}
	downloader := api.NewDownloader(client.Client)
	fileTranslator := api.NewFileTranslator(client.Client)
	uploader := api.NewUploader(client.Client)
	translationControl := api.NewTranslationControl(client.Client)
	mtSrv := srv.NewService(downloader, fileTranslator, uploader, translationControl)
	return mtSrv, client.ListAllFiles, nil
}
