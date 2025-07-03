package cmd

import (
	"github.com/Smartling/smartling-cli/services/helpers/client"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	sdk "github.com/Smartling/api-sdk-go"
	"github.com/kovetskiy/lorg"
)

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
		Directory:  operationDirectory,
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
func Client() (sdk.HttpAPIClient, error) {
	cnf, err := Config()
	if err != nil {
		return sdk.HttpAPIClient{}, err
	}
	client, err := client.CreateClient(CLIClientConfig(), cnf, uint8(verbose))
	if err != nil {
		return sdk.HttpAPIClient{}, err
	}
	return client, nil
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
