package info

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	projectscmdmocks "github.com/Smartling/smartling-cli/cmd/projects/mocks"
	projectsmocks "github.com/Smartling/smartling-cli/services/projects/mocks"

	"github.com/stretchr/testify/mock"
)

func TestNewInfoCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	projectsSrv := projectsmocks.NewMockService(t)
	projectsSrv.On("RunInfo").Run(func(args mock.Arguments) {
		fmt.Fprintln(buf, fmt.Sprintf("RunInfo was called with %d args", len(args)))
	}).Return(nil)

	initializer := projectscmdmocks.NewMockSrvInitializer(t)
	initializer.On("InitProjectsSrv").Return(projectsSrv, nil)

	cmd := NewInfoCmd(initializer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() returned an error: %v", err)
	}

	output := buf.String()
	expected := "RunInfo was called with 0 args"

	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}
