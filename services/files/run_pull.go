package files

import (
	"context"
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

	sdk "github.com/Smartling/api-sdk-go"
	sdkfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
	"github.com/reconquest/hierr-go"
	"golang.org/x/sync/errgroup"
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
func (s service) RunPull(ctx context.Context, params PullParams) error {
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
		files, err = globfiles.Remote(ctx, s.APIClient.ListAllFiles, s.Config.ProjectID, params.URI)
		if err != nil {
			return err
		}
	}

	group, groupCtx := errgroup.WithContext(ctx)
	if s.Config.Threads > 0 {
		group.SetLimit(int(s.Config.Threads))
	}

	for _, file := range files {
		group.Go(func() error {
			if err := groupCtx.Err(); err != nil {
				return nil
			}
			if err := s.downloadFileTranslations(groupCtx, params, file); err != nil {
				rlog.Error(err)
			}
			return nil
		})
	}

	_ = group.Wait()

	return nil
}

func (s service) downloadFileTranslations(ctx context.Context, params PullParams, file sdkfile.File) error {
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
	status, err := s.APIClient.GetFileStatus(ctx, projectID, file.FileURI)
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
			map[string]any{
				"FileURI": file.FileURI,
				"Locale":  locale.LocaleID,
			},
		)
		if err != nil {
			return err
		}

		progressPercent, err := locale.ProgressPercent(status.TotalStringCount)
		if err != nil {
			return err
		}
		path = filepath.Join(params.Directory, path)
		if progressThreshold > 0 && progressPercent < int(progressThreshold) {
			fmt.Printf("skipped %s %d%% (threshold: %s%%)\n", path, progressPercent, params.Progress)
			continue
		}

		err = helpers.DownloadFile(
			ctx,
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
