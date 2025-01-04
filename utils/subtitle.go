package utils

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	defaultRegularBodyWeight = 400 // 默认正常体字重
	defaultBoldBodyWeight    = 700 // 默认加粗体字重
	ASSHeader1               = `[Script Info]
; This is an Advanced Sub Station Alpha v4+ script.
Title:
ScriptType: v4.00+
Collisions: Normal
PlayDepth: 0

[V4+ Styles]`
	ASSHeader2 = `[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text`
)

var (
	srtSubtitlesPattern      = regexp.MustCompile(`@\d+@\d{2}:\d{2}:\d{2},\d{3} --> \d{2}:\d{2}:\d{2},\d{3}@`) // 用于文本是否为 SRT 字幕的正则表达式
	srtTimePattern           = regexp.MustCompile(`-?\d\d:\d\d:\d\d`)
	timeFormatPattern        = regexp.MustCompile(`\d(\d:\d{2}:\d{2}),(\d{2})\d`)
	arrowPattern             = regexp.MustCompile(`\s+-->\s+`)
	styleTagStartPattern     = regexp.MustCompile(`<([ubi])>`)
	styleTagEndPattern       = regexp.MustCompile(`</([ubi])>`)
	fontColorTagStartPattern = regexp.MustCompile(`<font\s+color="?#[\w]{2}([\w]{2})([\w]{2})([\w]{2})"?\s*>`)
	fontColorTagEndPattern   = regexp.MustCompile(`</font>`)
)

// 判断字幕是否为 SRT 格式
func IsSRT(content []byte) bool {
	content = bytes.ReplaceAll(content, []byte{'\r'}, []byte{})    // 去除 \r 保证多系统兼容
	content = bytes.ReplaceAll(content, []byte{'\n'}, []byte{'@'}) // 将 \n 替换为 @
	return srtSubtitlesPattern.Match(content)                      // 查找第一个匹配项
}

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
		subtitleContent        string                       // 字幕内容
		subtitleBuffer         []rune = make([]rune, 1, 50) // 字幕缓存区（某一行字幕未完成先存取到此处）
		currentSubtitleContent uint8  = 0                   // 一个时间下字幕的行数（2表示这一时间有2行字幕）
	)
	for index, line := range lines {
		if isInt(line) && srtTimePattern.MatchString(lines[index+1]) { // 这一行是 SRT 字幕的序列数且下一行是时间
			if len(subtitleBuffer) > 0 {
				subtitleContent += string(append(subtitleBuffer, '\n'))
				subtitleBuffer = subtitleBuffer[:0]
			}
			currentSubtitleContent = 0
		} else {
			if srtTimePattern.MatchString(line) { // 这一行是时间行
				subtitleBuffer = append(subtitleBuffer, []rune("Dialogue: 0,")...)
				subtitleBuffer = append(subtitleBuffer, []rune(strings.ReplaceAll(line, "-0", "0"))...)
				subtitleBuffer = append(subtitleBuffer, []rune(",Default,,0,0,0,,")...)
			} else {
				if currentSubtitleContent == 0 {
					subtitleBuffer = append(subtitleBuffer, []rune(line)...)
				} else {
					subtitleBuffer = append(subtitleBuffer, []rune(`\n`)...) // 同一时间多行字幕需要在一行中使用字面量 \n 表示换行
					subtitleBuffer = append(subtitleBuffer, []rune(line)...)
				}
				currentSubtitleContent += 1
			}
		}
	}
	subtitleContent += string(append(subtitleBuffer, '\n')) // 最后一行字幕

	subtitleContent = timeFormatPattern.ReplaceAllString(subtitleContent, "$1.$2")                 // 替换时间格式
	subtitleContent = arrowPattern.ReplaceAllString(subtitleContent, ",")                          // 替换箭头符号
	subtitleContent = styleTagStartPattern.ReplaceAllString(subtitleContent, `{\\$11}`)            // 替换样式标签
	subtitleContent = styleTagEndPattern.ReplaceAllString(subtitleContent, `{\\$10}`)              // 替换字体颜色标签
	subtitleContent = fontColorTagStartPattern.ReplaceAllString(subtitleContent, `{\\c&H$3$2$1&}`) // 替换字体颜色标签
	subtitleContent = fontColorTagEndPattern.ReplaceAllString(subtitleContent, "")                 // 删除字体结束标签
	return ASSHeader1 + "\n" + strings.Join(style, "\n") + "\n\n" + ASSHeader2 + "\n\n" + subtitleContent
}

