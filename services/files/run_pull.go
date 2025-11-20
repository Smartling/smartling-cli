package files

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Smartling/smartling-cli/services/helpers"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	"github.com/Smartling/smartling-cli/services/helpers/format"
	globfiles "github.com/Smartling/smartling-cli/services/helpers/glob_files"
	"github.com/Smartling/smartling-cli/services/helpers/reader"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	threadpool "github.com/Smartling/smartling-cli/services/helpers/thread_pool"

	sdk "github.com/Smartling/api-sdk-go"
	sdkfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
	"github.com/reconquest/hierr-go"
)

// PullParams is the parameters for the RunPull method.
type PullParams struct {
	URI       string
	All       bool
	Format    string
	Directory string
	Source    bool
	Locales   []string
	Progress  string
	Retrieve  string
}

func (p PullParams) validate() error {
	if p.URI == "" && !p.All {
		return fmt.Errorf("either uri or --all is required")
	}
	if p.All && p.URI != "" {
		return clierror.ErrIncompatibleParams("all", []string{"uri"})
	}
	return nil
}

// RunPull pulls translations for files from the Smartling based on the provided parameters.
func (s service) RunPull(params PullParams) error {
	if err := params.validate(); err != nil {
		return err
	}
	if params.Format == "" {
		params.Format = format.DefaultFilePullFormat
	}

	var (
		err   error
		files []sdkfile.File
	)
	if params.URI == "-" {
		files, err = reader.ReadFilesFromStdin()
		if err != nil {
			return err
		}
	} else {
		files, err = globfiles.Remote(s.APIClient.ListAllFiles, s.Config.ProjectID, params.URI)
		if err != nil {
			return err
		}
	}

	pool := threadpool.NewThreadPool(s.Config.Threads)

	for _, file := range files {
		// func closure required to pass different file objects to goroutines
		func(file sdkfile.File) {
			pool.Do(func() {
				err := s.downloadFileTranslations(params, file)
				if err != nil {
					rlog.Error(err)
				}
			})
		}(file)
	}

	pool.Wait()

	return nil
}

func (s service) downloadFileTranslations(params PullParams, file sdkfile.File) error {
	progress := strings.TrimSpace(params.Progress)
	progress = strings.TrimSpace(strings.TrimSuffix(progress, "%"))
	if progress == "" {
		progress = "0"
	}
	progressThreshold, err := strconv.ParseInt(progress, 10, 0)
	if err != nil {
		return hierr.Errorf(
			err,
			"unable to parse --progress as integer",
		)
	}

	retrievalType := sdk.RetrievalType(params.Retrieve)

	if params.Format == "" {
		params.Format = format.DefaultFileStatusFormat
	}

	projectID := s.Config.ProjectID
	status, err := s.APIClient.GetFileStatus(projectID, file.FileURI)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to retrieve file "%s" locales from project "%s"`,
			file.FileURI,
			projectID,
		)
	}

	var translations []sdkfile.FileStatusTranslation

	if params.Source {
		translations = []sdkfile.FileStatusTranslation{
			{LocaleID: ""},
		}
	} else {
		translations = status.Items
	}

	for _, locale := range translations {
		if len(params.Locales) > 0 {
			if !hasLocaleInList(locale.LocaleID, params.Locales) {
				continue
			}
		}

		useFormat := format.UsePullFormat
		if params.Format != "" {
			useFormat = func(_ config.FileConfig) string {
				return params.Format
			}
		}

		path, err := format.ExecuteFileFormat(
			s.Config,
			file,
			params.Format,
			useFormat,
			map[string]interface{}{
				"FileURI": file.FileURI,
				"Locale":  locale.LocaleID,
			},
		)
		if err != nil {
			return err
		}

		progressPercent := locale.ProgressPercent()
		if progressThreshold > 0 && progressPercent < int(progressThreshold) {
			fmt.Printf("skipped %s %s%% (threshold: %s%%)\n", path, progressPercent, params.Progress)
			continue
		}

		path = filepath.Join(params.Directory, path)

		err = helpers.DownloadFile(
			s.APIClient,
			projectID,
			file,
			locale.LocaleID,
			path,
			retrievalType,
		)
		if err != nil {
			return err
		}

		if params.Source {
			fmt.Printf("downloaded %s\n", path)
		} else {
			fmt.Printf("downloaded %s %d%%\n", path, progressPercent)
		}
	}

	return err
}

func hasLocaleInList(locale string, locales []string) bool {
	for _, filter := range locales {
		if strings.EqualFold(strings.ToLower(filter), strings.ToLower(locale)) {
			return true
		}
	}

	return false
}
