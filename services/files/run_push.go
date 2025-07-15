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
	sdkerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
	sdkfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
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
	Directives  []string
	JobIDOrName string
}

// RunPush uploads files to the Smartling project based on the provided parameters.
func (s service) RunPush(ctx context.Context, params PushParams) error {
	var (
		failedFiles []string
		project     = s.Config.ProjectID
		result      error
	)

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

	if params.JobIDOrName != "" {
		return s.runPushWithJob(ctx, params, files, project)
	}

	fileUris, err := getFileUris(s.Config.Path, params, files)
	if err != nil {
		return err
	}

	for fileID, file := range files {
		fileConfig, err := s.Config.GetFileConfig(file)
		if err != nil {
			return clierror.NewError(
				hierr.Errorf(
					err,
					`unable to retrieve file specific configuration`,
				),

				``,
			)
		}

		contents, err := os.ReadFile(file)
		if err != nil {
			return clierror.NewError(
				hierr.Errorf(
					err,
					`unable to read file contents "%s"`,
					file,
				),

				`Check that file exists and readable by current user.`,
			)
		}

		request := sdkfile.FileUploadRequest{
			File:               contents,
			Authorize:          params.Authorize,
			LocalesToAuthorize: params.Locales,
		}

		request.FileURI = fileUris[fileID]

		if fileConfig.Push.Type == "" {
			if params.FileType == "" {
				request.FileType = sdkfile.GetFileTypeByExtension(
					filepath.Ext(file),
				)

				if request.FileType == sdkfile.FileTypeUnknown {
					return clierror.NewError(
						fmt.Errorf(
							"unable to deduce file type from extension: %q",
							filepath.Ext(file),
						),

						`You need to specify file type via --type option.`,
					)
				}
			} else {
				request.FileType = sdkfile.FileType(params.FileType)
			}
		} else {
			request.FileType = sdkfile.FileType(fileConfig.Push.Type)
		}

		request.Smartling.Directives = fileConfig.Push.Directives

		for _, directive := range params.Directives {
			spec := strings.SplitN(directive, "=", 2)
			if len(spec) != 2 {
				return clierror.NewError(
					fmt.Errorf(
						"invalid directive specification: %q",
						directive,
					),

					`Should be in the form of <name>=<value>.`,
				)
			}

			if request.Smartling.Directives == nil {
				request.Smartling.Directives = map[string]string{}
			}

			request.Smartling.Directives[spec[0]] = spec[1]
		}

		response, err := s.APIClient.UploadFile(project, request)

		if err != nil {
			if returnError(err) {
				return clierror.NewError(
					err,
					fmt.Sprintf(`unable to upload file "%s"`, file),
					`Check, that you have enough permissions to upload file to`+
						` the specified project`,
				)
			}
			_, _ = fmt.Fprintln(os.Stderr, "Unable to upload file "+file)
			failedFiles = append(failedFiles, file)
		} else {
			status := "new"
			if response.Overwritten {
				status = "overwritten"
			}

			fmt.Printf(
				"%s (%s) %s [%d strings %d words]\n",
				fileUris[fileID],
				request.FileType,
				status,
				response.StringCount,
				response.WordCount,
			)
		}
	}

	if len(failedFiles) != 0 {
		result = clierror.NewError(fmt.Errorf("failed to upload %d files", len(failedFiles)), "failed to upload files "+strings.Join(failedFiles, ", "))
	}

	return result
}

func returnError(err error) bool {
	if errors.Is(err, sdkerror.NotAuthorizedError{}) {
		return true
	}

	for {
		smartlingAPIError, isSmartlingAPIError := err.(sdkerror.APIError)
		if isSmartlingAPIError {
			reasons := map[string]struct{}{
				"AUTHENTICATION_ERROR":   {},
				"AUTHORIZATION_ERROR":    {},
				"MAINTENANCE_MODE_ERROR": {},
			}

			_, stopExecution := reasons[smartlingAPIError.Code]
			return stopExecution
		}
		if err = errors.Unwrap(err); err == nil {
			return false
		}
	}
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

func (s service) runPushWithJob(ctx context.Context, params PushParams, files []string, project string) error {
	fileUris, err := getFileUris(s.Config.Path, params, files)
	if err != nil {
		return err
	}
	// create new job if params.JobIDOrName is not a valid UUID
	pattern := `^[a-z0-9]{12}$`
	var jobUID string
	if re := regexp.MustCompile(pattern); re.MatchString(params.JobIDOrName) {
		jobUID = params.JobIDOrName
	}
	var createJobResponse batchapi.CreateJobResponse
	if jobUID == "" {
		timeZoneName, err := timeZoneName()
		if err != nil {
			return err
		}
		payload := batchapi.CreateJobPayload{
			NameTemplate: params.JobIDOrName,
			Description:  params.JobIDOrName,
			//TargetLocaleIds: params.Locales,
			Mode:         batchapi.ReuseExistingMode,
			Salt:         batchapi.RandomAlphanumericSalt,
			TimeZoneName: timeZoneName,
		}
		createJobResponse, err = s.BatchApi.CreateJob(ctx, project, payload)
		if err != nil {
			return err
		}
		jobUID = createJobResponse.TranslationJobUID
	}

	createBatchResponse, err := s.BatchApi.Create(ctx, project, batchapi.CreateBatchPayload{
		Authorize:         true,
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
		payload := batchapi.UploadFilePayload{
			Filename:           fileUris[fileID],
			File:               content,
			FileType:           fileType,
			LocalesToAuthorize: params.Locales,
		}
		uploadFileResponse, err := s.BatchApi.UploadFile(ctx, project, createBatchResponse.BatchUID, payload)
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
		var processed bool
		for !processed {
			time.Sleep(pollingInterval)
			getStatusResponse, err := s.BatchApi.GetStatus(ctx, project, uploadFileResponse.Code)
			if err != nil {
				return clierror.UIError{
					Err:         err,
					Operation:   "GetStatus",
					Description: `unable to get status for file`,
					Fields: map[string]string{
						"upload file code": uploadFileResponse.Code,
						"code":             getStatusResponse.Code,
					},
				}
			}
			switch strings.ToLower(getStatusResponse.Status) {
			case "complete", "success":
				processed = true
				// TODO add failed statuses
			}
			errorsInFiles := make(map[string]string)
			for _, file := range getStatusResponse.Files {
				if file.Errors != "" {
					errorsInFiles[file.FileUri] = file.Errors
				}
			}
			if getStatusResponse.GeneralErrors != "" || len(errorsInFiles) > 0 {
				return clierror.UIError{
					Err:         errors.New(getStatusResponse.GeneralErrors),
					Operation:   "GetStatus",
					Description: `upload file status is not successful`,
					Fields:      errorsInFiles,
				}
			}
		}
	}
	return nil
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
