package gregx

import (
    "regexp"
)

// 正则表达式是否匹配
func IsMatch(val, pattern string) bool {
    match, err := regexp.Match(pattern, []byte(val))
    if err != nil {
        return false
    }
    return match
}