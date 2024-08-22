package pkg

import (
	"regexp"
	"strings"
)

// RemoveColorCodes 移除字符串中的 ANSI 颜色代码
func RemoveColorCodes(line string) string {
	re := regexp.MustCompile(`\033\[[0-9]*m`)
	return re.ReplaceAllString(line, "")
}

// 将字符串居中
func Center(s string, width int, fill string) string {
	if len(s) >= width {
		return s
	}
	padding := width - len(s)
	leftPadding := padding / 2
	rightPadding := padding - leftPadding
	return strings.Repeat(fill, leftPadding) + s + strings.Repeat(fill, rightPadding)
}
