package rlog

// Infof does logger.Infof
func Infof(format string, value ...any) {
	logger.Infof(format, value...)
}

// Error does logger.Error
func Error(value ...any) {
	logger.Error(value...)
}

// Errorf does logger.Errorf
func Errorf(format string, value ...any) {
	logger.Errorf(format, value...)
}

// Debugf does logger.Debugf
func Debugf(format string, value ...any) {
	logger.Debugf(format, value...)
}
