package list

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	projectscmdmocks "github.com/Smartling/smartling-cli/cmd/projects/mocks"
	projectsmocks "github.com/Smartling/smartling-cli/services/projects/mocks"

	"github.com/stretchr/testify/mock"
)

func TestNewListCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	projectsSrv := projectsmocks.NewMockService(t)
	projectsSrv.On("RunList", true).Run(func(args mock.Arguments) {
		fmt.Fprintln(buf, fmt.Sprintf("RunList was called with %d args", len(args)))
		fmt.Fprintln(buf, fmt.Sprintf("short: %v", args[0]))
	}).Return(nil)

	initializer := projectscmdmocks.NewMockSrvInitializer(t)
	initializer.On("InitProjectsSrv").Return(projectsSrv, nil)

	cmd := NewListCmd(initializer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--short", "true"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() returned an error: %v", err)
	}

	output := buf.String()
	expected := `RunList was called with 1 args
short: true
`

	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}
