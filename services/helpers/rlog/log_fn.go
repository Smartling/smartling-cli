package rlog

// Infof does logger.Infof
func Infof(format string, value ...interface{}) {
	logger.Infof(format, value...)
}

// Error does logger.Error
func Error(value ...interface{}) {
	logger.Error(value...)
}

// Errorf does logger.Errorf
func Errorf(format string, value ...interface{}) {
	logger.Errorf(format, value...)
}

// Debugf does logger.Debugf
func Debugf(format string, value ...interface{}) {
	logger.Debugf(format, value...)
}
