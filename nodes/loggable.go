package nodes

////////////////////////////////////////////////////////////////////////////////
type Loggable struct {
	Logger *Logger
}

////////////////////////////////////////////////////////////////////////////////
func (l *Loggable) SetLogger(logger *Logger) {
	l.Logger = logger
}