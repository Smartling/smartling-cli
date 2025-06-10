package client

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"regexp"

	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	sdk "github.com/Smartling/api-sdk-go"
	"github.com/kovetskiy/lorg"
	"github.com/reconquest/hierr-go"
)

var version = "1.7"

func CreateClient(clientConfig Config, config config.Config, verbose uint8) (*sdk.Client, error) {
	client := sdk.NewClient(config.UserID, config.Secret)

	var transport http.Transport

	if clientConfig.Insecure {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	if config.Proxy != "" && clientConfig.Proxy == "" {
		clientConfig.Proxy = config.Proxy
	}

	if clientConfig.Proxy != "" {
		proxy, err := url.Parse(clientConfig.Proxy)
		if err != nil {
			return nil, clierror.NewError(
				hierr.Errorf(
					err,
					"unable to parse specified proxy URL",
				),

				`Proxy should be valid URL, check CLI options and `+
					`config value.`,
			)
		}

		transport.Proxy = http.ProxyURL(proxy)
	}

	if clientConfig.SmartlingURL != "" {
		client.BaseURL = clientConfig.SmartlingURL
	}

	client.HTTP.Transport = &transport
	client.UserAgent = "smartling-cli/" + version

	setLogger(client, rlog.Logger(), verbose)

	rlog.HideRegexp(
		regexp.MustCompile(`"(?:access|refresh)Token": "([^"]+)"`),
	)

	err := client.Authenticate()
	if err != nil {
		return nil, clierror.NewError(
			err,
			`Your credentials are invalid. Double check it and try to run init.\n`,
		)
	}

	return client, nil
}

func setLogger(client *sdk.Client, logger lorg.Logger, verbosity uint8) {
	switch verbosity {
	case 0:
		return

	case 1:
		client.SetInfoLogger(logger.Infof)

	default:
		client.SetDebugLogger(logger.Debugf)
	}
}
