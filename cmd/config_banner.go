package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	output "github.com/Smartling/smartling-cli/output/projects"
	projectconfig "github.com/Smartling/smartling-cli/services/projects/config"

	"golang.org/x/term"
)

// ShowConfigBanner prints the --show-config banner to stdout when the
// flag is set, and — only on an interactive terminal — prompts the operator
// with "Continue? [y/N]:" before letting the command proceed. On a "no"
// answer the process exits with status 1; on any other path (non-TTY, no
// flag, init command, config resolution failure) the function returns
// silently and the command runs normally.
func ShowConfigBanner(ctx context.Context) error {
	if !showConfig || isInit {
		return nil
	}
	cfg, err := Config()
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	client, err := Client(ctx)
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}
	extendedConfig, err := projectconfig.FetchExtendedConfig(ctx, cfg, client.GetProjectDetails)
	if err != nil {
		return err
	}
	if !showConfigAndMaybePrompt(extendedConfig, os.Stdout, os.Stderr, os.Stdin, stdinIsTerminal()) {
		return errors.New("operation aborted by user")
	}
	return nil
}

func showConfigAndMaybePrompt(config projectconfig.Extended, stdoutW, stderrW io.Writer, stdinR io.Reader, stdinIsTerminal bool) bool {
	_ = output.RenderPlain(stdoutW, config)

	if !stdinIsTerminal {
		return true
	}

	_, _ = fmt.Fprint(stderrW, "Continue? [y/N]: ")
	if !confirmContinue(stdinR) {
		_, _ = fmt.Fprintln(stderrW, "aborted")
		return false
	}
	return true
}

func confirmContinue(r io.Reader) bool {
	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		return false
	}
	answer := strings.TrimSpace(scanner.Text())
	if answer == "" {
		return false
	}
	switch strings.ToLower(answer)[0] {
	case 'y':
		return true
	default:
		return false
	}
}

// stdinIsTerminal reports whether stdin is connected to an interactive terminal.
func stdinIsTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}
