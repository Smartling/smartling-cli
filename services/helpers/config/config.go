package config

import (
	"fmt"
	"os"

	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"

	"dario.cat/mergo"
	"github.com/gobwas/glob"
	"github.com/goccy/go-yaml"
	"github.com/reconquest/hierr-go"
)

// Source identifies where a configuration value was resolved from.
type Source string

// Source values, ordered by precedence (lowest first).
const (
	SourceDefault Source = "default"
	SourceConfig  Source = "config"
	SourceEnv     Source = "env"
	SourceFlag    Source = "flag"
)

// Sources records the origin of each configuration value that supports
// flag/env/file precedence. It is populated by BuildConfigFromFlags and is
// intentionally not serialized to YAML.
type Sources struct {
	UserID    Source
	AccountID Source
	ProjectID Source
}

func (s Sources) String() string {
	return fmt.Sprintf(
		"project=%s  account=%s  user=%s",
		s.ProjectID,
		s.AccountID,
		s.UserID,
	)
}

// FileConfig is the configuration from file.
type FileConfig struct {
	Pull struct {
		Format string `yaml:"format,omitzero"`
	} `yaml:"pull,omitzero"`
	Push struct {
		Type       string            `yaml:"type,omitzero"`
		Directives map[string]string `yaml:"directives,omitzero,flow"`
	} `yaml:"push,omitzero"`
}

// Config is the configuration for the Smartling CLI.
type Config struct {
	UserID    string `yaml:"user_id"`
	Secret    string `yaml:"secret"`
	AccountID string `yaml:"account_id"`
	ProjectID string `yaml:"project_id,omitzero"`
	Threads   uint32 `yaml:"threads,omitzero"`

	Files map[string]FileConfig `yaml:"files"`

	Proxy string `yaml:"proxy,omitzero"`

	Path string `yaml:"-"`

	Sources Sources `yaml:"-"`
}

// GetFileConfig returns the FileConfig for the given path.
func (config *Config) GetFileConfig(path string) (FileConfig, error) {
	var (
		match FileConfig
		found bool
	)

	for key, candidate := range config.Files {
		pattern, err := glob.Compile(key, '/')
		if err != nil {
			return FileConfig{}, clierror.NewError(
				hierr.Errorf(
					err,
					`unable to compile pattern from config file (key "%s")`,
					key,
				),

				`File match pattern is malformed. Check out help for more `+
					`information on globbing patterns.`,
			)
		}

		if pattern.Match(path) {
			match = candidate
			found = true
		}
	}

	defaults := config.Files["default"]

	if !found {
		return defaults, nil
	}

	err := mergo.Merge(&match, defaults)
	if err != nil {
		return FileConfig{}, clierror.NewError(
			hierr.Errorf(err, "unable to merge file config options"),
			`It's internal error. Consider reporting bug.`,
		)
	}

	return match, nil
}

// LoadConfigFromFile loads the configuration from the specified file.
func LoadConfigFromFile(filename string) (Config, error) {
	config := Config{
		Path: filename,
	}

	data, err := os.ReadFile(filename)
	if err != nil && os.IsNotExist(err) {
		return config, nil
	}
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, err
	}

	return config, nil
}
