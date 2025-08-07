package initialize

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"

	cmdmocks "github.com/Smartling/smartling-cli/cmd/init/mocks"
	srvmocks "github.com/Smartling/smartling-cli/services/init/mocks"

	"github.com/stretchr/testify/mock"
)

func TestNewInitCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	initSrv := srvmocks.NewMockService(t)
	initSrv.On("RunInit", false).Run(func(args mock.Arguments) {
		if _, err := fmt.Fprintf(buf, "RunInit was called with: %v\n", args[0]); err != nil {
			log.Panic(err)
		}
	}).Return(nil)

	initializer := cmdmocks.NewMockSrvInitializer(t)
	initializer.On("InitSrv").Return(initSrv, nil)

	cmd := NewInitCmd(initializer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() returned an error: %v", err)
	}

	output := buf.String()
	expected := "RunInit was called with: false"

	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}

func TestNewInitCmdDryRun(t *testing.T) {
	buf := new(bytes.Buffer)
	initSrv := srvmocks.NewMockService(t)
	initSrv.On("RunInit", true).Run(func(args mock.Arguments) {
		if _, err := fmt.Fprintf(buf, "RunInit was called with: %v\n", args[0]); err != nil {
			log.Panic(err)
		}
	}).Return(nil)

	initializer := cmdmocks.NewMockSrvInitializer(t)
	initializer.On("InitSrv").Return(initSrv, nil)

	cmd := NewInitCmd(initializer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{})
	cmd.SetArgs([]string{"--dry-run", "true"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() returned an error: %v", err)
	}

	output := buf.String()
	expected := "RunInit was called with: true"

	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}
