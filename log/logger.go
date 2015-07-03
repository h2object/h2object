package log

type Logger interface{
	Trace(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{}) 
	Error(format string, args ...interface{}) 
	Critical(format string, args ...interface{})
}




