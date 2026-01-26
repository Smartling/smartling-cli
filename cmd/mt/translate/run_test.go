package translate

import (
	"context"
	"errors"
	"testing"

	cmdmocks "github.com/Smartling/smartling-cli/cmd/mt/mocks"
	"github.com/Smartling/smartling-cli/output"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	srv "github.com/Smartling/smartling-cli/services/mt"
	srvmocks "github.com/Smartling/smartling-cli/services/mt/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRunGetFilesError(t *testing.T) {
	ctx := context.Background()
	initializer := cmdmocks.NewMockSrvInitializer(t)
	mtSrv := srvmocks.NewMockService(t)
	filesErr := errors.New("files error")

	params := srv.TranslateParams{
		InputDirectory: "/input/",
	}
	fileOrPattern := "*.txt"

	initializer.On("InitMTSrv").Return(mtSrv, nil)
	mtSrv.On("GetFiles", params.InputDirectory, fileOrPattern).Return(nil, filesErr)

	err := run(ctx, initializer, params, fileOrPattern, output.Params{})

	assert.Error(t, err)
	uiErr, ok := err.(clierror.UIError)
	assert.True(t, ok)
	assert.Equal(t, "get files", uiErr.Operation)
	assert.Equal(t, filesErr, uiErr.Err)
}

func TestRun(t *testing.T) {
	ctx := context.Background()
	initializer := cmdmocks.NewMockSrvInitializer(t)
	mtSrv := srvmocks.NewMockService(t)

	params := srv.TranslateParams{
		InputDirectory: "/input/",
	}
	fileOrPattern := "*.txt"
	files := []string{"file1.txt", "file2.txt"}

	initializer.On("InitMTSrv").Return(mtSrv, nil)
	mtSrv.On("GetFiles", params.InputDirectory, fileOrPattern).
		Return(files, nil)

	mtSrv.On("RunTranslate", mock.Anything, params, files, mock.Anything).
		Return([]srv.TranslateOutput{}, nil)

	err := run(ctx, initializer, params, fileOrPattern, output.Params{})

	assert.Nil(t, err)
}

func TestRunRunTranslateErr(t *testing.T) {
	ctx := context.Background()
	initializer := cmdmocks.NewMockSrvInitializer(t)
	mtSrv := srvmocks.NewMockService(t)

	params := srv.TranslateParams{
		InputDirectory: "/input/",
	}
	fileOrPattern := "*.txt"
	files := []string{"file1.txt", "file2.txt"}

	initializer.On("InitMTSrv").Return(mtSrv, nil)
	mtSrv.On("GetFiles", params.InputDirectory, fileOrPattern).
		Return(files, nil)

	err := clierror.UIError{
		Err:       errors.New("empty mtUid on start translation"),
		Operation: "Start translation",
		Fields: map[string]string{
			"file": "file1.txt",
			"code": "en",
		},
	}
	mtSrv.On("RunTranslate", mock.Anything, params, files, mock.Anything).
		Return([]srv.TranslateOutput{}, err)

	errRun := run(ctx, initializer, params, fileOrPattern, output.Params{})

	assert.NotNil(t, errRun)
	uiErrRun, ok := errRun.(clierror.UIError)
	assert.True(t, ok)

	assert.Equal(t, err.Operation, uiErrRun.Operation)
	assert.Equal(t, err.Description, uiErrRun.Description)
	assert.Equal(t, err.Err.Error(), uiErrRun.Err.Error())
}
