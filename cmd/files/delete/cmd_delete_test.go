package delete

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	cmdmocks "github.com/Smartling/smartling-cli/cmd/files/mocks"
	srvmocks "github.com/Smartling/smartling-cli/services/files/mocks"

	"github.com/stretchr/testify/mock"
)

func TestNewDeleteCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	filesSrv := srvmocks.NewMockService(t)
	uriArg := "https://example.com:8080/path/to/resource?search=a"
	filesSrv.On("RunDelete", mock.Anything, uriArg).Run(func(args mock.Arguments) {
		fmt.Fprintf(buf, "RunDelete was called with %d args\n", len(args))
		fmt.Fprintf(buf, "uri: %v\n", args[1])
	}).Return(nil)

	initializer := cmdmocks.NewMockSrvInitializer(t)
	initializer.On("InitFilesSrv", mock.Anything).Return(filesSrv, nil)

	cmd := NewDeleteCmd(initializer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{uriArg})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() returned an error: %v", err)
	}

	output := buf.String()
	expected := fmt.Sprintf(`RunDelete was called with 2 args
uri: %s
`, uriArg)

	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}
