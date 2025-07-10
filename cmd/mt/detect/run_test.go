package detect

import (
	"context"
	"errors"
	"testing"

	cmdmocks "github.com/Smartling/smartling-cli/cmd/mt/mocks"
	output "github.com/Smartling/smartling-cli/output/mt"
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

	params := srv.DetectParams{
		InputDirectory: "/input/",
		FileOrPattern:  "*.txt",
	}

	initializer.On("InitMTSrv").Return(mtSrv, nil)
	mtSrv.On("GetFiles", params.InputDirectory, params.FileOrPattern).Return(nil, filesErr)

	err := run(ctx, initializer, params, output.OutputParams{})

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

	params := srv.DetectParams{
		InputDirectory: "/input/",
		FileOrPattern:  "*.txt",
	}
	files := []string{"file1.txt", "file2.txt"}

	initializer.On("InitMTSrv").Return(mtSrv, nil)
	mtSrv.On("GetFiles", params.InputDirectory, params.FileOrPattern).
		Return(files, nil)

	mtSrv.On("RunDetect", mock.Anything, params, files, mock.Anything).
		Return([]srv.DetectOutput{}, nil)

	err := run(ctx, initializer, params, output.OutputParams{})

	assert.Nil(t, err)
}
