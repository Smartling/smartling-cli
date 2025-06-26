package cmd

import (
	"os"
	"strings"

	"github.com/Smartling/smartling-cli/services/helpers/client"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	sdk "github.com/Smartling/api-sdk-go"
	"github.com/kovetskiy/lorg"
	"github.com/spf13/cobra"
)

var (
	smartlingURL string
	configFile   string
	project      string
	account      string
	user         string
	secret       string
	directory    string
	threads      uint32
	insecure     bool
	proxy        string
	verbose      int

	isInit     bool
	isFiles    bool
	isProjects bool
	isList     bool
)

// NewRootCmd creates a new root command.
func NewRootCmd() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:     "smartling-cli",
		Short:   "Manage translation files using Smartling CLI.",
		Version: "1.7",
		Long: `Manage translation files using Smartling CLI.
                Complete documentation is available at https://www.smartling.com`,
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			configureLoggerVerbose()

			path := cmd.CommandPath()
			isInit = strings.HasPrefix(path, "my-cli init")
			isFiles = strings.HasPrefix(path, "my-cli files")
			isProjects = strings.HasPrefix(path, "my-cli projects")
			isList = strings.HasPrefix(path, "my-cli list")
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 && cmd.Flags().NFlag() == 0 {
				if err := cmd.Help(); err != nil {
					rlog.Error(err.Error())
					os.Exit(1)
				}
				return
			}
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
	rootCmd.PersistentFlags().StringVarP(&directory, "directory", "d", ".", `Sets directory to operate on, usually, to store or to
read files.  Depends on command.`)
	rootCmd.PersistentFlags().Uint32Var(&threads, "threads", 4, `If command can be executed concurrently, it will be
executed for at most <number> of threads.`)
	rootCmd.PersistentFlags().BoolVarP(&insecure, "insecure", "k", false, "Skip HTTPS certificate validation.")
	rootCmd.PersistentFlags().StringVar(&proxy, "proxy", "", "Use specified URL as proxy server.")
	rootCmd.PersistentFlags().StringVar(&smartlingURL, "smartling-url", "", `Specify base Smartling URL, merely for testing
purposes.`)
	rootCmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "Verbose logging")

	return rootCmd, nil
}

// ConfigureLogger initializes the logger with default settings.
func ConfigureLogger() {
	rlog.Init()
	rlog.ToggleRedact(true)
	rlog.SetFormat(lorg.NewFormat("* ${time} ${level:[%s]:right} %s"))
	rlog.SetIndentLines(true)
}

// CLIClientConfig returns a client.Config based on the CLI flags.
func CLIClientConfig() client.Config {
	return client.Config{
		Insecure:     insecure,
		Proxy:        proxy,
		SmartlingURL: smartlingURL,
	}
}

// Config returns a config.Config based on the CLI flags.
func Config() (config.Config, error) {
	params := config.Params{
		Directory:  directory,
		File:       configFile,
		User:       user,
		Secret:     secret,
		Account:    account,
		Project:    project,
		Threads:    threads,
		IsInit:     isInit,
		IsFiles:    isFiles,
		IsProjects: isProjects,
		IsList:     isList,
	}
	cnf, err := config.BuildConfigFromFlags(params)
	if err != nil {
		return config.Config{}, err
	}
	return cnf, nil
}

// Client creates a new Smartling API client based on the configuration and CLI params.
func Client() (sdk.Client, error) {
	cnf, err := Config()
	if err != nil {
		return sdk.Client{}, err
	}
	client, err := client.CreateClient(CLIClientConfig(), cnf, uint8(verbose))
	if err != nil {
		return sdk.Client{}, err
	}
	return *client, nil
}

// ConfigFile returns the path to the configuration file.
func ConfigFile() string {
	return configFile
}

func configureLoggerVerbose() {
	switch verbose {
	case 0:
		// nothing do to

	case 1:
		rlog.SetLevel(lorg.LevelInfo)

	case 2:
		rlog.SetLevel(lorg.LevelDebug)

	default:
		rlog.ToggleRedact(false)
		rlog.SetLevel(lorg.LevelDebug)
	}
}
