package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

func InitLogger(level string) *logrus.Logger {
	log := logrus.New()
	log.SetReportCaller(true)
	log.Out = os.Stdout
	switch level {
	case "error":
		log.Level = logrus.ErrorLevel
	default:
		log.Level = logrus.InfoLevel
	}

	log.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s():%d", frame.Function, frame.Line), filename
		},
		DisableColors:  true,
		DisableSorting: false,
	}

	return log
}
