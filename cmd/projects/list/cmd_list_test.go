package list

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	cmdmocks "github.com/Smartling/smartling-cli/cmd/projects/mocks"
	srvmocks "github.com/Smartling/smartling-cli/services/projects/mocks"

	"github.com/stretchr/testify/mock"
)

func TestNewListCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	projectsSrv := srvmocks.NewMockService(t)
	shortArg := true
	projectsSrv.On("RunList", shortArg).Run(func(args mock.Arguments) {
		fmt.Fprintln(buf, fmt.Sprintf("RunList was called with %d args", len(args)))
		fmt.Fprintln(buf, fmt.Sprintf("short: %v", args[0]))
	}).Return(nil)

	initializer := cmdmocks.NewMockSrvInitializer(t)
	initializer.On("InitProjectsSrv").Return(projectsSrv, nil)

	cmd := NewListCmd(initializer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--short", fmt.Sprintf("%v", shortArg)})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() returned an error: %v", err)
	}

	output := buf.String()
	expected := fmt.Sprintf(`RunList was called with 1 args
short: %v
`, shortArg)

	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}