type ASSFontStyle struct {
	Name   string // 字体名字
	Weight uint16 // 字重
	Italic bool   // 是否为意大利体（斜体）
}

// 分析 ASS 字幕中有哪些“字”，用于字体子集化
func AnalyseASS(assText string) (map[ASSFontStyle]SetInterface[rune], error) {
	var (
		state              uint8                               = 0 // 文本状态控制字
		assStyleNameIndex  int8                                = -1
		assFontNameIndex   int8                                = -1
		assBodyIndex       int8                                = -1
		assItalicIndex     int8                                = -1
		fontStyles         map[string]ASSFontStyle             = make(map[string]ASSFontStyle, 1)
		allFontStyleName   []string                            = []string{}
		firstFontStyleName string                              = ""
		assEventTextIndex  int8                                = -1
		assEventStyleIndex int8                                = -1
		subFontSets        map[ASSFontStyle]SetInterface[rune] = make(map[ASSFontStyle]SetInterface[rune], 1)
	)
	assText = strings.ReplaceAll(assText, "\r", "")
	for _, line := range strings.Split(assText, "\n") {
		if line == "" {
			continue
		} else if state == 0 && strings.HasPrefix(line, "[V4+ Styles]") {
			state = 1
		} else if state == 1 { // 这一行是字体格式
			if !strings.HasPrefix(line, "Format:") {
				return nil, fmt.Errorf("解析 ASS 字体 Styles 格式失败：%s", line)
			}
			stylesFormat := strings.Split(strings.TrimSpace(strings.Replace(line, "Format:", "", 1)), ",")
			assStyleNameIndex = int8(FindStringIndex(stylesFormat, "Name", true, true))    // ASS 中字体样式的名字
			assFontNameIndex = int8(FindStringIndex(stylesFormat, "Fontname", true, true)) // ASS 中样式使用字体的名字
			if assStyleNameIndex == -1 {
				return nil, fmt.Errorf("未找字体格式中的 Name ：%s", stylesFormat)
			}
			if assFontNameIndex == -1 {
				return nil, fmt.Errorf("未找字体格式中的 Fontname：%s", stylesFormat)
			}
			assBodyIndex = int8(FindStringIndex(stylesFormat, "Bold", true, true))
			assItalicIndex = int8(FindStringIndex(stylesFormat, "Italic", true, true))
			state = 2
		} else if state == 2 { // 这一行开始是字体样式
			if strings.HasPrefix(line, "Style:") {
				styleData := strings.Split(strings.TrimSpace(strings.Replace(line, "Style:", "", 1)), ",")
				styleName := strings.ReplaceAll(strings.TrimSpace(styleData[assStyleNameIndex]), "*", "")
				fontName := strings.ReplaceAll(strings.TrimSpace(styleData[assFontNameIndex]), "@", "")
				var (
					fontWeight uint16 = defaultRegularBodyWeight
					italic     bool   = false
				) // 字重默认 400
				if assBodyIndex != -1 && strings.TrimSpace(styleData[assBodyIndex]) == "1" { // 当该 ASS 字幕格式存在 Bold 属性且该样式属性设置为 "1" 时，将字重设置为 700
					fontWeight = defaultBoldBodyWeight
				}
				if assItalicIndex != -1 && strings.TrimSpace(styleData[assItalicIndex]) == "1" {
					italic = true
				}
				fontStyles[styleName] = ASSFontStyle{fontName, fontWeight, italic} // 将该样式存入 map 中
				allFontStyleName = append(allFontStyleName, styleName)
				if firstFontStyleName == "" {
					firstFontStyleName = styleName
				}
			} else if strings.HasPrefix(line, "[Events]") {
				state = 3 // 字体样式已结束
			}
		} else if state == 3 { // 开始解析 Event
			if !strings.HasPrefix(line, "Format:") {
				return nil, fmt.Errorf("解析字幕事件格式失败：%s", line)
			}
			eventFormat := strings.Split(strings.ReplaceAll(strings.Replace(line, "Format:", "", 1), " ", ""), ",")
			assEventTextIndex = int8(FindStringIndex(eventFormat, "Text", true, true))
			assEventStyleIndex = int8(FindStringIndex(eventFormat, "Style", true, true))
			if assEventTextIndex == -1 || assEventStyleIndex == -1 {
				return nil, fmt.Errorf("字幕事件格式中未找到 Text 或 Style：%s", line)
			}
			if int(assEventTextIndex) != len(eventFormat)-1 {
				return nil, fmt.Errorf("字幕事件格式中 Text 不是最后一个元素：%s", line)
			}
			state = 4
		} else if state == 4 { // 开始解析字幕具体内容
			if strings.HasPrefix(line, "Dialogue:") {
				parts := strings.Split(line, ",")
				text := strings.Join(parts[assEventTextIndex:], ",")                       // 获取字幕文本，需要考虑字幕中带有英文逗号的可能性
				defaultStyleName := strings.ReplaceAll(parts[assEventStyleIndex], "*", "") // 当前行默认样式
				if !Contains(allFontStyleName, defaultStyleName) {
					defaultStyleName = firstFontStyleName // 未找到对应样式，使用第一个样式
				}
				defaultStyle := fontStyles[defaultStyleName] // 当前样式
				currentStyle := defaultStyle                 // 当前样式
				parseTag := func(tag []rune) error {
					length := len(tag)
					if len(tag) == 0 {
						return nil
					}
					if length == 1 && tag[0] == 'r' {
						currentStyle = defaultStyle
					} else if length > 2 && tag[0] == 'f' && tag[1] == 'n' {
						// fmt.Println("Tag：", string(tag), "字体名：", strings.ReplaceAll(string(tag[2:]), "@", ""))
						currentStyle.Name = strings.ReplaceAll(string(tag[2:]), "@", "")
					} else if length == 2 && tag[0] == 'i' {
						if tag[1] == '1' {
							currentStyle.Italic = true
						} else if tag[1] == '0' {
							currentStyle.Italic = false
						} else {
							return fmt.Errorf("未知斜体状态：%s", string(tag))
						}
					} else if tag[0] == 'b' {
						if length == 2 {
							if tag[1] == '0' {
								currentStyle.Weight = defaultRegularBodyWeight
							} else if tag[1] == '1' {
								currentStyle.Weight = defaultBoldBodyWeight
							} else {
								return fmt.Errorf("未知加粗状态：%s", string(tag))
							}
						} else {
							for i := range tag[1:] {
								if tag[i] < '0' || tag[i] > '9' {
									return nil // b 后面不全是数字，不是加粗标签，忽略
								}
							}
							num, err := strconv.Atoi(string(tag[1:]))
							if err != nil {
								return fmt.Errorf("解析加粗数值失败：%s", string(tag))
							}
							currentStyle.Weight = uint16(num)
						}
					}
					return nil
				}

				var buffer []rune = make([]rune, 0, len(line))
				for _, char := range text {
					if char == '{' { // 可能是特殊样式的开始
						if subFontSets[currentStyle] == nil {
							subFontSets[currentStyle] = NewSet[rune]()
						}
						subFontSets[currentStyle].Adds(buffer...)
					} else if char == '}' && len(buffer) > 1 && buffer[0] == '{' && buffer[1] == '\\' {
						var tagStartIndex uint16 = 2
						for i, c := range buffer[2:] {
							if c == '\\' {
								err := parseTag(buffer[tagStartIndex : 2+i])
								if err != nil {
									return nil, err
								}
								tagStartIndex = uint16(2 + i + 1)
							}
						}
						err := parseTag(buffer[tagStartIndex:])
						if err != nil {
							return nil, err
						}
					} else {
						buffer = append(buffer, char)
					}
				}
				if len(buffer) != 0 {
					if subFontSets[currentStyle] == nil {
						subFontSets[currentStyle] = NewSet[rune]()
					}
					subFontSets[currentStyle].Adds(buffer...)
				}

			}
		}

	}
	return subFontSets, nil
}
