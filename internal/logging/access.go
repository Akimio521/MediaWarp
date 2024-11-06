package logging

import (
	"MediaWarp/internal/config"
	"MediaWarp/utils"
	"bytes"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type accessLoggerSetting struct{}

// 实现Format方法
func (s *accessLoggerSetting) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer == nil {
		b = &bytes.Buffer{}
	} else {
		b = entry.Buffer
	}

	fmt.Fprint(
		b,
		entry.Message+"\n",
	)
	return b.Bytes(), nil
}

func (s *accessLoggerSetting) Levels() []logrus.Level {
	return logrus.AllLevels
}

// HOOK
//
// 将日志写入文件
func (s *accessLoggerSetting) Fire(entry *logrus.Entry) error {
	if err := os.MkdirAll(config.LogDirWithDate(), os.ModePerm); err != nil {
		return err
	}

	accessLogFile, err := os.OpenFile(config.AccessLogPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer accessLogFile.Close()

	line, err := entry.String()
	if err != nil {
		return err
	}
	accessLogFile.WriteString(utils.RemoveColorCodes(line))
	return nil
}
