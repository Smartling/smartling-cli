package rlog

func Infof(format string, value ...interface{}) {
	logger.Infof(format, value...)
}

func Error(value ...interface{}) {
	logger.Error(value...)
}

func Errorf(format string, value ...interface{}) {
	logger.Errorf(format, value...)
}

func Debugf(format string, value ...interface{}) {
	logger.Debugf(format, value...)
}
