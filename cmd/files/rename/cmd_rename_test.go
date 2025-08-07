package rename_test

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"

	cmdmocks "github.com/Smartling/smartling-cli/cmd/files/mocks"
	"github.com/Smartling/smartling-cli/cmd/files/rename"
	srvmocks "github.com/Smartling/smartling-cli/services/files/mocks"

	"github.com/stretchr/testify/mock"
)

func TestNewRenameCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	filesSrv := srvmocks.NewMockService(t)
	oldArg := "_old"
	newArg := "_new"
	filesSrv.On("RunRename", oldArg, newArg).Run(func(args mock.Arguments) {
		if _, err := fmt.Fprintf(buf, "RunRename was called with %d args\n", len(args)); err != nil {
			log.Panic(err)
		}
		if _, err := fmt.Fprintf(buf, "old: %v\n", args[0]); err != nil {
			log.Panic(err)
		}
		if _, err := fmt.Fprintf(buf, "new: %v\n", args[1]); err != nil {
			log.Panic(err)
		}
	}).Return(nil)

	initializer := cmdmocks.NewMockSrvInitializer(t)
	initializer.On("InitFilesSrv").Return(filesSrv, nil)

	cmd := rename.NewRenameCmd(initializer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{
		oldArg,
		newArg,
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() returned an error: %v", err)
	}

	output := buf.String()
	expected := fmt.Sprintf(`RunRename was called with 2 args
old: %s
new: %v
`,
		oldArg,
		newArg,
	)

	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}
