package initialize

import (
	redactedlog "github.com/Smartling/smartling-cli/services/helpers/redacted_log"
)

// TODO replace legacy
var logger redactedlog.RedactedLog

func InitLogger(l *redactedlog.RedactedLog) {
	logger = *l
}
