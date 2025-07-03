package mt

import (
	"fmt"
	"slices"
	"strings"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	srv "github.com/Smartling/smartling-cli/services/mt"

	api "github.com/Smartling/api-sdk-go/api/mt"
	"github.com/spf13/cobra"
)

var (
	output         string
	allowedOutputs = []string{
		"table",
		"json",
		"simple",
	}
	joinedAllowedOutputs = strings.Join(allowedOutputs, ", ")
	outputMode           string
	allowedOutputModes   = []string{
		"dynamic",
		"static",
	}
	joinedAllowedOutputModes = strings.Join(allowedOutputModes, ", ")
)

// NewMTCmd returns new mt command
func NewMTCmd() *cobra.Command {
	mtCmd := &cobra.Command{
		Use:   "mt",
		Short: "File Machine Translations",
		Long:  `Machine Translations offers a simple way to upload files and execute actions on them without any complex setup required`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !slices.Contains(allowedOutputs, output) {
				return fmt.Errorf("invalid output: %s (allowed: %s)", output, joinedAllowedOutputs)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 && cmd.Flags().NFlag() == 0 {
				return cmd.Help()
			}
			return nil
		},
	}

	mtCmd.PersistentFlags().StringVar(&output, "output", "simple", "Output format: "+joinedAllowedOutputs)
	mtCmd.PersistentFlags().StringVar(&outputMode, "output-mode", "static", "Output mode: "+joinedAllowedOutputModes)

	return mtCmd
}

// SrvInitializer defines files service initializer
type SrvInitializer interface {
	InitMTSrv() (srv.Service, error)
}

// NewSrvInitializer returns new SrvInitializer implementation
func NewSrvInitializer() SrvInitializer {
	return srvInitializer{}
}

type srvInitializer struct{}

// InitMTSrv initializes `mt` service with the client and configuration.
func (i srvInitializer) InitMTSrv() (srv.Service, error) {
	client, err := rootcmd.Client()
	if err != nil {
		return nil, err
	}
	downloader := api.NewDownloader(client.Client)
	fileTranslator := api.NewFileTranslator(client.Client)
	uploader := api.NewUploader(client.Client)
	translationControl := api.NewTranslationControl(client.Client)
	mtSrv := srv.NewService(downloader, fileTranslator, uploader, translationControl)
	return mtSrv, nil
}
