package epiclogger

import (
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var (
	// baseLogger is the name of the standard logger in baseLoggerlib `log`
	baseLogger = &EpicLogger{log.NewEntry(log.New())}
	contextKey = "context"
)

// Base returns the default Logger logging to
func Base() *EpicLogger {
	return baseLogger
}

// NewLogger returns a new Logger logging to out.
func NewEpicLogger(w io.Writer) EpicLogger {
	l := log.New()
	l.Out = w
	return EpicLogger{log.NewEntry(l)}
}

func SetFormatter(formatter log.Formatter) {
	baseLogger.Logger.Formatter = formatter
}

// SetLevel sets the standard logger level.
func SetLevel(level log.Level) {
	baseLogger.Logger.SetLevel(level)
}

// AddHook adds a hook to the standard logger hooks.
func AddHook(hook log.Hook) {
	baseLogger.Logger.Hooks.Add(hook)
}

type EpicLogger struct {
	*log.Entry
}

func (e *EpicLogger) WithCtx(ctx context.Context) *EpicLogger {
	return &EpicLogger{e.Entry.WithField(contextKey, ctx)}
}

func (e *EpicLogger) WithField(key string, value interface{}) *EpicLogger {
	return &EpicLogger{e.Entry.WithField(key, value)}
}

func (e *EpicLogger) WithError(err error) *EpicLogger {
	return &EpicLogger{e.Entry.WithError(err)}
}

func (e *EpicLogger) WithFields(fields log.Fields) *EpicLogger {
	return &EpicLogger{e.Entry.WithFields(fields)}
}

func (e *EpicLogger) addServiceContext() *EpicLogger {
	podName := os.Getenv("POD_NAME")
	serviceNVersion := strings.Split(podName, "-")
	length := len(serviceNVersion)
	if podName != "" {
		service := strings.Join(serviceNVersion[:length-2], "-")
		version := serviceNVersion[length-2]
		return e.WithField("service", service).WithField("version", version)
	}
	return e
}

// Error logs a message at level Error on the standard logger.
func (e *EpicLogger) Error(args ...interface{}) {
	e.addServiceContext().Entry.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func (e *EpicLogger) Panic(args ...interface{}) {
	e.addServiceContext().Entry.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func (e *EpicLogger) Fatal(args ...interface{}) {
	e.addServiceContext().Entry.Fatal(args...)
}

// Errorf logs a message at level Error on the standard logger.
func (e *EpicLogger) Errorf(format string, args ...interface{}) {
	e.addServiceContext().Entry.Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func (e *EpicLogger) Panicf(format string, args ...interface{}) {
	e.addServiceContext().Entry.Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func (e *EpicLogger) Fatalf(format string, args ...interface{}) {
	e.addServiceContext().Entry.Fatalf(format, args...)
}

// Errorln logs a message at level Error on the standard logger.
func (e *EpicLogger) Errorln(args ...interface{}) {
	e.addServiceContext().Entry.Errorln(args...)
}

// Panicln logs a message at level Panic on the standard logger.
func (e *EpicLogger) Panicln(args ...interface{}) {
	e.addServiceContext().Entry.Panicln(args...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func (e *EpicLogger) Fatalln(args ...interface{}) {
	e.addServiceContext().Entry.Fatalln(args...)
}

// WithError creates an entry from the standard logger and adds a context to it, using the value defined in contextKey as key.
func WithCtx(ctx context.Context) *EpicLogger {
	return baseLogger.WithCtx(ctx)
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *EpicLogger {
	return baseLogger.WithField(log.ErrorKey, err)
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *EpicLogger {
	return baseLogger.WithField(key, value)
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields log.Fields) *EpicLogger {
	return baseLogger.WithFields(fields)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	baseLogger.Debug(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	baseLogger.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	baseLogger.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	baseLogger.Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	baseLogger.Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	baseLogger.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	baseLogger.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	baseLogger.Fatal(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	baseLogger.Debugf(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	baseLogger.Printf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	baseLogger.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	baseLogger.Warnf(format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	baseLogger.Warningf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	baseLogger.Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	baseLogger.Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	baseLogger.Fatalf(format, args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	baseLogger.Debugln(args...)
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	baseLogger.Println(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	baseLogger.Infoln(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	baseLogger.Warnln(args...)
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	baseLogger.Warningln(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	baseLogger.Errorln(args...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	baseLogger.Panicln(args...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	baseLogger.Fatalln(args...)
}
