package nodes

import (
	"fmt"

	"github.com/denkhaus/tcgl/applog"
	"github.com/segmentio/go-loggly"
)

////////////////////////////////////////////////////////////////////////////////
// Log levels to control the logging output.
const (
	LevelDebug = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

////////////////////////////////////////////////////////////////////////////////
type Logger struct {
	Loggly *loggly.Client
	Level  int
	Token  string
}

////////////////////////////////////////////////////////////////////////////////
// Debugf logs a message at debug level.
func (l *Logger) Debugf(format string, args ...interface{}) {
	applog.Debugf(format, args...)

	msg := loggly.Message{}
	msg["message"] = fmt.Sprintf(format, args...)
	l.Loggly.Debug(l.Token, msg)
}

////////////////////////////////////////////////////////////////////////////////
// Infof logs a message at info level.
func (l *Logger) Infof(format string, args ...interface{}) {
	applog.Infof(format, args...)

	msg := loggly.Message{}
	msg["message"] = fmt.Sprintf(format, args...)
	l.Loggly.Info(l.Token, msg)
}

////////////////////////////////////////////////////////////////////////////////
// Warningf logs a message at warning level.
func (l *Logger) Warningf(format string, args ...interface{}) {
	applog.Warningf(format, args...)

	msg := loggly.Message{}
	msg["message"] = fmt.Sprintf(format, args...)
	l.Loggly.Warn(l.Token, msg)
}

////////////////////////////////////////////////////////////////////////////////
// Errorf logs a message at error level.
func (l *Logger) Errorf(format string, args ...interface{}) {
	applog.Errorf(format, args...)

	msg := loggly.Message{}
	msg["message"] = fmt.Sprintf(format, args...)
	l.Loggly.Error(l.Token, msg)
}

////////////////////////////////////////////////////////////////////////////////
// Criticalf logs a message at critical level.
func (l *Logger) Criticalf(format string, args ...interface{}) {
	applog.Criticalf(format, args...)

	msg := loggly.Message{}
	msg["message"] = fmt.Sprintf(format, args...)
	l.Loggly.Critical(l.Token, msg)
}

////////////////////////////////////////////////////////////////////////////////
func (l *Logger) SetLevel(level int) {
	applog.SetLevel(level)
	l.Level = level
}

////////////////////////////////////////////////////////////////////////////////
func (e *Engine) GetLogger(token string) *Logger {
	logger := &Logger{Level: LevelDebug, Token: token}
	applog.SetLevel(applog.LevelDebug)

	if token, err := e.Config.GetLogglyTokenById("default"); err != nil {
		panic(err)
	} else {
		logger.Loggly = loggly.New(token)
	}

	return logger
}
