package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"

	// "timeutils"
	// "time"
	timeutils "rttmas-backend/pkg/utils/timeutils"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type Formatter struct{}

var log = logrus.New()
var infoColor = color.New(color.BgCyan, color.FgBlack)
var errColor = color.New(color.BgRed, color.FgBlack)
var fatalColor = color.New(color.BgMagenta, color.FgWhite)
var debugColor = color.New(color.BgBlack, color.FgWhite)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
	// log.SetFormatter(&logrus.TextFormatter{
	// 	DisableColors: false,
	// 	// FullTimestamp: true,
	// 	TimestampFormat: "2006-01-02 15:04:05",
	// })

	// stores log to file
	log_filename := fmt.Sprintf("logs/%s.log", timeutils.GetDatetime())
	file, err := os.OpenFile(log_filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err == nil {
		log.SetOutput(io.MultiWriter(os.Stdout, file))
	} else {
		log.Info(fmt.Sprintf("Failed to log to file: %s", log_filename))
	}
	log.SetFormatter(&Formatter{})

}

// SetLogLevel sets the log level for the logger.
func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

// Debug logs a message at the debug level.
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Info logs a message at the info level.
func Info(args ...interface{}) {
	log.Info(args...)
}

// Warning logs a message at the warning level.
func Warning(args ...interface{}) {
	log.Warning(args...)
}

// Error logs a message at the error level.
func Error(args ...interface{}) {
	log.Error(args...)
}

// Fatal logs a message at the fatal level and exits the program.
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Panic logs a message at the panic level and panics.
func Panic(args ...interface{}) {
	log.Panic(args...)
}

func (m *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	var logStr string

	logStr = generateLogStr(*entry)

	b.WriteString(logStr)
	return b.Bytes(), nil
}

func generateLogStr(e logrus.Entry) string {
	var logStr string
	timestamp := e.Time.Format("2006/01/02 - 15:04:05")
	switch e.Level {
	case logrus.InfoLevel:
		logStr = fmt.Sprintf("[LOG] %s |%s| %s\n", timestamp, infoColor.Sprintf(" %s ", "INF"), e.Message)
	case logrus.ErrorLevel:
		logStr = fmt.Sprintf("[LOG] %s |%s| %s\n", timestamp, errColor.Sprintf(" %s ", "ERR"), e.Message)
	case logrus.FatalLevel:
		logStr = fmt.Sprintf("[LOG] %s |%s| %s\n", timestamp, fatalColor.Sprintf(" %s ", "FAT"), e.Message)
	case logrus.DebugLevel:
		logStr = fmt.Sprintf("[LOG] %s |%s| %s\n", timestamp, debugColor.Sprintf(" %s ", "DBG"), e.Message)
	default:
		logStr = fmt.Sprintf("[LOG] %s |%s| %s\n", timestamp, debugColor.Sprintf(" %s ", e.Level), e.Message)

	}
	return logStr
}
