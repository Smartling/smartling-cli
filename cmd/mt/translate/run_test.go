package translate

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"

	cmdmocks "github.com/Smartling/smartling-cli/cmd/mt/mocks"
	output "github.com/Smartling/smartling-cli/output/mt"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	srv "github.com/Smartling/smartling-cli/services/mt"
	srvmocks "github.com/Smartling/smartling-cli/services/mt/mocks"

	"github.com/stretchr/testify/assert"
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

	err := run(ctx, initializer, params, fileOrPattern, output.OutputParams{})

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

	err := run(ctx, initializer, params, fileOrPattern, output.OutputParams{})

	assert.Nil(t, err)
}
