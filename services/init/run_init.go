package initialize

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	sdk "github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
	"github.com/tcnksm/go-input"
)

// RunInit initializes the Smartling CLI.
func (s service) RunInit(dryRun bool) error {
	fmt.Printf("Generating %s...\n\n", s.Config.Path)

	prompt := func(
		message string,
		value interface{},
		zero bool,
		hidden bool,
		variable interface{},
	) {
		display := regexp.MustCompile(`^(.{1,3}).*$`).ReplaceAllString(
			fmt.Sprint(value),
			`$1***`,
		)

		if !zero {
			message = fmt.Sprintf("%s [default is %q]", message, display)
		}

		read, err := input.DefaultUI().Ask(
			message,
			&input.Options{
				Default:     fmt.Sprint(value),
				Hide:        hidden,
				HideDefault: true,
			},
		)
		if err != nil {
			if input.ErrInterrupted == err {
				os.Exit(1)
			}
		}

		_, err = fmt.Sscanln(read, variable)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to scan input: %s", err.Error())
		}
	}

	var input config.Config

	prompt(
		"Smartling API V2.0 User Identifier",
		s.Config.UserID,
		s.Config.UserID == "",
		false,
		&input.UserID,
	)

	if input.UserID != "" {
		s.Config.UserID = input.UserID
	}

	prompt(
		"Smartling API V2.0 Token Secret",
		s.Config.Secret,
		s.Config.Secret == "",
		true,
		&input.Secret,
	)

	if input.Secret != "" {
		s.Config.Secret = input.Secret
	}

	prompt(
		"Account ID (optional)",
		s.Config.AccountID,
		s.Config.AccountID == "",
		false,
		&input.AccountID,
	)

	if input.AccountID != "" {
		s.Config.AccountID = input.AccountID
	}

	prompt(
		"Project ID",
		s.Config.ProjectID,
		s.Config.ProjectID == "",
		false,
		&input.ProjectID,
	)

	if input.ProjectID != "" {
		s.Config.ProjectID = input.ProjectID
	}

	var result bytes.Buffer
	err := configTemplate.Execute(&result, s.Config)
	if err != nil {
		return hierr.Errorf(
			err,
			"unable to compile config template",
		)
	}

	rlog.HideString(s.Config.Secret)
	rlog.HideString(s.Config.UserID)
	rlog.HideString(s.Config.AccountID)
	rlog.HideString(s.Config.ProjectID)

	fmt.Println("Testing connection to Smartling API...")

	err = s.Client.Authenticate()
	if err != nil {
		if _, ok := err.(sdk.NotAuthorizedError); ok {
			return clierror.NewError(
				errors.New("not authorized"),
				"Your credentials are invalid. Double check them and run "+
					"init command again",
			)
		}
		return clierror.NewError(
			hierr.Errorf(err, "failure while testing connection"),
			"Contact developer for more info",
		)
	}

	fmt.Println("Connection is successful.")

	if dryRun {
		fmt.Println()
		fmt.Println("Configured for Dry run. Not writing config file.")
		fmt.Println("New config is displayed below.")
		fmt.Println()

		fmt.Println(result.String())
	} else {
		err = os.WriteFile(s.Config.Path, result.Bytes(), 0644)
		if err != nil {
			return hierr.Errorf(
				err,
				"unable to write new config file",
			)
		}
	}

	return nil
}
