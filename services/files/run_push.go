package files

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	globfiles "github.com/Smartling/smartling-cli/services/helpers/glob_files"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	batchapi "github.com/Smartling/api-sdk-go/api/batches"
	sdktype "github.com/Smartling/api-sdk-go/helpers/file"
	"github.com/reconquest/hierr-go"
)

// PushParams holds the parameters for the RunPush method.
type PushParams struct {
	URI         string
	File        string
	Branch      string
	Locales     []string
	Authorize   bool
	Directory   string
	FileType    string
	Directives  map[string]string
	JobIDOrName string
}

// RunPush uploads files to the Smartling project based on the provided parameters.
func (s service) RunPush(ctx context.Context, params PushParams) error {
	if params.Branch == "@auto" {
		var err error
		params.Branch, err = getGitBranch()
		if err != nil {
			return hierr.Errorf(
				err,
				"unable to autodetect branch name",
			)
		}

		rlog.Infof("autodetected branch name: %s", params.Branch)
	}

	var patterns []string

	if params.File != "" {
		patterns = append(patterns, params.File)
	} else {
		for pattern, section := range s.Config.Files {
			if section.Push.Type != "" {
				patterns = append(patterns, pattern)
			}
		}
	}

	var files []string

	for _, pattern := range patterns {
		base, pattern := globfiles.GetDirectoryFromPattern(pattern)
		chunk, err := globfiles.LocallyFunc(
			params.Directory,
			base,
			pattern,
		)
		if err != nil {
			return clierror.NewError(
				hierr.Errorf(
					err,
					`unable to find matching files to upload`,
				),

				`Check, that specified pattern is valid and refer to help for`+
					` more information about glob patterns.`,
			)
		}

		files = append(files, chunk...)
	}

	if len(files) == 0 {
		return clierror.NewError(
			fmt.Errorf(`no files found by specified patterns`),

			`Check command line pattern if any and configuration file for`+
				` more patterns to search for.`,
		)
	}

	if params.URI != "" && len(files) > 1 {
		return clierror.NewError(
			fmt.Errorf(
				`more than one file is matching speciifed pattern and <uri>`+
					` is specified too`,
			),

			`Either remove <uri> argument or make sure that only one file`+
				` is matching mask.`,
		)
	}

	return s.runPush(ctx, params, files, s.Config.ProjectID)
}

func getGitBranch() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", hierr.Errorf(
			err,
			"unable to get current working directory",
		)
	}

	for {
		if dir == "/" {
			return "", hierr.Errorf(
				err,
				"no git repository can be found containing current directory",
			)
		}

		_, err := os.Stat(filepath.Join(dir, ".git"))
		if err != nil {
			if !os.IsNotExist(err) {
				return "", hierr.Errorf(
					err,
					`unable to get stats for "%s"`,
					dir,
				)
			}

			dir = filepath.Dir(dir)

			continue
		}
		break
	}

	head, err := os.ReadFile(filepath.Join(dir, ".git", "HEAD"))
	if err != nil {
		return "", hierr.Errorf(
			err,
			"unable to read git HEAD",
		)
	}

	return filepath.Base(strings.TrimSpace(string(head))), nil
}

