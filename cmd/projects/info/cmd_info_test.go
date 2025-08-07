package info

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"

	cmdmocks "github.com/Smartling/smartling-cli/cmd/projects/mocks"
	srvmocks "github.com/Smartling/smartling-cli/services/projects/mocks"

	"github.com/stretchr/testify/mock"
)

func TestNewInfoCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	projectsSrv := srvmocks.NewMockService(t)
	projectsSrv.On("RunInfo").Run(func(args mock.Arguments) {
		if _, err := fmt.Fprintf(buf, "RunInfo was called with %d args\n", len(args)); err != nil {
			log.Panic(err)
		}
	}).Return(nil)

	initializer := cmdmocks.NewMockSrvInitializer(t)
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
