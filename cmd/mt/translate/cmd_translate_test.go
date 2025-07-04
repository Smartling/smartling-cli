package translate

import (
	"bytes"
	"strings"
	"testing"

	mtmocks "github.com/Smartling/smartling-cli/cmd/mt/mocks"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/mt"
)

func TestNewTranslateCmdNoArgs(t *testing.T) {
	rlog.Init()
	buf := new(bytes.Buffer)
	params := srv.TranslateParams{
		InputDirectory: "/inputs/",
		TargetLocales:  []string{"en"},
	}

	initializer := mtmocks.NewMockSrvInitializer(t)

	cmd := NewTranslateCmd(initializer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{
		"--input-directory", params.InputDirectory,
		"--target-locale", params.TargetLocales[0],
	})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() expected with error")
	}

	output := buf.String()
	expected := "wrong argument quantity"

	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}

func TestNewTranslateCmdTargetLocale(t *testing.T) {
	rlog.Init()
	buf := new(bytes.Buffer)
	params := srv.TranslateParams{
		InputDirectory: "/inputs/",
	}

	initializer := mtmocks.NewMockSrvInitializer(t)

	cmd := NewTranslateCmd(initializer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{
		"--input-directory", params.InputDirectory,
	})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() expected with error")
	}

	output := buf.String()
	expected := "target-locale"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
	expected = "not set"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
	expected = "required flag(s)"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}
