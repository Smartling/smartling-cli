package main

import (
	"crypto/tls"
	"fmt"
	"github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/client"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	smartling "github.com/Smartling/api-sdk-go"
	"github.com/docopt/docopt-go"
	"github.com/kovetskiy/lorg"
	"github.com/reconquest/hierr-go"
)

var version = "1.7"

var usage = `smartling-cli - manage translation files using Smartling.

Add --help option to command to get detailed help for specific command.

Usage:
  smartling-cli [options] [-v]... init --help
  smartling-cli [options] [-v]... init [--dry-run]
  smartling-cli [options] [-v]... projects list --help
  smartling-cli [options] [-v]... projects list [--short]
  smartling-cli [options] [-v]... projects info --help
  smartling-cli [options] [-v]... projects info
  smartling-cli [options] [-v]... projects locales --help
  smartling-cli [options] [-v]... projects locales [--source] [--short] [--format=]
  smartling-cli [options] [-v]... files list --help
  smartling-cli [options] [-v]... files list [--format=] [--short] [<uri>]
  smartling-cli [options] [-v]... files (pull|get) --help
  smartling-cli [options] [-v]... files (pull|get) [--locale=]... [--directory=] [--source] [--format=]
                                               [--progress=] [--retrieve=] [<uri>]
  smartling-cli [options] [-v]... files push --help
  smartling-cli [options] [-v]... files push [(--authorize|--locale=...)] [--branch=] [--type=]
                                         [--directory=] [--directive=]... [<file>] [<uri>]
  smartling-cli [options] [-v]... files rename --help
  smartling-cli [options] [-v]... files rename <old-uri> <new-uri>
  smartling-cli [options] [-v]... files status --help
  smartling-cli [options] [-v]... files status [--directory=] [--format=] [<uri>]
  smartling-cli [options] [-v]... files delete --help
  smartling-cli [options] [-v]... files delete <uri>
  smartling-cli [options] [-v]... files import --help
  smartling-cli [options] [-v]... files import <uri> <file> <locale>
                                           [(--published|--post-translation)]
                                           [--type=] [--overwrite]
  smartling-cli --help

Commands:
  init                    Prepares project to work with Smartling,
                           essentially, assisting user in creating
                           configuration file.
   --dry-run              Do not actually write file, just output it
                           on stdout.
  projects                Used to access various project sub-commands.
   list                   Lists projects for current account.
    -s --short            Display only project IDs.
   info                   Get project details about specific project.
   locales                Display list of target locales.
    -s --short            Display only target locale IDs.
    --format <format>     Use specified format for listing locales.
                           [format: $PROJECTS_LOCALES_FORMAT]
  files                   Used to access various files sub-commands.
   status <uri>           Shows file translation status.
    --format <format>     Specifies format to use for file status output.
                           [default: $FILE_STATUS_FORMAT]
    --directory <dir>     Use another directory as reference to check for
                           local files.
   list <uri>             Lists files from specified project.
    -s --short            Output only file URI.
    --format <format>     Specifies format to use for file list output.
                           [default: $FILE_LIST_FORMAT]
   pull <uri>             Pulls specified files from server.
    --source              Pulls source file as well.
    --progress <done>     Pulls only translations that are at least specified
                           percent of work complete.
    --retrieve <type>     Retrieval type: pending, published, pseudo
                           or contextMatchingInstrumented.
    -d --directory <dir>  Download all files to specified directory.
    --format <format>     Can be used to format path to downloaded files.
                           Note, that single file can be translated in
                           different locales, so format should include locale
                           to create several file paths.
                           [default: $FILE_PULL_FORMAT]
   push <file> <uri>      Uploads specified file into Smartling platform.
    -z --authorize        Automatically authorize all locales in specified
                           file. Incompatible with -l option.
    -l --locale <locale>  Authorize only specified locales.
    -b --branch <branch>  Prepend specified text to the file uri.
    -t --type <type>      Specifies file type which will be used instead of
                           automatically deduced from extension.
    -r --directive <dir>  Specifies one or more directives to use in push
                           request.
   rename <old> <new>     Renames given file by old URI into new URI.
   delete <uri>           Deletes given file from Smartling. This operation
                           can not be undone, so use with care.
   import <uri> <file>    Imports translations for given original file URI with
          <locale>          given locale. Original file mush present on server
                           prior to import.
    --published           Translated content will be published.
    --post-translation    Translated content will be imported into first step
                           of translation. If there are none, it will be
                           published.
    --type <type>         Specify file type. If option is not given, file type
                           will be deduced from extension.
    --overwrite           Overwrite any existing translations.


Options:
  -h --help               Show this help.
  -c --config <file>      Config file in YAML format.
                           By default CLI will look for file named
                           "smartling.yml" in current directory and in all
                           intermediate parents, emulating git behavior.
  -p --project <project>  Project ID to operate on.
                           This option overrides config value "project_id".
  -a --account <account>  Account ID to operate on.
                           This option overrides config value "account_id".
  --user <user>           User ID which will be used for authentication.
                           This option overrides config value "user_id".
  --secret <secret>       Token Secret which will be used for authentication.
                           This option overrides config value "secret".
  -s --short              Use short list output, usually outputs only first
                           column, e.g. file URI in case of files list.
  -l --locale <locale>    Sets locale to filter by or operate upon. Depends on
                           command.
  -d --directory <dir>    Sets directory to operate on, usually, to store or to
                           read files.  Depends on command.  [default: .]
  -z --authorize          Authorize all locales while pushing file.
                           Incompatible with -l option.
  -b --branch <branch>    Prepend specified value to the file URI.
  -t --type <type>        Specify file type. Depends on command.
  -r --directive <dir>    Directives to add to push request in form of
                           <name>=<value>.
  --dry-run               Do not actually perform action, just log it.
  --threads <number>      If command can be executed concurrently, it will be
                           executed for at most <number> of threads.
                           [default: 4]
  -k --insecure           Skip HTTPS certificate validation.
  --proxy <url>           Use specified URL as proxy server.
  --smartling-url <url>   Specify base Smartling URL, merely for testing
                           purposes.
  -v --verbose            Sets verbosity level for logging messages. Specify
                           flag several time to increase verbosity. Useful
                           when debugging and investigating unexpected
                           behavior.
`

