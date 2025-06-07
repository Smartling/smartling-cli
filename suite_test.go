package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"sort"
	"strings"
	"testing"

	sdk "github.com/Smartling/api-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	codeSuccess = "SUCCESS"
)

type MainSuite struct {
	suite.Suite

	Mock struct {
		Server  *httptest.Server
		Handler http.HandlerFunc
	}
}

func (suite *MainSuite) SetupSuite() {
	suite.Mock.Server = httptest.NewUnstartedServer(
		http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				auth, err := handleAuthentication(writer, request)
				if err != nil {
					panic(err)
				}

				if auth {
					return
				}

				if suite.Mock.Handler != nil {
					suite.Mock.Handler(writer, request)
				}
			},
		),
	)

	suite.Mock.Server.Config.SetKeepAlivesEnabled(false)
	suite.Mock.Server.StartTLS()
}

func (suite *MainSuite) TearDownSuite() {
	suite.Mock.Server.Close()
}

func (suite *MainSuite) run(opts ...interface{}) (bool, string, string) {
	var (
		success = true
		stdout  = &bytes.Buffer{}
		stderr  = &bytes.Buffer{}
	)

	args := []string{
		"-test.run=Test_Run",
		"--",
		"--insecure",
		"--smartling-url",
		suite.Mock.Server.URL,
	}

	for _, opt := range opts {
		switch opt := opt.(type) {
		case string:
			args = append(args, opt)
		}
	}

	cmd := exec.Command(
		os.Args[0],
		args...,
	)

	cmd.Env = append(cmd.Env, "_TEST_RUN=1")
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err, ok := err.(*exec.ExitError); ok {
		success = err.Success()
	}

	stdoutString := strings.TrimSuffix(stdout.String(), "PASS\n")
	stderrString := stderr.String()

	return success, stdoutString, stderrString
}

func (suite *MainSuite) assertStdout(output []string, args ...interface{}) {
	success, stdout, _ := suite.run(args...)

	assert.True(suite.T(), success)

	sorted := strings.Split(stdout, "\n")
	sort.Strings(sorted)
	sort.Strings(output)

	assert.Equal(
		suite.T(),
		strings.Join(output, "\n"),
		strings.TrimSpace(strings.Join(sorted, "\n")),
	)
}

func TestMainSuite(t *testing.T) {
	suite.Run(t, &MainSuite{})
}

func Test_Run(t *testing.T) {
	if os.Getenv("_TEST_RUN") != "1" {
		t.SkipNow()
	}

	var (
		skip = true
		args []string
	)

	for _, arg := range os.Args {
		if !skip {
			args = append(args, arg)
		}

		if arg == "--" {
			skip = false
		}
	}

	os.Args = []string{os.Args[0]}
	os.Args = append(os.Args, args...)

	main()
}

func writeSmartlingReply(
	writer http.ResponseWriter,
	code string,
	reply interface{},
) error {
	var payload struct {
		Response struct {
			Code string
			Data interface{}
		}
	}

	payload.Response.Code = code
	payload.Response.Data = reply

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	return json.NewEncoder(writer).Encode(payload)
}

func handleAuthentication(
	writer http.ResponseWriter,
	request *http.Request,
) (bool, error) {
	if !strings.HasSuffix(request.URL.Path, "/authenticate") {
		return false, nil
	}

	err := writeSmartlingReply(
		writer,
		codeSuccess,
		struct {
			AccessToken      string
			ExpiresIn        int
			RefreshExpiresIn int
		}{
			"somenewtoken",
			480,
			3660,
		},
	)

	if err != nil {
		return false, err
	}

	return true, err
}

func utc(
	timestamp string,
) sdk.UTC {
	var utc sdk.UTC

	data, _ := json.Marshal(timestamp)

	err := utc.UnmarshalJSON([]byte(data))
	if err != nil {
		panic(err)
	}

	return utc
}
