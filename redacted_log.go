package main

import (
	"fmt"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	"io"
	"os"
	"regexp"

	"github.com/kovetskiy/lorg"
)

type redactedLog struct {
	*lorg.Log

	writer *redactedWriter
}

func NewRedactedLog() *redactedLog {
	log := &redactedLog{
		Log:    lorg.NewLog(),
		writer: &redactedWriter{},
	}

	log.SetOutput(log.writer)

	return log
}

func (log *redactedLog) ToggleRedact(enable bool) {
	log.writer.enabled = enable
}

func (log *redactedLog) HideRegexp(pattern *regexp.Regexp) {
	log.writer.patterns = append(log.writer.patterns, pattern)
}

func (log *redactedLog) HideString(value string) {
	pattern := regexp.MustCompile(
		fmt.Sprintf(
			"(%s)",
			regexp.QuoteMeta(value),
		),
	)

	log.writer.patterns = append(log.writer.patterns, pattern)
}

func (log *redactedLog) HideFromConfig(config config.Config) {
	log.HideString(config.Secret)
	log.HideString(config.UserID)
	log.HideString(config.AccountID)
	log.HideString(config.ProjectID)
}

func (log *redactedLog) GetWriter() io.Writer {
	return log.writer
}

type redactedWriter struct {
	patterns []*regexp.Regexp
	enabled  bool
}

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