var (
	logger = NewRedactedLog()
)

const (
	defaultConfigName = "smartling.yml"

	defaultProjectsLocalesFormat = `{{.LocaleID}}\t{{.Description}}\t{{.Enabled}}\n`
	defaultFilesListFormat       = `{{.FileURI}}\t{{.LastUploaded}}\t{{.FileType}}\n`
	defaultFileStatusFormat      = `{{name .FileURI}}{{with .Locale}}_{{.}}{{end}}{{ext .FileURI}}`
	defaultFilePullFormat        = `{{name .FileURI}}{{with .Locale}}_{{.}}{{end}}{{ext .FileURI}}`
)

func main() {
	usage = os.Expand(usage, func(key string) string {
		switch key {
		case "FILE_LIST_FORMAT":
			return defaultFilesListFormat

		case "FILE_PULL_FORMAT":
			return defaultFilePullFormat

		case "FILE_STATUS_FORMAT":
			return defaultFileStatusFormat

		case "PROJECTS_LOCALES_FORMAT":
			return defaultProjectsLocalesFormat
		}

		return key
	})

	//
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
	//

	args, err := docopt.ParseArgs(usage, nil, "smartling "+version)
	if err != nil {
		panic(err)
	}

	logger.ToggleRedact(true)

	verbose := cmd.Verbose()

	switch verbose {
	case 0:
		// nothing do to

	case 1:
		logger.SetLevel(lorg.LevelInfo)

	case 2:
		logger.SetLevel(lorg.LevelDebug)

	default:
		logger.ToggleRedact(false)
		logger.SetLevel(lorg.LevelDebug)
	}

	logger.SetFormat(lorg.NewFormat("* ${time} ${level:[%s]:right} %s"))
	logger.SetIndentLines(true)

	config, err := buildConfigFromFlags(args)
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	switch {
	case args["init"].(bool):
		err = doInit(config, args)

	case args["projects"].(bool):
		err = doProjects(config, args)

	case args["files"].(bool):
		err = doFiles(config, args)

	default:
		showHelp(args)
	}

	if err != nil {
		reportError(err)
		os.Exit(1)
	}
}

func reportError(err error) {
	switch err := err.(type) {
	case clierror.ProjectNotFoundError, clierror.Error:
		fmt.Fprintln(logger.GetWriter(), err)

	default:
		fmt.Fprintf(logger.GetWriter(), "ERROR: %s", err)
	}
}

func findConfig(name string) (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	var path string

	for {
		path = filepath.Join(dir, name)

		logger.Debugf("looking for config file in: %q", dir)

		_, err = os.Stat(path)
		if err != nil {
			if !os.IsNotExist(err) {
				return "", hierr.Errorf(err, "unable to find config file: %q", path)
			}
		} else {
			logger.Debugf("config file found: %q", path)

			return path, nil
		}

		if dir == "/" {
			break
		}

		dir = filepath.Dir(dir)
	}

	return "", fmt.Errorf(
		"no configuration file %q found",
		name,
	)
}

