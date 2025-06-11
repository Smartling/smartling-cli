package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Smartling/smartling-cli/services/files"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	globfiles "github.com/Smartling/smartling-cli/services/helpers/glob_files"

	sdk "github.com/Smartling/api-sdk-go"
	"github.com/Smartling/smartling-cli/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type request struct {
	response string
	code     int
}

type roundTripFunc func(req *http.Request) *http.Response

func (function roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return function(req), nil
}

func mockHTTPClient(function roundTripFunc) *http.Client {
	return &http.Client{
		Transport: function,
	}
}

func TestPushStopUnauthorized(t *testing.T) {
	params := getPushParams("README.md README.md")

	httpClient := getMockHTTPClient([]request{{"{}", 401}})

	mockGlobber(params)
	defer func() {
		globfiles.GlobFilesLocally = globfiles.LocallyFn
	}()

	client := getClient(httpClient)

	filesSrv := files.NewService(&client, getConfig(), config.FileConfig{})

	err := filesSrv.RunPush(params)

	assert.True(t, errors.Is(err, sdk.NotAuthorizedError{}))
}

func TestPushContinueFakeError(t *testing.T) {
	params := getPushParams("README.md README.md")

	mockGlobber(params)
	defer func() {
		globfiles.GlobFilesLocally = globfiles.LocallyFn
	}()

	client := &mocks.ClientInterface{}
	client.On("UploadFile", "test", mock.Anything).
		Return(nil, sdk.APIError{Cause: errors.New("some error")}).
		Times(2)

	filesSrv := files.NewService(client, getConfig(), config.FileConfig{})
	err := filesSrv.RunPush(params)
	assert.EqualError(
		t,
		err,
		"ERROR: failed to upload 2 files\n\nfailed to upload files README.md, README.md")
	client.AssertExpectations(t)
}

func TestPushStopApiError(t *testing.T) {
	params := getPushParams("README.md README.md")

	mockGlobber(params)
	defer func() {
		globfiles.GlobFilesLocally = globfiles.LocallyFn
	}()

	client := &mocks.ClientInterface{}
	expectedError := sdk.APIError{
		Cause: errors.New("some error"),
		Code:  "MAINTENANCE_MODE_ERROR",
	}
	client.On("UploadFile", "test", mock.Anything).
		Return(nil, expectedError).
		Once()

	filesSrv := files.NewService(client, getConfig(), config.FileConfig{})
	err := filesSrv.RunPush(params)

	assert.True(t, errors.Is(err, expectedError))
	client.AssertExpectations(t)
}

func getMockHTTPClient(responses []request) *http.Client {
	responseCount := 0
	return mockHTTPClient(func(_ *http.Request) *http.Response {
		var response string
		var statusCode int
		header := make(http.Header)
		header.Add("Content-Type", "application/json")
		if responseCount >= len(responses) {
			response = responses[len(responses)-1].response
			statusCode = responses[len(responses)-1].code
		} else {
			response = responses[responseCount].response
			statusCode = responses[responseCount].code
			responseCount++
		}
		return &http.Response{
			StatusCode: statusCode,
			Body:       io.NopCloser(bytes.NewBufferString(response)),
			Header:     header,
		}
	})
}

func getPushParams(file string) files.PushParams {
	return files.PushParams{
		Authorize:  false,
		Directory:  "",
		File:       file,
		Directives: nil,
	}
}

func getConfig() config.Config {
	fileConfig := make(map[string]config.FileConfig)
	fileConfig["default"] = config.FileConfig{
		Push: struct {
			Type       string            `yaml:"type,omitempty"`
			Directives map[string]string `yaml:"directives,omitempty,flow"`
		}{Type: "md"},
	}
	return config.Config{
		UserID:    "test",
		Secret:    "test",
		ProjectID: "test",
		Files:     fileConfig,
	}
}

func getClient(httpClient *http.Client) sdk.Client {
	client := sdk.NewClient("test", "test")
	client.HTTP = httpClient

	return *client
}

func mockGlobber(params files.PushParams) {
	globfiles.GlobFilesLocally = func(_, _, _ string) ([]string, error) {
		return strings.Split(params.File, " "), nil
	}
}
