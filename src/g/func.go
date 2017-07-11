package g

// 三元运算符
func TriIf(c bool, t, f interface{}) interface{} {
    if c {
        return t
    }
    return f
}