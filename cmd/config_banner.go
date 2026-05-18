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
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	projectconfig "github.com/Smartling/smartling-cli/services/projects/config"

	"golang.org/x/term"
)

// ShowConfigBanner prints the --show-config banner to stdout when the
// flag is set, and — only on an interactive terminal — prompts the operator
// with "Continue? [y/N]:" before letting the command proceed.
func ShowConfigBanner(ctx context.Context) error {
	if !showConfig || isInit {
		return nil
	}
	config, err := Config()
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	var extendedConfig projectconfig.Extended
	extendedConfig.InjectConfig(config)
	client, err := Client(ctx)
	if err != nil {
		rlog.Errorf("failed to create API client: %s", err)
	} else if config.ProjectID != "" {
		projectDetails, err := client.GetProjectDetails(ctx, config.ProjectID)
		if err != nil {
			rlog.Errorf("failed to fetch project details: %s", err)
		}
		if projectDetails != nil {
			extendedConfig.InjectProject(*projectDetails)
		}
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
	switch strings.ToLower(strings.TrimSpace(scanner.Text())) {
	case "y", "yes":
		return true
	default:
		return false
	}
}

// stdinIsTerminal reports whether stdin is connected to an interactive terminal.
func stdinIsTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}