func (s service) runPush(ctx context.Context, params PushParams, files []string, projectID string) error {
	fileUris, err := getFileUris(s.Config.Path, params, files)
	if err != nil {
		return err
	}
	// create new job if params.JobIDOrName is not a valid UUID
	pattern := `^[a-z0-9]{12}$`
	var jobUID string
	if re := regexp.MustCompile(pattern); params.JobIDOrName != "" && re.MatchString(params.JobIDOrName) {
		jobUID = params.JobIDOrName
	}
	var createJobResponse batchapi.CreateJobResponse
	if jobUID == "" {
		timeZoneName, err := timeZoneName()
		if err != nil {
			return err
		}
		nameTemplate := params.JobIDOrName
		if nameTemplate == "" {
			nameTemplate = defaultJobNameTemplate
		}
		payload := batchapi.CreateJobPayload{
			NameTemplate:    nameTemplate,
			Description:     params.JobIDOrName,
			TargetLocaleIds: params.Locales,
			Mode:            batchapi.ReuseExistingMode,
			Salt:            batchapi.RandomAlphanumericSalt,
			TimeZoneName:    timeZoneName,
		}
		createJobResponse, err = s.BatchApi.CreateJob(ctx, projectID, payload)
		if err != nil {
			return err
		}
		jobUID = createJobResponse.TranslationJobUID
	}

	createBatchResponse, err := s.BatchApi.Create(ctx, projectID, batchapi.CreateBatchPayload{
		Authorize:         params.Authorize,
		TranslationJobUID: jobUID,
		FileUris:          fileUris,
	})
	if err != nil {
		return err
	}

	for fileID, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return clierror.UIError{
				Err:       err,
				Operation: "ReadFile",
				Description: `Unable to read file contents.
Check that file exists and readable by current user.`,
				Fields: map[string]string{
					"file": file,
				},
			}
		}
		fileType, found := sdktype.TypeByExt[filepath.Ext(file)]
		if !found {
			rlog.Debugf("unknown file type: %s", file)
		}
		locales := params.Locales
		if len(locales) == 0 {
			locales, err = s.getLocales(projectID)
			if err != nil {
				return err
			}
		}
		payload := batchapi.UploadFilePayload{
			Filename:           fileUris[fileID],
			File:               content,
			FileType:           fileType,
			FileUri:            fileUris[fileID],
			LocalesToAuthorize: locales,
			Directives:         params.Directives,
		}
		uploadFileResponse, err := s.BatchApi.UploadFile(ctx, projectID, createBatchResponse.BatchUID, payload)
		if err != nil {
			return clierror.UIError{
				Err:         err,
				Operation:   "UploadFile",
				Description: fmt.Sprintf(`unable to upload file "%s"`, file),
				Fields: map[string]string{
					"Filename": fileUris[fileID],
					"FileType": fileType.String(),
				},
			}
		}
		rlog.Debugf("uploaded file %v", uploadFileResponse)
		fmt.Printf(
			"%s (%s) %s [code: %s]\n",
			fileUris[fileID],
			payload.FileType,
			"uploaded",
			uploadFileResponse.Code,
		)
	}
	fmt.Println("batch processing is started")
	started := time.Now()
	var processed bool
	for !processed {
		if time.Since(started) > pollingDuration {
			return errors.New("timeout exceeded for polling batch status: " + createBatchResponse.BatchUID)
		}
		time.Sleep(pollingInterval)
		getStatusResponse, err := s.BatchApi.GetStatus(ctx, projectID, createBatchResponse.BatchUID)
		if err != nil {
			return clierror.UIError{
				Err:         err,
				Operation:   "GetStatus",
				Description: `unable to get status for batch`,
				Fields: map[string]string{
					"code": getStatusResponse.Code,
				},
			}
		}
		if strings.ToLower(getStatusResponse.Status) == "completed" {
			processed = true
		}
		errorsInFiles := make(map[string]string)
		for _, file := range getStatusResponse.Files {
			if strings.ToLower(file.Status) == "completed" {
				continue
			}
			if file.Errors != "" && file.Errors != "{}" {
				errorsInFiles[file.FileUri] = file.Errors
			}
		}
		if (getStatusResponse.GeneralErrors != "" && getStatusResponse.GeneralErrors != "{}") || len(errorsInFiles) > 0 {
			return clierror.UIError{
				Err:         errors.New(getStatusResponse.GeneralErrors),
				Operation:   "GetStatus",
				Description: `errors occurred during batch processing`,
				Fields:      errorsInFiles,
			}
		}
	}
	fmt.Println("batch is processed successfully")
	return nil
}

func (s service) getLocales(project string) ([]string, error) {
	var locales []string
	projectDetails, err := s.APIClient.GetProjectDetails(project)
	if err != nil {
		return nil, err
	}
	if projectDetails == nil {
		return nil, fmt.Errorf("no project details found for project: " + project)
	}
	for _, targetLocale := range projectDetails.TargetLocales {
		locales = append(locales, targetLocale.LocaleID)
	}
	return locales, nil
}

func getFileUris(configPath string, params PushParams, files []string) ([]string, error) {
	base, err := filepath.Abs(configPath)
	if err != nil {
		return nil, clierror.NewError(
			hierr.Errorf(
				err,
				`unable to resolve absolute path to config`,
			),

			`It's internal error, please, contact developer for more info`,
		)
	}
	base = filepath.Dir(base)

	branch := params.Branch
	if branch != "" {
		branch = strings.TrimSuffix(params.Branch, "/") + "/"
	}

	res := make([]string, len(files))
	for i, file := range files {
		name, err := filepath.Abs(file)
		if err != nil {
			return nil, clierror.NewError(
				hierr.Errorf(
					err,
					`unable to resolve absolute path to file: %q`,
					file,
				),
				`Check, that file exists and you have proper permissions `+
					`to access it.`,
			)
		}

		if relPath, err := filepath.Rel(base, name); err != nil || strings.HasPrefix(relPath, "..") {
			return nil, clierror.NewError(
				errors.New(
					`you are trying to push file outside project directory`,
				),
				`Check file path and path to configuration file and try again.`,
			)
		}

		name, err = filepath.Rel(base, name)
		if err != nil {
			return nil, clierror.NewError(
				hierr.Errorf(
					err,
					`unable to resolve relative path to file: %q`,
					file,
				),

				`Check, that file exists and you have proper permissions `+
					`to access it.`,
			)
		}

		uri := params.URI
		if uri == "" {
			uri = name
		}

		if uri == "" {
			uri = file
		}
		res[i] = branch + uri
	}
	return res, nil
}

func timeZoneName() (string, error) {
	location := time.Now().Location().String()
	if location != time.Local.String() && strings.ToLower(location) != "" {
		return location, nil
	}
	resp, err := http.Get("https://ipapi.co/json/")
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			rlog.Debugf("failed to close response body: %v", err)
		}
	}()

	type IPInfo struct {
		Timezone string `json:"timezone"`
	}
	var info IPInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", err
	}

	return info.Timezone, nil
}
