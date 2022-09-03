package util

import (
	"strings"
	"unicode"
)

// UpperFirstLetter 所有首字母大小
func UpperFirstLetter(s string) string {
	temp := strings.Split(s, "_")
	var result string
	for y := 0; y < len(temp); y++ {
		vv := []rune(temp[y])
		for i := 0; i < len(vv); i++ {
			// 不是大写才转换
			if i == 0 && !unicode.IsUpper(vv[i]) {
				vv[i] -= 32
				result += string(vv[i]) // + string(vv[i+1])
			} else {
				result += string(vv[i]) // + string(vv[i+1])
			}
		}
	}
	return result
}

// Substr 字符串截取
func Substr(s string, start, end int) string {
	r := []rune(s)
	return string(r[start:end])
}

// HidePhone 手机号码 加 *
func HidePhone(s string) string {
	return Substr(s, 0, 3) + "****" + Substr(s, 7, len(s))
}

// HideIdentity 身份证号码 加 *
func HideIdentity(s string) string {
	return Substr(s, 0, 6) + "******" + Substr(s, 12, len(s))
}
