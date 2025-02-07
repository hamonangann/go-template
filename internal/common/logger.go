package common

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Log Logger

type Logger interface {
	Debug(args ...any)
	Debugf(format string, args ...any)
	Info(args ...any)
	Infof(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	WithField(key string, value any) Logger
	WithFields(fields map[string]any) Logger
}

func SetLogger(log Logger) {
	Log = log
}

type LogrusEntry struct {
	entry *logrus.Entry
}

func (l *LogrusEntry) Debugf(format string, args ...any) {
	l.entry.Debugf(format, args...)
}

func (l *LogrusEntry) Debug(args ...any) {
	l.entry.Debug(args...)
}

func (l *LogrusEntry) Infof(format string, args ...any) {
	l.entry.Infof(format, args...)
}

func (l *LogrusEntry) Info(args ...any) {
	l.entry.Info(args...)
}

func (l *LogrusEntry) Warn(args ...any) {
	l.entry.Warn(args...)
}

func (l *LogrusEntry) Warnf(format string, args ...any) {
	l.entry.Warnf(format, args...)
}

func (l *LogrusEntry) Errorf(format string, args ...any) {
	l.entry.Errorf(format, args...)
}

func (l *LogrusEntry) Error(args ...any) {
	l.entry.Error(args...)
}

func (l *LogrusEntry) Fatalf(format string, args ...any) {
	l.entry.Fatalf(format, args...)
}

func (l *LogrusEntry) Fatal(args ...any) {
	l.entry.Fatal(args...)
}

func (l *LogrusEntry) WithField(key string, value any) Logger {
	return &LogrusEntry{l.entry.WithField(key, value)}
}

func (l *LogrusEntry) WithFields(args map[string]any) Logger {
	return &LogrusEntry{l.entry.WithFields(args)}
}

type LogrusLogger struct {
	log *logrus.Logger
}

func NewLogrusLogger() *LogrusLogger {
	log := logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		PrettyPrint:     true,
	})

	log.SetLevel(logrus.InfoLevel)
	log.SetOutput(os.Stdout)

	return &LogrusLogger{log: log}
}

func (l *LogrusLogger) Debugf(format string, args ...any) {
	l.log.Debugf(format, args...)
}

func (l *LogrusLogger) Debug(args ...any) {
	l.log.Debug(args...)
}

func (l *LogrusLogger) Infof(format string, args ...any) {
	l.log.Infof(format, args...)
}

func (l *LogrusLogger) Info(args ...any) {
	l.log.Info(args...)
}

func (l *LogrusLogger) Warn(args ...any) {
	l.log.Warn(args...)
}

func (l *LogrusLogger) Warnf(format string, args ...any) {
	l.log.Warnf(format, args...)
}

func (l *LogrusLogger) Errorf(format string, args ...any) {
	l.log.Errorf(format, args...)
}

func (l *LogrusLogger) Error(args ...any) {
	l.log.Error(args...)
}

func (l *LogrusLogger) Fatalf(format string, args ...any) {
	l.log.Fatalf(format, args...)
}

func (l *LogrusLogger) Fatal(args ...any) {
	l.log.Fatal(args...)
}

func (l *LogrusLogger) WithField(key string, value any) Logger {
	return &LogrusEntry{entry: l.log.WithField(key, value)}
}

func (l *LogrusLogger) WithFields(args map[string]any) Logger {
	return &LogrusEntry{entry: l.log.WithFields(args)}
}
