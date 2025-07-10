package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

var (
	smartlingURL       string
	configFile         string
	project            string
	account            string
	user               string
	secret             string
	operationDirectory string
	threads            uint32
	insecure           bool
	proxy              string
	verbose            int

	isInit     bool
	isFiles    bool
	isProjects bool
	isList     bool
)

// NewRootCmd creates a new root command.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "smartling-cli",
		Short:   "Manage translation files using Smartling CLI.",
		Version: "2.1",
		Long: `Manage translation files using Smartling CLI.
                Complete documentation is available at https://www.smartling.com`,
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			configureLoggerVerbose()

			path := cmd.CommandPath()
			isInit = strings.HasPrefix(path, "smartling-cli init")
			isFiles = strings.HasPrefix(path, "smartling-cli files")
			isProjects = strings.HasPrefix(path, "smartling-cli projects")
			isList = strings.HasPrefix(path, "smartling-cli list")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 && cmd.Flags().NFlag() == 0 {
				if err := cmd.Help(); err != nil {
					return err
				}
				return nil
			}
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", `Config file in YAML format.
By default CLI will look for file named
"smartling.yml" in current directory and in all
intermediate parents, emulating git behavior.`)
	rootCmd.PersistentFlags().StringVarP(&project, "project", "p", "", `Project ID to operate on.
This option overrides config value "project_id".`)
	rootCmd.PersistentFlags().StringVarP(&account, "account", "a", "", `Account ID to operate on.
This option overrides config value "account_id".`)
	rootCmd.PersistentFlags().StringVar(&user, "user", "", `User ID which will be used for authentication.
This option overrides config value "user_id".`)
	rootCmd.PersistentFlags().StringVar(&secret, "secret", "", `Token Secret which will be used for authentication.
This option overrides config value "secret".`)
	rootCmd.PersistentFlags().StringVar(&operationDirectory, "operation-directory", ".", `Sets directory to operate on, usually, to store or to
read files.  Depends on command.`)
	rootCmd.PersistentFlags().Uint32Var(&threads, "threads", 4, `If command can be executed concurrently, it will be
executed for at most <number> of threads.`)
	rootCmd.PersistentFlags().BoolVarP(&insecure, "insecure", "k", false, "Skip HTTPS certificate validation.")
	rootCmd.PersistentFlags().StringVar(&proxy, "proxy", "", "Use specified URL as proxy server.")
	rootCmd.PersistentFlags().StringVar(&smartlingURL, "smartling-url", "", `Specify base Smartling URL, merely for testing
purposes.`)
	rootCmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "Verbose logging")

	return rootCmd
}
