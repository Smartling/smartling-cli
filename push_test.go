package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Smartling/smartling-cli/mocks"

	smartling "github.com/Smartling/api-sdk-go"
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

func mockHttpClient(function roundTripFunc) *http.Client {
	return &http.Client{
		Transport: function,
	}
}

func TestPushStopUnauthorized(t *testing.T) {
	args := getArgs("README.md README.md")

	httpClient := getMockHttpClient([]request{{"{}", 401}})

	mockGlobber(args)
	defer func() {
		globFilesLocally = globFilesLocallyFunc
	}()

	client := getClient(httpClient)

	err := doFilesPush(&client, getConfig(), args)

	assert.True(t, errors.Is(err, smartling.NotAuthorizedError{}))
}

func TestPushContinueFakeError(t *testing.T) {
	args := getArgs("README.md README.md")

	mockGlobber(args)
	defer func() {
		globFilesLocally = globFilesLocallyFunc
	}()

	client := &mocks.ClientInterface{}
	client.On("UploadFile", "test", mock.Anything).
		Return(nil, smartling.APIError{Cause: errors.New("some error")}).
		Times(2)

	err := doFilesPush(client, getConfig(), args)
	assert.EqualError(
		t,
		err,
		"ERROR: failed to upload 2 files\n\nfailed to upload files README.md, README.md")
	client.AssertExpectations(t)
}

func TestPushStopApiError(t *testing.T) {
	args := getArgs("README.md README.md")

	mockGlobber(args)
	defer func() {
		globFilesLocally = globFilesLocallyFunc
	}()

	client := &mocks.ClientInterface{}
	expectedError := smartling.APIError{
		Cause: errors.New("some error"),
		Code:  "MAINTENANCE_MODE_ERROR",
	}
	client.On("UploadFile", "test", mock.Anything).
		Return(nil, expectedError).
		Once()

	err := doFilesPush(client, getConfig(), args)

	assert.True(t, errors.Is(err, expectedError))
	client.AssertExpectations(t)
}

func getMockHttpClient(responses []request) *http.Client {
	responseCount := 0
	return mockHttpClient(func(req *http.Request) *http.Response {
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

func getArgs(file string) map[string]interface{} {
	args := make(map[string]interface{})
	args["--authorize"] = false
	args["--directory"] = ""
	args["<file>"] = file
	return args
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

func getClient(httpClient *http.Client) smartling.Client {
	client := smartling.NewClient("test", "test")
	client.HTTP = httpClient

	return *client
}

func mockGlobber(args map[string]interface{}) {
	globFilesLocally = func(
		directory string,
		base string,
		mask string,
	) ([]string, error) {
		return strings.Split(fmt.Sprintf("%s", args["<file>"]), " "), nil
	}
}
