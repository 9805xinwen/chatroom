package utils

import (
	"errors"
	"regexp"
)

////////////////////////////////////////////////////////////////////////
//                            字符串操作实现                             //
//--------------------------------------------------------------------//
// 内部函数:                                                            //
// DoubleQuotedStringsMarch(str string)                               //
////////////////////////////////////////////////////////////////////////

//寻找引号内的字符串的正则表达式
var stringExtractRegexpCompile *regexp.Regexp = regexp.MustCompile("\"([^\"]*)\"")

func DoubleQuotedStringsMarch(str string) (string, error) {
	result := stringExtractRegexpCompile.FindAllStringSubmatch(str, 1)
	if len(result) > 0 && len(result[0]) > 1 {
		return result[0][1], nil
	} else {
		return "", errors.New("没有匹配到引号内的内容")
	}
}
