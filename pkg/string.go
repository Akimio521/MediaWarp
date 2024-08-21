package pkg

import (
	"regexp"
)

// RemoveColorCodes 移除字符串中的 ANSI 颜色代码
func RemoveColorCodes(line string) string {
	re := regexp.MustCompile(`\033\[[0-9]*m`)
	return re.ReplaceAllString(line, "")
}
