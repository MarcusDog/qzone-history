package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// extractStringBetween 提取两个字符串之间的内容
func extractStringBetween(sourceString, startString, endString string) string {
	startIndex := strings.Index(sourceString, startString)
	if startIndex < 0 {
		return ""
	}
	startIndex += len(startString)

	relativeEndIndex := strings.Index(sourceString[startIndex:], endString)
	if relativeEndIndex < 0 {
		return ""
	}

	return sourceString[startIndex : startIndex+relativeEndIndex]
}

// replaceMultipleSpaces 去除多余的空格
func replaceMultipleSpaces(s string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(s, " ")
}

// ProcessOldHTML 替换十六进制编码并处理HTML
func ProcessOldHTML(message string) string {
	re := regexp.MustCompile(`\\x[0-9a-fA-F]{2}`)
	newText := re.ReplaceAllStringFunc(message, func(hex string) string {
		byteValue, err := strconv.ParseUint(hex[2:], 16, 8)
		if err != nil {
			return hex
		}
		return string([]byte{byte(byteValue)})
	})

	startString := "html:'"
	endString := "',opuin"
	newText = extractStringBetween(newText, startString, endString)
	newText = replaceMultipleSpaces(newText)
	newText = strings.NewReplacer(
		`\/`, `/`,
		`\"`, `"`,
		`\'`, `'`,
		`\\`, `\`,
	).Replace(newText)
	return newText
}
