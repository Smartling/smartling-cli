package list

import (
	"bytes"
	"fmt"
	"log"
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
		if _, err := fmt.Fprintf(buf, "RunList was called with %d args\n", len(args)); err != nil {
			log.Panic(err)
		}
		if _, err := fmt.Fprintf(buf, "format: %v\n", args[0]); err != nil {
			log.Panic(err)
		}
		if _, err := fmt.Fprintf(buf, "short: %v\n", args[1]); err != nil {
			log.Panic(err)
		}
		if _, err := fmt.Fprintf(buf, "uri: %v\n", args[2]); err != nil {
			log.Panic(err)
		}
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
