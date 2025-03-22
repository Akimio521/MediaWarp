package logging

import (
	"MediaWarp/constants"
	"MediaWarp/internal/config"
	"MediaWarp/utils"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type serviceLoggerSetting struct{}

func (s *serviceLoggerSetting) Format(entry *logrus.Entry) ([]byte, error) {
	// 根据日志级别设置颜色
	var colorCode uint8
	switch entry.Level {
	case logrus.DebugLevel:
		colorCode = constants.COLOR_BLUE
	case logrus.InfoLevel:
		colorCode = constants.COLOR_GREEN
	case logrus.WarnLevel:
		colorCode = constants.COLOR_YELLOW
	case logrus.ErrorLevel:
		colorCode = constants.COLOR_RED
	default:
		colorCode = constants.COLOR_GRAY
	}

	// 设置文本Buffer
	var b *bytes.Buffer
	if entry.Buffer == nil {
		b = &bytes.Buffer{}
	} else {
		b = entry.Buffer
	}
	// 时间格式化
	formatTime := entry.Time.Format(constants.FORMATE_TIME)

	fmt.Fprintf(
		b,
		"\033[3%dm【%s】\033[0m\t%s | %s\n", // 长度需要算是上控制字符的长度
		colorCode,
		strings.ToUpper(entry.Level.String()),
		formatTime,
		entry.Message,
	)
	return b.Bytes(), nil
}

func (s *serviceLoggerSetting) Levels() []logrus.Level {
	return []logrus.Level{logrus.ErrorLevel, logrus.WarnLevel}
}

// HOOK
//
// 将日志写入文件
func (s *serviceLoggerSetting) Fire(entry *logrus.Entry) error {
	if err := os.MkdirAll(config.LogDirWithDate(), os.ModePerm); err != nil {
		return err
	}

	serviceLogFile, err := os.OpenFile(config.ServiceLogPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
