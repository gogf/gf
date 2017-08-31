package gutil


// 合并两个slice
func MergeSlice(dst []byte, src []byte) []byte {
    if len(dst) == 0 {
        return src
    } else {
        buffer := make([]byte, len(dst) + len(src))
        copy(buffer, dst)
        copy(buffer[len(dst):], src)
        return buffer
    }
}

// 便利数组查找字符串索引位置，如果不存在则返回数组长度
func StringSearch (a []string, s string) int {
    for i, v := range a {
        if s == v {
            return i
        }
    }
    return len(a)
}

// 判断字符串是否在数组中
func StringInArray (a []string, s string) bool {
    return StringSearch(a, s) != len(a)
}



