package client

import (
	redactedlog "github.com/Smartling/smartling-cli/services/helpers/redacted_log"
)

// TODO replace legacy
var redactedLogger redactedlog.RedactedLog

func InitLogger(l redactedlog.RedactedLog) {
	redactedLogger = l
}