func buildConfigFromFlags(args map[string]interface{}) (config.Config, error) {
	var (
		directory, _ = args["--directory"].(string)
	)

	var err error

	path, _ := args["--config"].(string)
	if path == "" {
		path, err = findConfig(
			filepath.Join(directory, defaultConfigName),
		)
		if err != nil {
			if !args["init"].(bool) {
				return config.Config{}, clierror.NewError(
					err,

					`Ensure, that config file exists either in the current `+
						`directory or in any parent directory.`,
				)
			} else {
				path = "smartling.yml"
			}
		}
	}

	config, err := config.loadConfigFromFile(path)
	if err != nil {
		return config, clierror.NewError(
			hierr.Errorf(err, `failed to load configuration file "%s".`, path),
			`Check configuration file contents according to documentation.`,
		)
	}

	if config.UserID == "" {
		config.UserID = os.Getenv("SMARTLING_USER_ID")
	}

	if config.Secret == "" {
		config.Secret = os.Getenv("SMARTLING_SECRET")
	}

	if config.ProjectID == "" {
		config.ProjectID = os.Getenv("SMARTLING_PROJECT_ID")
	}

	if args["--user"] != nil {
		config.UserID = args["--user"].(string)
	}

	if args["--secret"] != nil {
		config.Secret = args["--secret"].(string)
	}

	if args["--account"] != nil {
		config.AccountID = args["--account"].(string)
	}

	if args["--project"] != nil {
		config.ProjectID = args["--project"].(string)
	}

	if !args["init"].(bool) {
		if config.UserID == "" {
			return config, MissingConfigValueError{
				ConfigPath: config.path,
				EnvVarName: "SMARTLING_USER_ID",
				ValueName:  "user ID",
				OptionName: "user",
				KeyName:    "user_id",
			}
		}

		if config.Secret == "" {
			return config, MissingConfigValueError{
				ConfigPath: config.path,
				EnvVarName: "SMARTLING_SECRET",
				ValueName:  "token secret",
				OptionName: "secret",
				KeyName:    "secret",
			}
		}
	}

	logger.HideFromConfig(config)

	switch {
	case args["files"].(bool), args["projects"].(bool) && !args["list"].(bool):
		if config.ProjectID == "" {
			return config, MissingConfigValueError{
				ConfigPath: config.path,
				EnvVarName: "SMARTLING_PROJECT_ID",
				ValueName:  "project ID",
				OptionName: "project",
				KeyName:    "project_id",
			}
		}
	}

	threads, err := strconv.ParseInt(args["--threads"].(string), 10, 0)
	if err != nil {
		return config, InvalidConfigValueError{
			ValueName:   "threads",
			Description: "should be valid integer number",
		}
	}

	if threads <= 0 {
		return config, InvalidConfigValueError{
			ValueName:   "threads",
			Description: "should be positive integer number",
		}
	}

	if config.Threads == 0 {
		config.Threads = int(threads)
	}

	return config, nil
}

func createClient(config config.Config, cliClientConfig client.Config) (*smartling.Client, error) {
	client := smartling.NewClient(config.UserID, config.Secret)

	var transport http.Transport

	if cliClientConfig.Insecure {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	if config.Proxy != "" && cliClientConfig.Proxy == "" {
		cliClientConfig.Proxy = config.Proxy
	}

	if cliClientConfig.Proxy != "" {
		proxy, err := url.Parse(cliClientConfig.Proxy)
		if err != nil {
			return nil, clierror.NewError(
				hierr.Errorf(
					err,
					"unable to parse specified proxy URL",
				),

				`Proxy should be valid URL, check CLI options and `+
					`config value.`,
			)
		}

		transport.Proxy = http.ProxyURL(proxy)
	}

	if cliClientConfig.SmartlingURL != "" {
		client.BaseURL = cliClientConfig.SmartlingURL
	}

	client.HTTP.Transport = &transport
	client.UserAgent = "smartling-cli/" + version

	setLogger(client, logger, cmd.Verbose())

	logger.HideRegexp(
		regexp.MustCompile(`"(?:access|refresh)Token": "([^"]+)"`),
	)

	err := client.Authenticate()
	if err != nil {
		return nil, clierror.NewError(
			err,
			`Your credentials are invalid. Double check it and try to run init.\n`,
		)
	}

	return client, nil
}

func doProjects(config config.Config, args map[string]interface{}, cliClientConfig client.Config) error {
	client, err := createClient(config, cliClientConfig)
	if err != nil {
		return err
	}

	switch {
	case args["list"].(bool):
		if config.AccountID == "" {
			return MissingConfigValueError{
				ConfigPath: config.path,
				EnvVarName: "SMARTLING_PROJECT_ID",
				ValueName:  "account ID",
				OptionName: "account",
				KeyName:    "account_id",
			}
		}
	}

	switch {
	case args["list"].(bool):
		return doProjectsList(client, config, args)

	case args["info"].(bool):
		return doProjectsInfo(client, config)

	case args["locales"].(bool):
		return doProjectsLocales(client, config, args)

	}

	return nil
}

func doFiles(config config.Config, args map[string]interface{}, cliClientConfig client.Config) error {
	client, err := createClient(config, cliClientConfig)
	if err != nil {
		return err
	}

	switch {
	case args["list"].(bool):
		return doFilesList(client, config, args)

	case args["pull"].(bool), args["get"].(bool):
		return doFilesPull(client, config, args)

	case args["push"].(bool):
		return doFilesPush(client, config, args)

	case args["status"].(bool):
		return doFilesStatus(client, config, args)

	case args["delete"].(bool):
		return doFilesDelete(client, config, args)

	case args["rename"].(bool):
		return doFilesRename(client, config, args)

	case args["import"].(bool):
		return doFilesImport(client, config, args)
	}

	return nil
}
