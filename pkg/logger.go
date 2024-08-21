package pkg

import (
	"bytes"
	"fmt"
	"path"

	"github.com/sirupsen/logrus"
)

const (
	C_RED    int = 1
	C_YELLOW int = 3
	C_GREEN  int = 2
	C_BLUE   int = 4
	C_WRITE  int = 7
	C_GRAY   int = 8
)

type LoggerManager struct{}

var logger = logrus.New()

type loggerFormatter struct{}

func (f *loggerFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var color int
	switch entry.Level {
	case logrus.ErrorLevel:
		color = C_RED
	case logrus.WarnLevel:
		color = C_YELLOW
	case logrus.InfoLevel:
		color = C_GREEN
	case logrus.DebugLevel:
		color = C_WRITE
	default:
		color = C_GRAY
	}

	// 设置Buffer缓冲区
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// 设置时间格式
	formatTime := entry.Time.Format("2006-01-02 15:04:05")

	// 设置格式
	fmt.Fprintf(b, "\033[3%dm[%s]\t[%s]%s:%s\033[0m\n", color, entry.Level, formatTime, path.Base(entry.Caller.File), entry.Message)
	return b.Bytes(), nil
}

func InitLger() *logrus.Logger {
	logger.SetReportCaller(true)
	logger.SetFormatter(&loggerFormatter{})

	return logger
}
