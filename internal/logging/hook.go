package logging

import (
	"MediaWarp/internal/config"
	"MediaWarp/utils"
	"os"

	"github.com/sirupsen/logrus"
)

type LoggerFileHook struct {
	isService bool
}

func NewLoggerFileHook(isService bool) *LoggerFileHook {
	return &LoggerFileHook{
		isService: isService,
	}
}

func (h *LoggerFileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// HOOK
//
// 将日志写入文件
func (h *LoggerFileHook) Fire(entry *logrus.Entry) error {
	if err := os.MkdirAll(config.LogDirWithDate(), os.ModePerm); err != nil {
		return err
	}
	var filename string
	if h.isService {
		filename = config.ServiceLogPath()
	} else {
		filename = config.AccessLogPath()
	}
	serviceLogFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer serviceLogFile.Close()

	line, err := entry.String()
	if err != nil {
		return err
	}
	serviceLogFile.WriteString(utils.RemoveColorCodes(line))
	return nil
}

var _ logrus.Hook = (*LoggerFileHook)(nil)
