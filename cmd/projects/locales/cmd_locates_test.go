package locales

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	cmdmocks "github.com/Smartling/smartling-cli/cmd/projects/mocks"
	"github.com/Smartling/smartling-cli/services/projects"
	srvmocks "github.com/Smartling/smartling-cli/services/projects/mocks"

	"github.com/stretchr/testify/mock"
)

func TestNewLocatesCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	projectsSrv := srvmocks.NewMockService(t)
	params := projects.LocalesParams{
		Format: "any",
		Short:  true,
		Source: false,
	}
	projectsSrv.On("RunLocales", params).Run(func(args mock.Arguments) {
		fmt.Fprintln(buf, fmt.Sprintf("RunLocales was called with %d args", len(args)))
		fmt.Fprintln(buf, fmt.Sprintf("params: %v", args[0]))
	}).Return(nil)

	initializer := cmdmocks.NewMockSrvInitializer(t)
	initializer.On("InitProjectsSrv").Return(projectsSrv, nil)

	cmd := NewLocatesCmd(initializer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{
		"--short", fmt.Sprintf("%v", params.Short),
		"--format", params.Format,
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() returned an error: %v", err)
	}

	output := buf.String()
	expected := fmt.Sprintf(`RunLocales was called with 1 args
params: %v
`, params)

	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}
