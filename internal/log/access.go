package log

import (
	"MediaWarp/pkg"
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

func (s *accessLoggerSetting) Fire(entry *logrus.Entry) error {
	accessLogFile, err := os.OpenFile(cfg.AccessLogPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer accessLogFile.Close()
	line, err := entry.String()
	if err != nil {
		return err
	}
	accessLogFile.WriteString(pkg.RemoveColorCodes(line))
	return nil
}
