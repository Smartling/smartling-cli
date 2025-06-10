package rlog

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/kovetskiy/lorg"
)

var logger *RedactedLog

type RedactedLog struct {
	*lorg.Log

	writer *redactedWriter
}

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

func ToggleRedact(enable bool) {
	logger.writer.enabled = enable
}

func HideRegexp(pattern *regexp.Regexp) {
	logger.writer.patterns = append(logger.writer.patterns, pattern)
}

func HideFromConfig(value string) {
	logger.hideString(value)
}

func HideString(value string) {
	logger.hideString(value)
}

func SetFormat(format lorg.Formatter) {
	logger.SetFormat(format)
}

func SetIndentLines(value bool) {
	logger.SetIndentLines(value)
}

func SetLevel(level lorg.Level) {
	logger.SetLevel(level)
}

func GetWriter() io.Writer {
	return logger.writer
}

func Logger() *RedactedLog {
	return logger
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
