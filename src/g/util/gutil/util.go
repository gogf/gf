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



