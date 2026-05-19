package pull

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	cmdmocks "github.com/Smartling/smartling-cli/cmd/files/mocks"
	"github.com/Smartling/smartling-cli/services/files"
	srvmocks "github.com/Smartling/smartling-cli/services/files/mocks"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	// Production main() calls rlog.Init(); the Run handler in NewPullCmd
	// reaches rootcmd.Config() which calls rlog.Debugf — without Init the
	// global logger is nil and the test panics.
	rlog.Init()
	m.Run()
}

func TestNewPullCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	filesSrv := srvmocks.NewMockService(t)
	params := files.PullParams{
		URI:       "https://example.com:8080/path/to/resource?search=a",
		Format:    "txt",
		Directory: "/texts/",
		Source:    true,
		Locales:   []string{"en-US", "fr-FR"},
		Progress:  "20%",
		Retrieve:  "none",
		Threads:   20,
	}
	filesSrv.On("RunPull", mock.Anything, params).Run(func(args mock.Arguments) {
		if _, err := fmt.Fprintf(buf, "RunPull was called with %d args\n", len(args)); err != nil {
			t.Fatal(err)
		}
		if _, err := fmt.Fprintf(buf, "params: %v\n", args[1]); err != nil {
			t.Fatal(err)
		}
	}).Return(nil)

	initializer := cmdmocks.NewMockSrvInitializer(t)
	initializer.On("InitFilesSrv", mock.Anything).Return(filesSrv, nil)

	cmd := NewPullCmd(initializer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{
		params.URI,
		"--format", params.Format,
		"--directory", params.Directory,
		"--source", fmt.Sprintf("%v", params.Source),
		"--locale", params.Locales[0],
		"--locale", params.Locales[1],
		"--progress", params.Progress,
		"--retrieve", params.Retrieve,
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() returned an error: %v", err)
	}

	output := buf.String()
	expected := fmt.Sprintf(`RunPull was called with 2 args
params: %v
`, params)

	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}
