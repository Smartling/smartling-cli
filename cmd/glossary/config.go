package glossary

import (
	"os"

	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	"github.com/Smartling/smartling-cli/services/helpers/config"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

// FileConfig defines file config for glossary
type FileConfig struct {
	Glossary struct {
		Export ExportConfig `yaml:"export,omitzero"`
		Create CreateConfig `yaml:"create,omitzero"`
		Import ImportConfig `yaml:"import,omitzero"`
	} `yaml:"glossary,omitzero"`
}

// ImportConfig mirrors the flags accepted by `glossary import`. Each field
// maps to a flag on the command (and to a Smartling Glossary Import API field).
type ImportConfig struct {
	ArchiveMode bool   `yaml:"archive_mode,omitzero"`
	MediaType   string `yaml:"media_type,omitzero"`
}

type ExportConfig struct {
	FileType      string             `yaml:"file_type,omitzero"`
	TbxVersion    string             `yaml:"tbx_version,omitzero"`
	FocusLocaleID string             `yaml:"focus_locale_id,omitzero"`
	LocaleIDs     []string           `yaml:"locale_ids,omitzero"`
	SkipEntries   bool               `yaml:"skip_entries,omitzero"`
	Filter        ExportFilterConfig `yaml:"filter,omitzero"`
}

type ExportFilterConfig struct {
	Query                      string               `yaml:"query,omitzero"`
	LocaleIDs                  []string             `yaml:"locale_ids,omitzero"`
	EntryUIDs                  []string             `yaml:"entry_uids,omitzero"`
	EntryState                 string               `yaml:"entry_state,omitzero"`
	MissingTranslationLocaleID string               `yaml:"missing_translation_locale_id,omitzero"`
	PresentTranslationLocaleID string               `yaml:"present_translation_locale_id,omitzero"`
	DntLocaleID                string               `yaml:"dnt_locale_id,omitzero"`
	ReturnFallbackTranslations bool                 `yaml:"return_fallback_translations,omitzero"`
	LabelsType                 string               `yaml:"labels_type,omitzero"`
	DntTermSet                 bool                 `yaml:"dnt_term_set,omitzero"`
	Created                    CreatedConfig        `yaml:"created,omitzero"`
	CreatedBy                  CreatedByConfig      `yaml:"created_by,omitzero"`
	LastModified               LastModifiedConfig   `yaml:"last_modified,omitzero"`
	LastModifiedBy             LastModifiedByConfig `yaml:"last_modified_by,omitzero"`
	Paging                     PagingConfig         `yaml:"paging,omitzero"`
	Sorting                    SortingConfig        `yaml:"sorting,omitzero"`
}

type CreatedConfig struct {
	Level string `yaml:"level,omitzero"`
	Type  string `yaml:"type,omitzero"`
	Date  string `yaml:"date,omitzero"`
}

type CreatedByConfig struct {
	Level   string   `yaml:"level,omitzero"`
	UserIDs []string `yaml:"user_ids,omitzero"`
}

type LastModifiedConfig struct {
	Level string `yaml:"level,omitzero"`
	Type  string `yaml:"type,omitzero"`
	Date  string `yaml:"date,omitzero"`
}

type LastModifiedByConfig struct {
	Level   string   `yaml:"level,omitzero"`
	UserIDs []string `yaml:"user_ids,omitzero"`
}

type PagingConfig struct {
	Offset int `yaml:"offset,omitzero"`
	Limit  int `yaml:"limit,omitzero"`
}

type SortingConfig struct {
	Field     string `yaml:"field,omitzero"`
	Direction string `yaml:"direction,omitzero"`
	LocaleID  string `yaml:"locale_id,omitzero"`
}

// CreateConfig mirrors the flags accepted by `glossary create`. Each field
// maps to a flag on the command (and to a Smartling Glossary Create API field).
type CreateConfig struct {
	VerificationMode bool     `yaml:"verification_mode,omitzero"`
	LocaleIDs        []string `yaml:"locale_ids,omitzero"`
	// FallbackLocales is a list of "<fallbackLocaleId>:<localeId>[,<localeId>...]"
	// strings — same shape as the --fallback-locale CLI flag.
	FallbackLocales []string `yaml:"fallback_locales,omitzero"`
}

// BindFileConfig binds glossary file config
func BindFileConfig(cmd *cobra.Command) (FileConfig, error) {
	dir := resolve.ConfigDirectory(cmd)
	filename := resolve.ConfigFile(cmd)
	path, err := config.GetPath(dir, filename, false)
	if err != nil {
		return FileConfig{}, err
	}
	var cfg FileConfig
	data, err := os.ReadFile(path)
	if err != nil && os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
