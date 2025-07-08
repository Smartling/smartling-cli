package mt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	globfiles "github.com/Smartling/smartling-cli/services/helpers/glob_files"
	"github.com/Smartling/smartling-cli/services/helpers/pointer"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	api "github.com/Smartling/api-sdk-go/api/mt"
)

// TranslateParams is the parameters for the RunTranslate method.
type TranslateParams struct {
	SourceLocale     string
	TargetLocales    []string
	InputDirectory   string
	OutputDirectory  string
	Directives       map[string]string
	Progress         bool
	OverrideFileType string
	AccountUID       api.AccountUID
}

func (s service) RunTranslate(ctx context.Context, p TranslateParams, files []string, updates chan any) ([]TranslateOutput, error) {
	var res []TranslateOutput

	for fileID, file := range files {
		contents, err := getContent(p.InputDirectory, file)
		if err != nil {
			return nil, err
		}
		fileType, found := api.FileTypeByExt[filepath.Ext(file)]
		if !found {
			rlog.Debugf("unknown file type: %s", file)
		}
		request := api.UploadFileRequest{
			File:               contents,
			LocalesToAuthorize: []string{p.SourceLocale},
			FileType:           fileType,
			Directives:         p.Directives,
		}
		uploadFileResponse, err := s.uploader.UploadFile(p.AccountUID, filepath.Base(file), request)
		if err != nil {
			return nil, err
		}

		update := TranslateUpdates{ID: uint32(fileID), Upload: pointer.NewP(true)}
		updates <- update

		params := api.StartParams{
			SourceLocaleIO:  p.SourceLocale,
			TargetLocaleIDs: p.TargetLocales,
		}
		translatorStartResponse, err := s.fileTranslator.Start(p.AccountUID, uploadFileResponse.FileUID, params)
		if err != nil {
			return nil, err
		}

		update.Translate = pointer.NewP("start")
		updates <- update

		if translatorStartResponse.MtUID == "" {
			return nil, clierror.UIError{
				Err:       errors.New("empty mtUid on start translation"),
				Operation: "Start translation",
				Fields: map[string]string{
					"startTranslationCode": translatorStartResponse.Code,
					"file":                 file,
					"uploadCode":           uploadFileResponse.Code,
					"FileUID":              string(uploadFileResponse.FileUID),
				},
				Description: "Translation cannot start. Check if the file is supported and if the source/target locale is valid.",
			}
		}

		var processed bool
		for !processed {
			progressResponse, err := s.fileTranslator.Progress(p.AccountUID, uploadFileResponse.FileUID, translatorStartResponse.MtUID)
			if err != nil {
				return nil, err
			}

			update.Translate = pointer.NewP(progressResponse.State)
			updates <- update

			switch strings.ToUpper(progressResponse.State) {
			case api.QueuedTranslatedState, api.ProcessingTranslatedState:
				time.Sleep(pollingIntervalSeconds)
				continue
			case api.FailedTranslatedState, api.CanceledTranslatedState, api.CompletedTranslatedState:
				processed = true
			default:
				processed = true
			}
			if progressResponse.State != api.CompletedTranslatedState {
				break
			}
			var localeIDs []string
			for _, localeProcessStatus := range progressResponse.LocaleProcessStatuses {
				filename := filepath.Base(file)
				ext := filepath.Ext(filename)
				name := strings.TrimSuffix(filename, ext)
				res = append(res, TranslateOutput{
					File:      filename,
					Locale:    localeProcessStatus.LocaleID,
					Name:      name,
					Ext:       ext,
					Directory: filepath.Dir(file),
				})
				localeIDs = append(localeIDs, localeProcessStatus.LocaleID)
			}

			update.Locale = pointer.NewP(strings.Join(localeIDs, ","))
			updates <- update

			for _, localeProcessStatus := range progressResponse.LocaleProcessStatuses {
				reader, err := s.downloader.File(p.AccountUID, uploadFileResponse.FileUID, translatorStartResponse.MtUID, localeProcessStatus.LocaleID)
				if err != nil {
					return nil, err
				}
				ext := filepath.Ext(file)
				filenameLocale := strings.TrimSuffix(file, ext) + "_" + localeProcessStatus.LocaleID + ext
				outputDirectory, err := filepath.Abs(p.OutputDirectory)
				if err != nil {
					return nil, clierror.UIError{
						Err:         err,
						Operation:   "get absolute output directory",
						Description: "unable to get absolute path for output directory",
						Fields: map[string]string{
							"outputDirectory": p.OutputDirectory,
						},
					}
				}
				if err := os.MkdirAll(outputDirectory, 0755); err != nil {
					return nil, clierror.UIError{
						Err:         err,
						Operation:   "create output directory",
						Description: "unable to create output directory",
						Fields: map[string]string{
							"outputDirectory": outputDirectory,
						},
					}
				}
				if err := saveToFile(reader, filepath.Join(outputDirectory, filepath.Base(filenameLocale))); err != nil {
					return nil, err
				}
			}
			update.Download = pointer.NewP(true)
			updates <- update
		}
	}
	return res, nil
}

func (s service) GetFiles(inputDirectory, fileOrPattern string) ([]string, error) {
	base, pattern := globfiles.GetDirectoryFromPattern(fileOrPattern)
	files, err := globfiles.LocallyFunc(
		inputDirectory,
		base,
		pattern,
	)

	if err != nil {
		return nil, clierror.UIError{
			Err:       err,
			Operation: "globfiles.LocallyFunc",
			Description: `Unable to find matching files to upload.
Check, that specified pattern is valid and refer to help for more information about glob patterns.`,
			Fields: map[string]string{
				"base":    base,
				"pattern": pattern,
			},
		}
	}

	if len(files) == 0 {
		return nil, clierror.UIError{
			Err:       errors.New(`no files found by specified patterns`),
			Operation: "check files",
			Description: `Check command line pattern if any and configuration file for` +
				` more patterns to search for.`,
			Fields: map[string]string{
				"inputDirectory": inputDirectory,
				"fileOrPattern":  fileOrPattern,
			},
		}
	}
	return files, nil
}

// TranslateOutput is translate output
type TranslateOutput struct {
	File      string
	Locale    string
	Name      string
	Ext       string
	Directory string
}

// TranslateUpdates defines updates
type TranslateUpdates struct {
	ID        uint32
	Locale    *string
	Upload    *bool
	Translate *string
	Download  *bool
}

func saveToFile(r io.Reader, filepath string) error {
	outFile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() {
		if err := outFile.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	_, err = io.Copy(outFile, r)
	return err
}
