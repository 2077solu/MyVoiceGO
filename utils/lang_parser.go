package utils

import (
	"regexp"
	"unicode"
)

// TextLanguage 表示文本的语言类型
type TextLanguage string

const (
	Chinese  TextLanguage = "zh"
	Japanese TextLanguage = "ja"
	English  TextLanguage = "en"
	Mixed    TextLanguage = "mixed"
	Unknown  TextLanguage = "unknown"
)

// DetectLanguage 检测文本的主要语言
func DetectLanguage(text string) TextLanguage {
	if text == "" {
		return Unknown
	}

	// 统计各种语言字符的数量
	chineseCount := 0
	japaneseCount := 0
	englishCount := 0
	totalCount := 0

	// 中文正则表达式（包括简体和繁体）
	chineseRegex := regexp.MustCompile(`[一-鿿㐀-䶿𠀀-𪛟𪜀-𫜿𫝀-𫠟𫠠-𬺯豈-﫿㌀-㏿︰-﹏豈-﫿丽-𯨟]`)

	// 日语平假名和片假名正则表达式
	japaneseRegex := regexp.MustCompile(`[぀-ゟ゠-ヿ]`)

	// 英文正则表达式
	englishRegex := regexp.MustCompile(`[a-zA-Z]`)

	// 遍历文本中的每个字符
	for _, r := range text {
		if unicode.Is(unicode.Han, r) || chineseRegex.MatchString(string(r)) {
			chineseCount++
		} else if japaneseRegex.MatchString(string(r)) {
			japaneseCount++
		} else if englishRegex.MatchString(string(r)) {
			englishCount++
		}
		totalCount++
	}

	// 计算各语言字符的百分比
	chinesePercent := float64(chineseCount) / float64(totalCount) * 100
	japanesePercent := float64(japaneseCount) / float64(totalCount) * 100
	englishPercent := float64(englishCount) / float64(totalCount) * 100

	// 判断主要语言
	// 如果某种语言占比超过60%，则认为该语言是主要语言
	if chinesePercent > 60 {
		return Chinese
	} else if japanesePercent > 60 {
		return Japanese
	} else if englishPercent > 60 {
		return English
	}

	// 如果没有明显的主要语言，则判断是否为混合语言
	// 如果至少有两种语言的占比都超过20%，则认为是混合语言
	languageCount := 0
	if chinesePercent > 20 {
		languageCount++
	}
	if japanesePercent > 20 {
		languageCount++
	}
	if englishPercent > 20 {
		languageCount++
	}

	if languageCount >= 2 {
		return Mixed
	}

	// 如果仍然无法确定，则返回未知
	return Unknown
}

// DetectLanguagesWithPercentage 返回文本中各语言的详细信息和百分比
func DetectLanguagesWithPercentage(text string) map[TextLanguage]float64 {
	if text == "" {
		return map[TextLanguage]float64{
			Chinese:  0,
			Japanese: 0,
			English:  0,
		}
	}

	// 统计各种语言字符的数量
	chineseCount := 0
	japaneseCount := 0
	englishCount := 0
	totalCount := 0

	// 中文正则表达式（包括简体和繁体）
	chineseRegex := regexp.MustCompile(`[一-鿿㐀-䶿𠀀-𪛟𪜀-𫜿𫝀-𫠟𫠠-𬺯豈-﫿㌀-㏿︰-﹏豈-﫿丽-𯨟]`)

	// 日语平假名和片假名正则表达式
	japaneseRegex := regexp.MustCompile(`[぀-ゟ゠-ヿ]`)

	// 英文正则表达式
	englishRegex := regexp.MustCompile(`[a-zA-Z]`)

	// 遍历文本中的每个字符
	for _, r := range text {
		if unicode.Is(unicode.Han, r) || chineseRegex.MatchString(string(r)) {
			chineseCount++
		} else if japaneseRegex.MatchString(string(r)) {
			japaneseCount++
		} else if englishRegex.MatchString(string(r)) {
			englishCount++
		}
		totalCount++
	}

	// 计算各语言字符的百分比
	chinesePercent := float64(chineseCount) / float64(totalCount) * 100
	japanesePercent := float64(japaneseCount) / float64(totalCount) * 100
	englishPercent := float64(englishCount) / float64(totalCount) * 100

	return map[TextLanguage]float64{
		Chinese:  chinesePercent,
		Japanese: japanesePercent,
		English:  englishPercent,
	}
}
