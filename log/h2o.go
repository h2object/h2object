package log

import (
	"strings"
	"os"
	"path/filepath"
	l4g "github.com/h2object/log4go"
)

type H2OLogger struct{
	logger l4g.Logger
}

func NewH2OLogger() *H2OLogger {
	return &H2OLogger{
		logger: make(l4g.Logger),
	}
}

func (log *H2OLogger) SetConsole(enable bool) {
	if enable {
		log.logger.AddFilter("console", l4g.TRACE, l4g.NewConsoleLogWriter())
	}	
}
func (log *H2OLogger) SetFileLog(fn string, level string, rotateMaxSize int, rotateMaxLine int, rotateDaily bool) {
	dir, _ := filepath.Split(fn)
	os.MkdirAll(dir, os.ModePerm)

	flw := l4g.NewFileLogWriter(fn, true)
	flw.SetFormat("[%D %T] [%L] %M")
	flw.SetRotateSize(rotateMaxSize)
	flw.SetRotateLines(rotateMaxLine)
	flw.SetRotateDaily(rotateDaily)
	switch strings.ToLower(level) {
	case "fine":
		log.logger.AddFilter("file", l4g.FINE, flw)
	case "debug":
		log.logger.AddFilter("file", l4g.DEBUG, flw)
	case "trace":
		log.logger.AddFilter("file", l4g.TRACE, flw)
	case "info":
		log.logger.AddFilter("file", l4g.INFO, flw)
	case "warn":
		log.logger.AddFilter("file", l4g.WARNING, flw)
	case "error":
		log.logger.AddFilter("file", l4g.ERROR, flw)
	case "critical":
		log.logger.AddFilter("file", l4g.CRITICAL, flw)
	}
}

func (log *H2OLogger) Close() {
	log.logger.Close()
}

func (log *H2OLogger) Trace(format string, args ...interface{}) {
	log.logger.Trace(format, args...)
}
func (log *H2OLogger) Debug(format string, args ...interface{}) {
	log.logger.Debug(format, args...)	
}
func (log *H2OLogger) Info(format string, args ...interface{}) {
	log.logger.Info(format, args...)	
}
func (log *H2OLogger) Warn(format string, args ...interface{}) {
	log.logger.Warn(format, args...)
} 
func (log *H2OLogger) Error(format string, args ...interface{}) {
	log.logger.Error(format, args...)
} 
func (log *H2OLogger) Critical(format string, args ...interface{}) {
	log.logger.Info(format, args...)
}

func (log *H2OLogger) Infof(format string, args ...interface{}) {
	log.logger.Info(format, args...)
}

func (log *H2OLogger) Warningf(format string, args ...interface{}) {
	log.logger.Warn(format, args...)
}

func (log *H2OLogger) Errorf(format string, args ...interface{}) {
	log.logger.Error(format, args...)
}

