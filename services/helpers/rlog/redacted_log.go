package rlog

import (
	"fmt"
	"os"
	"regexp"

	"github.com/kovetskiy/lorg"
)

var logger *RedactedLog

// RedactedLog is a custom logger with a writer.
type RedactedLog struct {
	*lorg.Log

	writer *redactedWriter
}

// Init initializes RedactedLog.
func Init() {
	logger = &RedactedLog{
		Log:    lorg.NewLog(),
		writer: &redactedWriter{},
	}
	logger.SetOutput(logger.writer)
}

func (log *RedactedLog) hideString(value string) {
	pattern := regexp.MustCompile(
		fmt.Sprintf(
			"(%s)",
			regexp.QuoteMeta(value),
		),
	)

	logger.writer.patterns = append(logger.writer.patterns, pattern)
}

// ToggleRedact enables or disables the redaction of sensitive information in logs.
func ToggleRedact(enable bool) {
	logger.writer.enabled = enable
}

// HideRegexp adds a regular expression pattern to the logger's redaction list.
func HideRegexp(pattern *regexp.Regexp) {
	logger.writer.patterns = append(logger.writer.patterns, pattern)
}

// HideString adds a string to the logger's redaction list, which will be replaced in the logs.
func HideString(value string) {
	logger.hideString(value)
}

// SetFormat sets formatting for the logger.
func SetFormat(format lorg.Formatter) {
	logger.SetFormat(format)
}

// SetIndentLines sets indent.
func SetIndentLines(value bool) {
	logger.SetIndentLines(value)
}

// SetLevel sets the logging level for the logger.
func SetLevel(level lorg.Level) {
	logger.SetLevel(level)
}

// Logger returns the RedactedLog instance.
func Logger() *RedactedLog {
	return logger
}

type redactedWriter struct {
	patterns []*regexp.Regexp
	enabled  bool
}

// Write write without sensitive information
func (writer redactedWriter) Write(buffer []byte) (int, error) {
	if !writer.enabled {
		return os.Stderr.Write(buffer)
	}

	output := string(buffer)

	placeholder := "***"

	for _, pattern := range writer.patterns {
		output = pattern.ReplaceAllStringFunc(
			output,
			func(value string) string {
				i := pattern.FindStringSubmatchIndex(value)
				if len(i) < 4 {
					return value
				}

				if len(value) < i[2]+3 {
					return value
				}

				// NOTE: Cut out first 3 characters of first regexp submatch,
				// NOTE: which identifies secret.
				return value[:i[2]+3] + placeholder + value[i[3]:]
			},
		)
	}

	return os.Stderr.Write([]byte(output))
}
