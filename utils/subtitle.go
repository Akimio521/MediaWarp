package utils

import (
	"regexp"
	"strings"
)

var srtRegexp = regexp.MustCompile(`@\d+@\d{2}:\d{2}:\d{2},\d{3} --> \d{2}:\d{2}:\d{2},\d{3}@`)

// 判断字幕是否为 SRT 格式
func IsSRT(text string) bool {
	text = strings.ReplaceAll(text, "\r", "")                        // 去除 \r 保证多系统兼容
	lines := strings.Split(text, "\n")                               // 按行分割
	matches := srtRegexp.FindAllString(strings.Join(lines, "@"), -1) // 查找所有匹配项
	return len(matches) > 0
}

var (
	srtTimeRegxp    = regexp.MustCompile(`-?\d\d:\d\d:\d\d`)
	timeFormatRe    = regexp.MustCompile(`\d(\d:\d{2}:\d{2}),(\d{2})\d`)
	arrowRe         = regexp.MustCompile(`\s+-->\s+`)
	styleOpenTagRe  = regexp.MustCompile(`<([ubi])>`)
	styleCloseTagRe = regexp.MustCompile(`</([ubi])>`)
	fontColorTagRe  = regexp.MustCompile(`<font\s+color="?#[\w]{2}([\w]{2})([\w]{2})([\w]{2})"?\s*>`)
	fontCloseTagRe  = regexp.MustCompile(`</font>`)
)

const (
	ASSHeader1 = `[Script Info]
; This is an Advanced Sub Station Alpha v4+ script.
Title:
ScriptType: v4.00+
Collisions: Normal
PlayDepth: 0

[V4+ Styles]`
	ASSHeader2 = `[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text`
)

// 将 SRT 字幕转换成 ASS 字幕
//
// srtText: SRT 格式字幕文本
// style: ASS 字幕样式
func SRT2ASS(srtText string, style []string) string {
	srtText = strings.ReplaceAll(srtText, "\r", "")
	var lines []string
	for _, line := range strings.Split(srtText, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	var (
		subtitleContent        string = "" // 字幕内容
		subtitleBuffer         string = "" // 字幕缓存区（某一行字幕未完成先存取到此处）
		currentSubtitleContent uint8  = 0  // 一个时间下字幕的行数（2表示这一时间有2行字幕）
	)
	for index, line := range lines {
		if isInt(line) && srtTimeRegxp.Match([]byte(lines[index+1])) { // 这一行是 SRT 字幕的序列数且下一行是时间
			if subtitleBuffer != "" {
				subtitleContent += subtitleBuffer + "\n"
				subtitleBuffer = ""
			}
			currentSubtitleContent = 0
		} else {
			if srtTimeRegxp.Match([]byte(line)) {
				line = strings.ReplaceAll(line, "-0", "0")
				subtitleBuffer += "Dialogue: 0," + line + ",Default,,0,0,0,,"
			} else {
				if currentSubtitleContent == 0 {
					subtitleBuffer += line
				} else {
					subtitleBuffer += "\\n" + line // 同一时间多行字幕需要在一行中使用字面量 \n 表示换行
				}
				currentSubtitleContent += 1
			}
		}
	}
	subtitleContent += subtitleBuffer + "\n" // 最后一行字幕

	subtitleContent = timeFormatRe.ReplaceAllString(subtitleContent, "$1.$2")            // 替换时间格式
	subtitleContent = arrowRe.ReplaceAllString(subtitleContent, ",")                     // 替换箭头符号
	subtitleContent = styleOpenTagRe.ReplaceAllString(subtitleContent, `{\\$11}`)        // 替换样式标签
	subtitleContent = styleCloseTagRe.ReplaceAllString(subtitleContent, `{\\$10}`)       // 替换字体颜色标签
	subtitleContent = fontColorTagRe.ReplaceAllString(subtitleContent, `{\\c&H$3$2$1&}`) // 替换字体颜色标签
	subtitleContent = fontCloseTagRe.ReplaceAllString(subtitleContent, "")               // 删除字体结束标签
	return ASSHeader1 + "\n" + strings.Join(style, "\n") + "\n\n" + ASSHeader2 + "\n\n" + subtitleContent
}
