package common

import (
	"io"
	"os"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/disgord/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type GlobalLogger struct {
	instance *zap.Logger
}

var (
	// Logger represents a globally usable logger.
	Logger *GlobalLogger

	// DisGordLogger represents a clone of Logger with a few specifications for Disgord.
	DisGordLogger *logger.LoggerZap
)

// InitLogger initializes the global logger.
// If the logfile already exists, all the output will be appended.
func InitLogger(debug bool) {
	file, err := os.OpenFile("disgotify.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	conf := zap.NewProductionConfig()

	if debug {
		conf.Development = true
		conf.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	writeSyncer := zapcore.AddSync(io.Writer(file))
	logger, _ := conf.Build(
		zap.ErrorOutput(writeSyncer),
		zap.AddCallerSkip(1),
	)

	Logger = &GlobalLogger{
		instance: logger,
	}
	DisGordLogger = disgord.DefaultLoggerWithInstance(logger.With(
		zap.String("lib", constant.Name),
		zap.String("ver", constant.Version)))
}

// getMessage is a slightly modified version of DisGord's logging wrapper for zap.
// All credit goes to its contributors.
func (l *GlobalLogger) getMessage(v ...interface{}) string {
	var message string
	for i := range v {
		var str string
		switch t := v[i].(type) {
		case string:
			str = t
		case error:
			str = t.Error()
		default:
			str = fmt.Sprint(v[i])
		}

		if message != "" {
			message += " " + str
		} else {
			message = str
		}
	}

	return message
}

// Debug logs a message at DebugLevel and flushes any buffered log entries.
func (l *GlobalLogger) Debug(v ...interface{}) {
	l.instance.Debug(l.getMessage(v...))
	_ = l.instance.Sync()
}

// Info logs a message at InfoLevel and flushes any buffered log entries.
func (l *GlobalLogger) Info(v ...interface{}) {
	l.instance.Info(l.getMessage(v...))
	_ = l.instance.Sync()
}

// Warn logs a message at WarnLevel and flushes any buffered log entries.
func (l *GlobalLogger) Warn(v ...interface{}) {
	l.instance.Warn(l.getMessage(v...))
	_ = l.instance.Sync()
}

// Error logs a message at ErrorLevel and flushes any buffered log entries.
func (l *GlobalLogger) Error(v ...interface{}) {
	l.instance.Error(l.getMessage(v...))
	_ = l.instance.Sync()
}

// Panic logs a message at PanicLevel, flushes any buffered log entries and panics.
func (l *GlobalLogger) Panic(v ...interface{}) {
	defer l.instance.Sync()
	l.instance.Panic(l.getMessage(v...))
}

// Fatal logs a message at FatalLevel, flushes any buffered log entries and calls os.Exit(1).
func (l *GlobalLogger) Fatal(v ...interface{}) {
	defer l.instance.Sync()
	l.instance.Fatal(l.getMessage(v...))
}
