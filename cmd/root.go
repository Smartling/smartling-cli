package cmd

import (
	"github.com/Smartling/smartling-cli/cmd/files"
	"github.com/Smartling/smartling-cli/cmd/init"
	"github.com/Smartling/smartling-cli/services/helpers/client"
	"github.com/spf13/cobra"
)

var (
	verbose      uint8
	config       string
	project      string
	account      string
	user         string
	secret       string
	short        string
	locale       string
	directory    string
	authorize    string
	branch       string
	typ          string
	directive    string
	dryRun       string
	threads      uint32
	insecure     bool
	proxy        string
	smartlingURL string
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "smartling-cli",
		Short:   "Manage translation files using Smartling CLI.",
		Version: "1.7",
		Long: `Manage translation files using Smartling CLI.
                Complete documentation is available at https://www.smartling.com`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "", `Config file in YAML format.
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
	rootCmd.PersistentFlags().StringVarP(&short, "short", "s", "", `Use short list output, usually outputs only first
column, e.g. file URI in case of files list.`)
	rootCmd.PersistentFlags().StringVarP(&locale, "locale", "l", "", "Sets locale to filter by or operate upon. Depends on command.")
	rootCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "", `Sets directory to operate on, usually, to store or to
read files.  Depends on command.  [default: .]`)
	rootCmd.PersistentFlags().StringVarP(&authorize, "authorize", "z", "", `Authorize all locales while pushing file.
Incompatible with -l option.`)
	rootCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "", "Prepend specified value to the file URI.")
	rootCmd.PersistentFlags().StringVarP(&typ, "type", "t", "", "Specify file type. Depends on command.")
	rootCmd.PersistentFlags().StringVarP(&directive, "directive", "r", "", `Directives to add to push request in form of
<name>=<value>.`)
	rootCmd.PersistentFlags().StringVar(&dryRun, "dry-run", "", "Do not actually perform action, just log it.")
	rootCmd.PersistentFlags().Uint32Var(&threads, "threads", 0, `If command can be executed concurrently, it will be
executed for at most <number> of threads.
[default: 4]`)
	rootCmd.PersistentFlags().BoolVarP(&insecure, "insecure", "k", false, "Skip HTTPS certificate validation.")
	rootCmd.PersistentFlags().StringVar(&proxy, "proxy", "", "Use specified URL as proxy server.")
	rootCmd.PersistentFlags().StringVar(&smartlingURL, "smartling-url", "", `Specify base Smartling URL, merely for testing
purposes.`)
	rootCmd.PersistentFlags().Uint8VarP(&verbose, "verbose", "v", 0, "Verbose logging")

	rootCmd.AddCommand(init.NewInitCmd())
	rootCmd.AddCommand(files.NewFilesCmd())

	return rootCmd
}

func Verbose() uint8 {
	return verbose
}

func ClientConfig() client.Config {
	return client.Config{
		Insecure:     insecure,
		Proxy:        proxy,
		SmartlingURL: smartlingURL,
	}
}
