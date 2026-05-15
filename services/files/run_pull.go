package files

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/Smartling/smartling-cli/services/helpers"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	"github.com/Smartling/smartling-cli/services/helpers/format"
	globfiles "github.com/Smartling/smartling-cli/services/helpers/glob_files"
	"github.com/Smartling/smartling-cli/services/helpers/reader"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	sdk "github.com/Smartling/api-sdk-go"
	sdkjob "github.com/Smartling/api-sdk-go/api/job"
	sdkfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
	"github.com/gobwas/glob"
	"github.com/reconquest/hierr-go"
	"golang.org/x/sync/errgroup"
)

// PullParams is the parameters for the RunPull method.
type PullParams struct {
	URI       string
	JobUID    string
	All       bool
	Format    string
	Directory string
	Source    bool
	Locales   []string
	Resume    bool
	DryRun    bool
	Progress  string
	Retrieve  string
}

func (p *PullParams) setDefaultFormatIfEmpty() {
	if p.Format != "" {
		return
	}
	if p.JobUID != "" {
		p.Format = format.DefaultFilePullJobFormat
		return
	}
	p.Format = format.DefaultFilePullFormat
}

func (p *PullParams) validate() error {
	if p.URI == "" && p.JobUID == "" && !p.All {
		return fmt.Errorf("either uri or --job-uid or --all is required")
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
	params.setDefaultFormatIfEmpty()

	var (
		err        error
		files      []sdkfile.File
		jobLocales []string
	)
	switch {
	case params.JobUID != "":
		files, jobLocales, err = s.enumerateJobFiles(ctx, params.JobUID)
	case params.URI == "-":
		files, err = reader.ReadFilesFromStdin()
	default:
		files, err = globfiles.Remote(ctx, s.APIClient.ListAllFiles, s.Config.ProjectID, params.URI)
	}
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no files found matching the provided parameters")
	}

	// When --job-uid is combined with a URI glob, filter the job's file list
	// down to URIs that match the pattern.
	if params.JobUID != "" && params.URI != "" && params.URI != "-" {
		filtered, err := filterFilesByGlob(files, params.URI)
		if err != nil {
			return err
		}
		if len(filtered) == 0 {
			return fmt.Errorf("job %q has no files matching uri pattern %q", params.JobUID, params.URI)
		}
		files = filtered
	}

	if params.JobUID != "" {
		if len(jobLocales) == 0 {
			return fmt.Errorf("job %q has no target locales; nothing to download", params.JobUID)
		}
		params.Locales = filterLocales(jobLocales, params.Locales)
		if len(params.Locales) == 0 {
			return fmt.Errorf("job %q has no target locales matching the requested --locale filters", params.JobUID)
		}
	}

	if params.DryRun {
		return s.printDryRun(files, params)
	}

	group, groupCtx := errgroup.WithContext(ctx)
	if s.Config.Threads > 0 {
		group.SetLimit(int(s.Config.Threads))
	}
	var failed atomic.Int32
	for _, file := range files {
		group.Go(func() error {
			if err := groupCtx.Err(); err != nil {
				return nil
			}
			if err := s.downloadFileTranslations(groupCtx, params, file); err != nil {
				failed.Add(1)
				rlog.Error(err)
			}
			return nil
		})
	}
	_ = group.Wait()
	if n := failed.Load(); n > 0 {
		return fmt.Errorf("%d file(s) failed to download; see log for details", n)
	}
	return nil
}

// printDryRun writes the resolved file × locale matrix to stdout without
// calling GetFileStatus or downloading anything.
func (s service) printDryRun(files []sdkfile.File, params PullParams) error {
	for _, file := range files {
		locales := params.Locales
		if params.Source {
			locales = append([]string{""}, locales...)
		}
		for _, locale := range locales {
			path, err := s.renderPullPath(file, locale, params)
			if err != nil {
				return err
			}
			fmt.Println(filepath.Join(params.Directory, path))
		}
	}
	return nil
}

// renderPullPath produces the on-disk relative path for a file/locale pair
// using the pull format template, including the JobUID variable.
func (s service) renderPullPath(file sdkfile.File, locale string, params PullParams) (string, error) {
	useFormat := format.UsePullFormat
	if params.Format != "" {
		useFormat = func(_ config.FileConfig) string {
			return params.Format
		}
	}
	return format.ExecuteFileFormat(
		s.Config,
		file,
		params.Format,
		useFormat,
		map[string]any{
			"FileURI": file.FileURI,
			"Locale":  locale,
			"JobUID":  params.JobUID,
		},
	)
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

		path, err := s.renderPullPath(file, locale.LocaleID, params)
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

		if params.Resume {
			if _, err := os.Stat(path); err == nil {
				fmt.Printf("skipped %s (already exists)\n", path)
				continue
			}
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

// enumerateJobFiles resolves the file × target-locale matrix for a job by
// calling the Jobs API.
func (s service) enumerateJobFiles(ctx context.Context, jobUID string) ([]sdkfile.File, []string, error) {
	var (
		job      sdkjob.GetJobResponse
		jobFiles []sdkjob.JobFile
		group    errgroup.Group
	)
	projectID := s.Config.ProjectID
	group.Go(func() error {
		var err error
		job, err = s.JobApi.GetJob(ctx, projectID, jobUID)
		return err
	})
	group.Go(func() error {
		var err error
		jobFiles, err = s.JobApi.ListFiles(ctx, projectID, jobUID)
		return err
	})
	if err := group.Wait(); err != nil {
		if errors.Is(err, sdkjob.ErrNotFound) {
			return nil, nil, fmt.Errorf("job %q not found in project %q", jobUID, projectID)
		}
		return nil, nil, fmt.Errorf("unable to fetch job %q in project %q: %w", jobUID, projectID, err)
	}

	files := make([]sdkfile.File, 0, len(jobFiles))
	for _, jf := range jobFiles {
		files = append(files, sdkfile.File{FileURI: jf.FileURI})
	}
	return files, job.TargetLocaleIDs, nil
}

// filterFilesByGlob keeps only files whose FileURI matches the provided glob
// pattern. Uses the same gobwas/glob delimiter as globfiles.Remote so the
// pattern behavior is identical for both code paths.
func filterFilesByGlob(files []sdkfile.File, uri string) ([]sdkfile.File, error) {
	pattern, err := glob.Compile(uri, '/')
	if err != nil {
		return nil, fmt.Errorf("invalid uri glob pattern %q: %w", uri, err)
	}
	out := make([]sdkfile.File, 0, len(files))
	for _, f := range files {
		if pattern.Match(f.FileURI) {
			out = append(out, f)
		}
	}
	return out, nil
}

// filterLocales returns the subset of locales (preserving order) that
// also appears in filter, matched case-insensitively. If filter is
// empty, locales is returned unchanged.
func filterLocales(locales, filter []string) []string {
	if len(filter) == 0 {
		return locales
	}
	var res []string
	for _, locale := range locales {
		if hasLocaleInList(locale, filter) {
			res = append(res, locale)
		}
	}
	return res
}
