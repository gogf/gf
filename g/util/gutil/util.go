package gutil

// 便利数组查找字符串索引位置，如果不存在则返回-1
func StringSearch (a []string, s string) int {
    for i, v := range a {
        if s == v {
            return i
        }
    }
    return -1
}

// 判断字符串是否在数组中
func StringInArray (a []string, s string) bool {
    return StringSearch(a, s) != -1
}



