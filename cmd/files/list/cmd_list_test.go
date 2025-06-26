package list

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	cmdmocks "github.com/Smartling/smartling-cli/cmd/files/mocks"
	srvmocks "github.com/Smartling/smartling-cli/services/files/mocks"

	"github.com/stretchr/testify/mock"
)

func TestNewListCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	filesSrv := srvmocks.NewMockService(t)
	formatTypeArg := "any"
	shortArg := true
	uriArg := "https://example.com:8080/path/to/resource?search=a"
	filesSrv.On("RunList", formatTypeArg, shortArg, uriArg).Run(func(args mock.Arguments) {
		fmt.Fprintln(buf, fmt.Sprintf("RunList was called with %d args", len(args)))
		fmt.Fprintln(buf, fmt.Sprintf("format: %v", args[0]))
		fmt.Fprintln(buf, fmt.Sprintf("short: %v", args[1]))
		fmt.Fprintln(buf, fmt.Sprintf("uri: %v", args[2]))
	}).Return(nil)

	initializer := cmdmocks.NewMockSrvInitializer(t)
	initializer.On("InitFilesSrv").Return(filesSrv, nil)

	cmd := NewListCmd(initializer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{
		uriArg,
		"--short",
		"--format", formatTypeArg,
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() returned an error: %v", err)
	}

	output := buf.String()
	expected := fmt.Sprintf(`RunList was called with 3 args
format: %s
short: %v
uri: %s
`,
		formatTypeArg,
		shortArg,
		uriArg,
	)

	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}
