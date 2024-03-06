package log

func Debug(v ...any) {
	logger.Debug(v...)
}

func Debugf(format string, v ...any) {
	logger.Debugf(format, v...)
}

func Info(v ...any) {
	logger.Info(v...)
}

func Infof(format string, v ...any) {
	logger.Infof(format, v...)
}

func Warn(v ...any) {
	logger.Warn(v...)
}

func Warnf(format string, v ...any) {
	logger.Warnf(format, v...)
}

func Error(v ...any) {
	logger.Error(v...)
}

func Errorf(format string, v ...any) {
	logger.Errorf(format, v...)
}

func SetLevel(level int) {
	logger.SetLevel(level)
}

func SetLogger(v Logger) {
	logger = v
}
