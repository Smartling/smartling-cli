package main

import (
	"github.com/Smartling/api-sdk-go"
	"github.com/kovetskiy/lorg"
)

func setLogger(client *smartling.Client, logger lorg.Logger, verbosity uint8) {
	switch verbosity {
	case 0:
		return

	case 1:
		client.SetInfoLogger(logger.Infof)

	default:
		client.SetDebugLogger(logger.Debugf)
	}
}
