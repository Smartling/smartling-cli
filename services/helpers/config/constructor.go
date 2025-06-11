package config

import (
	"fmt"
	"os"
	"path/filepath"

	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/reconquest/hierr-go"
)

const (
	defaultConfigName = "smartling.yml"
)

// Params is parameters for building a configuration object.
type Params struct {
	Directory  string
	File       string
	User       string
	Secret     string
	Account    string
	Project    string
	Threads    uint32
	IsInit     bool
	IsFiles    bool
	IsProjects bool
	IsList     bool
}

// BuildConfigFromFlags returns a Config object based on the provided parameters,
// and an error if any.
func BuildConfigFromFlags(params Params) (Config, error) {
	var err error

	path := params.File
	if path == "" {
		path, err = findConfig(
			filepath.Join(params.Directory, defaultConfigName),
		)
		if err != nil {
			if !params.IsInit {
				return Config{}, clierror.NewError(
					err,

					`Ensure, that config file exists either in the current `+
						`directory or in any parent directory.`,
				)
			}
			path = "smartling.yml"
		}
	}

	config, err := LoadConfigFromFile(path)
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

	if params.User != "" {
		config.UserID = params.User
	}

	if params.Secret != "" {
		config.Secret = params.Secret
	}

	if params.Account != "" {
		config.AccountID = params.Account
	}

	if params.Project != "" {
		config.ProjectID = params.Project
	}

	if !params.IsInit {
		if config.UserID == "" {
			return config, clierror.MissingConfigValueError{
				ConfigPath: config.Path,
				EnvVarName: "SMARTLING_USER_ID",
				ValueName:  "user ID",
				OptionName: "user",
				KeyName:    "user_id",
			}
		}

		if config.Secret == "" {
			return config, clierror.MissingConfigValueError{
				ConfigPath: config.Path,
				EnvVarName: "SMARTLING_SECRET",
				ValueName:  "token secret",
				OptionName: "secret",
				KeyName:    "secret",
			}
		}
	}

	rlog.HideString(config.Secret)
	rlog.HideString(config.UserID)
	rlog.HideString(config.AccountID)
	rlog.HideString(config.ProjectID)

	switch {
	case params.IsFiles, params.IsProjects && !params.IsList:
		if config.ProjectID == "" {
			return config, clierror.MissingConfigValueError{
				ConfigPath: config.Path,
				EnvVarName: "SMARTLING_PROJECT_ID",
				ValueName:  "project ID",
				OptionName: "project",
				KeyName:    "project_id",
			}
		}
	}

	if config.Threads == 0 {
		config.Threads = params.Threads
	}

	return config, nil
}

func findConfig(name string) (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	fmt.Println("dir = ", dir)

	var path string

	for {
		path = filepath.Join(dir, name)

		rlog.Debugf("looking for config file in: %q", dir)

		_, err = os.Stat(path)
		if err != nil {
			if !os.IsNotExist(err) {
				return "", hierr.Errorf(err, "unable to find config file: %q", path)
			}
		} else {
			rlog.Debugf("config file found: %q", path)

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
