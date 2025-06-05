package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/client"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	"os"
	"regexp"

	smartling "github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
	"github.com/tcnksm/go-input"
)

func doInit(config config.Config, args map[string]interface{}, cliClientConfig client.Config) error {
	fmt.Printf("Generating %s...\n\n", config.path)

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
			fmt.Fprintf(os.Stderr, "failed to scan input: "+err.Error())
		}
	}

	var input config.Config

	prompt(
		"Smartling API V2.0 User Identifier",
		config.UserID,
		config.UserID == "",
		false,
		&input.UserID,
	)

	if input.UserID != "" {
		config.UserID = input.UserID
	}

	prompt(
		"Smartling API V2.0 Token Secret",
		config.Secret,
		config.Secret == "",
		true,
		&input.Secret,
	)

	if input.Secret != "" {
		config.Secret = input.Secret
	}

	prompt(
		"Account ID (optional)",
		config.AccountID,
		config.AccountID == "",
		false,
		&input.AccountID,
	)

	if input.AccountID != "" {
		config.AccountID = input.AccountID
	}

	prompt(
		"Project ID",
		config.ProjectID,
		config.ProjectID == "",
		false,
		&input.ProjectID,
	)

	if input.ProjectID != "" {
		config.ProjectID = input.ProjectID
	}

	var result bytes.Buffer
	err := config.configTemplate.Execute(&result, config)
	if err != nil {
		return hierr.Errorf(
			err,
			"unable to compile config template",
		)
	}

	logger.HideFromConfig(config)

	fmt.Println("Testing connection to Smartling API...")

	client, err := createClient(config, cliClientConfig)
	if err != nil {
		return hierr.Errorf(
			err,
			"unable to create client for testing connection",
		)
	}

	err = client.Authenticate()
	if err != nil {
		if _, ok := err.(smartling.NotAuthorizedError); ok {
			return clierror.NewError(
				errors.New("not authorized"),
				"Your credentials are invalid. Double check them and run "+
					"init command again",
			)
		} else {
			return clierror.NewError(
				hierr.Errorf(err, "failure while testing connection"),
				"Contact developer for more info",
			)
		}
	}

	fmt.Println("Connection is successful.")

	if args["--dry-run"].(bool) {
		fmt.Println()
		fmt.Println("Configured for Dry run. Not writing config file.")
		fmt.Println("New config is displayed below.")
		fmt.Println()

		fmt.Println(result.String())
	} else {
		err = os.WriteFile(config.path, result.Bytes(), 0644)
		if err != nil {
			return hierr.Errorf(
				err,
				"unable to write new config file",
			)
		}
	}

	return nil
}
