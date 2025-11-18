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
	smfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
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

func (s service) RunTranslate(ctx context.Context, params TranslateParams, files []string, updates chan any) ([]TranslateOutput, error) {
	var res []TranslateOutput

	for fileID, file := range files {
		rlog.Debugf("Running translate for file %s", file)
		contents, err := getContent(params.InputDirectory, file)
		if err != nil {
			return nil, err
		}
		var fileType api.Type
		var found bool
		if params.OverrideFileType != "" {
			fileType, found = smfile.ParseType(api.FirstType, api.LastType, params.OverrideFileType)
			if !found {
				rlog.Debugf("unknown override file type: %s", params.OverrideFileType)
			}
		}
		if !found {
			var found bool
			fileType, found = api.TypeByExt[filepath.Ext(file)]
			if !found {
				rlog.Debugf("unknown file type for file: %s", file)
			}
		}
		request := api.UploadFileRequest{
			File:               contents,
			LocalesToAuthorize: []string{params.SourceLocale},
			FileType:           fileType,
			Directives:         params.Directives,
		}
		rlog.Debugf("start upload")
		uploadFileResponse, err := s.uploader.UploadFile(params.AccountUID, filepath.Base(file), request)
		if err != nil {
			return nil, err
		}
		rlog.Debugf("finish upload")

		update := TranslateUpdates{ID: uint32(fileID * len(params.TargetLocales)), Upload: pointer.NewP(true)}
		updates <- update

		if params.SourceLocale == "" {
			rlog.Debugf("detect language")
			detectFileLanguageResponse, err := s.translationControl.DetectFileLanguage(params.AccountUID, uploadFileResponse.FileUID)
			if err != nil {
				return nil, err
			}
			started := time.Now()
			var processed bool
			for !processed {
				if time.Since(started) > pollingDuration {
					return nil, errors.New("timeout exceeded for polling detect file language progress: FileUID:" + string(uploadFileResponse.FileUID))
				}
				rlog.Debugf("check detection progress")
				detectionProgressResponse, err := s.translationControl.DetectionProgress(params.AccountUID, uploadFileResponse.FileUID, detectFileLanguageResponse.LanguageDetectionUID)
				if err != nil {
					return nil, err
				}

				rlog.Debugf("detection progress state: %s", detectionProgressResponse.State)
				switch strings.ToUpper(detectionProgressResponse.State) {
				case api.QueuedTranslatedState, api.ProcessingTranslatedState:
					time.Sleep(pollingInterval)
					continue
				case api.FailedTranslatedState, api.CanceledTranslatedState, api.CompletedTranslatedState:
					processed = true
				default:
					processed = true
				}
				if detectionProgressResponse.State != api.CompletedTranslatedState {
					rlog.Debugf("detection progress break on incomplete state: %s", detectionProgressResponse.State)
					break
				}
				if len(detectionProgressResponse.DetectedSourceLanguages) > 0 {
					params.SourceLocale = detectionProgressResponse.DetectedSourceLanguages[0].LanguageID
				}
			}
		}

		startParams := api.StartParams{
			SourceLocaleID:  params.SourceLocale,
			TargetLocaleIDs: params.TargetLocales,
		}
		rlog.Debugf("start translation")
		translatorStartResponse, err := s.fileTranslator.Start(params.AccountUID, uploadFileResponse.FileUID, startParams)
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

		started := time.Now()
		var processed bool
		for !processed {
			if time.Since(started) > pollingDuration {
				return nil, errors.New("timeout exceeded for polling file translation progress FileUID:" + string(uploadFileResponse.FileUID))
			}
			rlog.Debugf("check translation progress")
			progressResponse, err := s.fileTranslator.Progress(params.AccountUID, uploadFileResponse.FileUID, translatorStartResponse.MtUID)
			if err != nil {
				return nil, err
			}

			update.Translate = pointer.NewP(progressResponse.State)
			updates <- update

			rlog.Debugf("progress state: %s", progressResponse.State)
			switch strings.ToUpper(progressResponse.State) {
			case api.QueuedTranslatedState, api.ProcessingTranslatedState:
				time.Sleep(pollingInterval)
				continue
			case api.FailedTranslatedState, api.CanceledTranslatedState, api.CompletedTranslatedState:
				processed = true
			default:
				processed = true
			}
			if progressResponse.State != api.CompletedTranslatedState {
				break
			}
			for updateID, localeProcessStatus := range progressResponse.LocaleProcessStatuses {
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
				update.ID = uint32(fileID*len(params.TargetLocales) + updateID)
				update.Locale = pointer.NewP(localeProcessStatus.LocaleID)
				updates <- update

				rlog.Debugf("download start")
				reader, err := s.downloader.File(params.AccountUID, uploadFileResponse.FileUID, translatorStartResponse.MtUID, localeProcessStatus.LocaleID)
				if err != nil {
					return nil, err
				}
				rlog.Debugf("download finished")
				filenameLocale := strings.TrimSuffix(file, ext) + "_" + localeProcessStatus.LocaleID + ext
				outputDirectory, err := filepath.Abs(params.OutputDirectory)
				if err != nil {
					return nil, clierror.UIError{
						Err:         err,
						Operation:   "get absolute output directory",
						Description: "unable to get absolute path for output directory",
						Fields: map[string]string{
							"outputDirectory": params.OutputDirectory,
						},
					}
				}
				if err := os.MkdirAll(outputDirectory, 0o755); err != nil {
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
				update.ID = uint32(fileID*len(params.TargetLocales) + updateID)
				update.TranslatedFile = pointer.NewP(filepath.Base(filenameLocale))
				update.Download = pointer.NewP(true)
				updates <- update
			}
		}
	}
	return res, nil
}

func (s service) GetFiles(inputDirectory, fileOrPattern string) ([]string, error) {
	rlog.Debugf("get files")
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
	ID             uint32
	Locale         *string
	Upload         *bool
	Translate      *string
	TranslatedFile *string
	Download       *bool
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
