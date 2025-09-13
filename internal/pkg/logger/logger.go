package logger

import (
	"github.com/sirupsen/logrus"
)

type Level logrus.Level

// These are the different logging levels. You can set the logging level to log
// on your instance of logger
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = Level(logrus.PanicLevel)
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel = Level(logrus.FatalLevel)
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel = Level(logrus.ErrorLevel)
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel = Level(logrus.WarnLevel)
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel = Level(logrus.InfoLevel)
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel = Level(logrus.DebugLevel)
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel = Level(logrus.TraceLevel)
)

type Logger struct {
	impl *logrus.Logger
}

func New(level Level) *Logger {
	l := &Logger{
		impl: logrus.New(),
	}
	l.impl.SetFormatter(&logrus.JSONFormatter{})

	l.SetLevel(level)
	return l
}

func (l *Logger) SetLevel(level Level) {
	l.impl.SetLevel(logrus.Level(level))
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.impl.Printf(format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.impl.Logf(logrus.InfoLevel, format, args...)
}

func (l *Logger) Exit(code int) {
	l.impl.Exit(code)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.impl.Logf(logrus.ErrorLevel, format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.impl.Logf(logrus.FatalLevel, format, args...)
}
