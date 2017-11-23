package internal

import (
    "fmt"
    "errors"
    "strings"
    "strconv"
)

// 这是一个使用go进行json语法解析的解析器，效率没有官方的json解析高，仅作学习参考

const (
    gJSON_CHAR_BRACE_LEFT         = rune('{')
    gJSON_CHAR_BRACE_RIGHT        = rune('}')
    gJSON_CHAR_BRACKET_LEFT       = rune('[')
    gJSON_CHAR_BRACKET_RIGHT      = rune(']')
    gJSON_CHAR_QUOTATION          = rune('\\')
    gJSON_CHAR_COMMA              = rune(',')
    gJSON_CHAR_COLON              = rune(':')
    gJSON_CHAR_DOUBLE_QUOTE_MARK  = rune('"')
)

const (
    gJSON_TOKEN_BRACE_LEFT        = rune('{')
    gJSON_TOKEN_BRACE_RIGHT       = rune('}')
    gJSON_TOKEN_BRACKET_LEFT      = rune('[')
    gJSON_TOKEN_BRACKET_RIGHT     = rune(']')
    gJSON_TOKEN_COMMA             = rune(',')
    gJSON_TOKEN_COLON             = rune(':')
    gJSON_TOKEN_STRING            = rune('"')
    gJSON_TOKEN_NUMBER            = rune('0')
)

// json关联数组(哈希表)
type JsonMap   map[string]interface{}
// json索引数组(普通数组，从0开始索引)
type JsonArray []interface{}

// JSON数据对象
type gJsonNode struct {
    m JsonMap
    a JsonArray
}

// JSON语义token
type gJsonToken struct {
    token      []rune // token字符串
    tokenType  rune   // token类型
    tokenindex int    // token在原始字符串中的索引位置
}

// JSON解析结构对象
type gJsonParser struct {
    content []rune       // 需要解析json字符串(通过string转换为[]rune)
    tokens  []gJsonToken // 存放解析content后的json token数组
    root    *gJsonNode   // json根节点
    pointer *gJsonNode   // 指向当前正在解析的json节点
}

// 解析json字符串
func Decode(j *string) (*gJsonParser, error) {
    p   := &gJsonParser{content:[]rune(*j)}
    err := p.parse()
    if err == nil {
        return p, err
    } else {
        return nil, err
    }
}

// 判断所给字符串是否为数字
func isNumeric(s string) bool  {
    for i :=0; i < len(s); i++ {
        if s[i] < byte('0') || s[i] > byte('9') {
            return false
        }
    }
    return true
}

// 获得一个键值对关联数组/哈希表，方便操作，不需要自己做类型转换
// 注意，如果获取的值不存在，或者类型与json类型不匹配，那么将会返回nil
func (p *gJsonParser) GetMap(pattern string) JsonMap {
    result := p.Get(pattern)
    if result != nil {
        if r, ok := result.(JsonMap); ok {
            return r
        }
    }
    return nil
}

// 获得一个数组[]interface{}，方便操作，不需要自己做类型转换
// 注意，如果获取的值不存在，或者类型与json类型不匹配，那么将会返回nil
func (p *gJsonParser) GetArray(pattern string) JsonArray {
    result := p.Get(pattern)
    if result != nil {
        if r, ok := result.(JsonArray); ok {
            return r
        }
    }
    return nil
}


// 根据约定字符串方式访问json解析数据，参数形如： "items.name.first", "list.0"
// 返回的结果类型的interface{}，因此需要自己做类型转换
// 如果找不到对应节点的数据，返回nil
func (p *gJsonParser) Get(pattern string) interface{} {
    var result interface{}
    pointer  := p.root
    array    := strings.Split(pattern, ".")
    length   := len(array)
    for i:= 0; i < length; i++ {
        // 优先判断数组
        if isNumeric(array[i]) {
            n, err := strconv.Atoi(array[i])
            if err == nil && len(pointer.a) > n {
                if i == length - 1 {
                    result = pointer.a[n]
                    break;
                } else {
                    if p, ok := pointer.a[n].(*gJsonNode); ok {
                        pointer = p
                        continue
                    }
                }
            }
        }
        // 其次判断哈希表，如果一个键在数组及map中均不存在，直接返回nil
        if v, ok := pointer.m[array[i]]; ok {
            if i == length - 1 {
                result = v
            } else {
                if p, ok := v.(*gJsonNode); ok {
                    pointer = p
                    continue
                }
            }
        } else {
            return nil
        }
    }
    // 处理结果，如果是gJsonNode类型，那么需要做转换
    if r, ok := result.(*gJsonNode); ok {
        if len(r.m) < 1 {
            return r.a
        } else {
            return r.m
        }
    }
    return result
}

// 遍历json字符串数组，并且判断转义
func (p *gJsonParser) getNextChar(c rune, f int) int {
    for i := f + 1; i < len(p.content); i++ {
        if p.content[i] == c {
            if i > 0 && p.content[i - 1] != gJSON_CHAR_QUOTATION {
                return i
            }
        } else {
            switch p.content[i] {
                case gJSON_CHAR_DOUBLE_QUOTE_MARK:
                    r := p.getNextChar(gJSON_CHAR_DOUBLE_QUOTE_MARK, i)
                    if r > 0 {
                        i = r
                    }
            }
        }
    }
    return 0
}

// 判断字符是否为数字
func (p *gJsonParser) isCharNumber(c rune) bool {
    if c >= rune('0') && c <= rune('9') {
        return true
    }
    return false
}

// 按照json语法对保存的字符串进行解析
func (p *gJsonParser) parse() error {
    // 首先将字符串解析成token进行保存
    for i := 0; i < len(p.content); i++ {
        if p.isCharNumber(p.content[i]) {
            j := i + 1
            for ; j < len(p.content); j++ {
                if !p.isCharNumber(p.content[j]) {
                    break;
                }
            }
            p.tokens = append(p.tokens, gJsonToken {
                token:      p.content[i:j],
                tokenType:  gJSON_TOKEN_NUMBER,
                tokenindex: i,
            })
            i = j - 1
        } else {
            switch p.content[i] {
                case gJSON_CHAR_DOUBLE_QUOTE_MARK:
                    r := p.getNextChar(gJSON_CHAR_DOUBLE_QUOTE_MARK, i)
                    if r > 0 {
                        // 注意这里需要去掉字符串两边的双引号
                        p.tokens = append(p.tokens, gJsonToken {
                            token:      p.content[i+1:r],
                            tokenType:  gJSON_TOKEN_STRING,
                            tokenindex: i,
                        })
                        i = r
                    }
                case gJSON_CHAR_COLON:
                    p.tokens = append(p.tokens, gJsonToken{token: p.content[i:i+1], tokenType: gJSON_TOKEN_COLON,         tokenindex: i})
                case gJSON_CHAR_COMMA:
                    p.tokens = append(p.tokens, gJsonToken{token: p.content[i:i+1], tokenType: gJSON_TOKEN_COMMA,         tokenindex: i})
                case gJSON_CHAR_BRACE_LEFT:
                    p.tokens = append(p.tokens, gJsonToken{token: p.content[i:i+1], tokenType: gJSON_TOKEN_BRACE_LEFT,    tokenindex: i})
                case gJSON_CHAR_BRACE_RIGHT:
                    p.tokens = append(p.tokens, gJsonToken{token: p.content[i:i+1], tokenType: gJSON_TOKEN_BRACE_RIGHT,   tokenindex: i})
                case gJSON_CHAR_BRACKET_LEFT:
                    p.tokens = append(p.tokens, gJsonToken{token: p.content[i:i+1], tokenType: gJSON_TOKEN_BRACKET_LEFT,  tokenindex: i})
                case gJSON_CHAR_BRACKET_RIGHT:
                    p.tokens = append(p.tokens, gJsonToken{token: p.content[i:i+1], tokenType: gJSON_TOKEN_BRACKET_RIGHT, tokenindex: i})

                default:
                    c := string(p.content[i])
                    if c != " " && c != "\r" && c != "\n" && c != "\t" {
                        return errors.New(fmt.Sprintf("json parse error: invalid char '%s' at index %d", c, i))
                    }
            }
        }
    }
    // 最后对解析后的token转换为go变量
    return p.parseTokenNodeToVar(0, len(p.tokens) - 1)
}

// 获取json范围字符包含范围最右侧的索引位置
func (p *gJsonParser)getTokenBorderRightIndex(token rune, from int) int {
    switch token {
        case gJSON_TOKEN_BRACE_LEFT:
            leftCount := 0
            for i := from + 1; i < len(p.tokens); i++ {
                if p.tokens[i].tokenType == gJSON_TOKEN_BRACE_LEFT {
                    leftCount ++
                } else if p.tokens[i].tokenType == gJSON_TOKEN_BRACE_RIGHT {
                    if leftCount < 1 {
                        return i
                    } else {
                        leftCount--
                    }
                }
            }
        case gJSON_CHAR_BRACKET_LEFT:
            leftCount := 0
            for i := from + 1; i < len(p.tokens); i++ {
                if p.tokens[i].tokenType == gJSON_CHAR_BRACKET_LEFT {
                    leftCount ++
                } else if p.tokens[i].tokenType == gJSON_CHAR_BRACKET_RIGHT {
                    if leftCount < 1 {
                        return i
                    } else {
                        leftCount--
                    }
                }
            }
    }
    return 0
}

// 将解析过后的json token转换为go变量
func (p *gJsonParser) parseTokenNodeToVar(left int, right int) error {
    //fmt.Println("================================")
    //for i := left; i <= right; i++ {
    //    fmt.Println(string(p.tokens[i].token))
    //}
    for i := left; i <= right; i++ {
        //fmt.Println(string(p.tokens[i].token))
        switch p.tokens[i].tokenType {
            case gJSON_TOKEN_BRACE_LEFT:
                fallthrough
            case gJSON_TOKEN_BRACKET_LEFT:
                node := newJsonNode()
                // 判断根节点
                if p.root == nil {
                    p.root    = node
                    p.pointer = node
                }
                // 判断层级关系
                borderRight := p.getTokenBorderRightIndex(p.tokens[i].tokenType, i)
                if borderRight < 1 {
                    return errors.New(fmt.Sprintf("json parse error: unclosed tag '%s' at index %d", string(p.tokens[i].token), p.tokens[i].tokenindex))
                }
                if i > 1 && (
                    p.tokens[i-1].tokenType == gJSON_TOKEN_COLON &&
                    p.tokens[i-2].tokenType == gJSON_TOKEN_STRING) {
                    // json赋值操作
                    oldptr        := p.pointer
                    k             := string(p.tokens[i-2].token)
                    p.pointer.m[k] = node
                    p.pointer      = node
                    err           := p.parseTokenNodeToVar(i + 1, borderRight - 1)
                    if err != nil {
                        return err
                    } else {
                        i         = borderRight
                        p.pointer = oldptr
                    }

                } else if i > 0 && (
                    p.tokens[i-1].tokenType == gJSON_TOKEN_COMMA ||
                    p.tokens[i-1].tokenType == gJSON_TOKEN_BRACE_LEFT ||
                    p.tokens[i-1].tokenType == gJSON_TOKEN_BRACKET_LEFT) {
                    // json数组操作
                    oldptr     := p.pointer
                    p.pointer.a = append(p.pointer.a, node)
                    p.pointer   = node
                    err        := p.parseTokenNodeToVar(i + 1, borderRight - 1)
                    if err != nil {
                        return err
                    } else {
                        i         = borderRight
                        p.pointer = oldptr
                    }
                } else {
                    // json层级关系
                    p.pointer = node
                    err      := p.parseTokenNodeToVar(i + 1, borderRight - 1)
                    if err != nil {
                        return err
                    } else {
                        i = borderRight
                    }
                }

            case gJSON_TOKEN_STRING:
                fallthrough
            case gJSON_TOKEN_NUMBER:
                if i > 0 && p.tokens[i-1].tokenType == gJSON_TOKEN_COLON {
                    k := string(p.tokens[i-2].token)
                    v := string(p.tokens[i].token)
                    p.pointer.m[k] = v
                } else if p.tokens[i+1].tokenType != gJSON_TOKEN_COLON {
                    p.pointer.a = append(p.pointer.a, string(p.tokens[i].token))
                }

            case gJSON_TOKEN_COLON:
                if i < 1 || (p.tokens[i-1].tokenType != gJSON_TOKEN_STRING) {
                    return errors.New(fmt.Sprintf("json parse error: invalid charactar '%s' at index %d", string(p.tokens[i].token), p.tokens[i].tokenindex))
                }

            case gJSON_TOKEN_COMMA:
                if (p.tokens[i+1].tokenType != gJSON_TOKEN_STRING &&
                    p.tokens[i+1].tokenType != gJSON_TOKEN_NUMBER &&
                    p.tokens[i+1].tokenType != gJSON_TOKEN_BRACE_LEFT &&
                    p.tokens[i+1].tokenType != gJSON_TOKEN_BRACKET_LEFT) ||
                    (i < 1 || (
                    p.tokens[i-1].tokenType != gJSON_TOKEN_STRING &&
                    p.tokens[i-1].tokenType != gJSON_TOKEN_NUMBER &&
                    p.tokens[i-1].tokenType != gJSON_TOKEN_BRACE_RIGHT &&
                    p.tokens[i-1].tokenType != gJSON_TOKEN_BRACKET_RIGHT)) {
                    return errors.New(fmt.Sprintf("json parse error: invalid charactar '%s' at index %d", string(p.tokens[i].token), p.tokens[i].tokenindex))
                }
        }
    }
    return nil
}

// 打印出所有的token(测试用)
func (p *gJsonParser)printTokens() {
    for _, v := range p.tokens {
        fmt.Println(string(v.token))
    }
}

// 格式化打印根节点
func (p *gJsonParser)Print() {
    if len(p.root.m) > 0 {
        fmt.Println("{")
    } else {
        fmt.Println("[")
    }
    p.printNode(p.pointer, "\t")
    if len(p.root.m) > 0 {
        fmt.Println("}")
    } else {
        fmt.Println("]")
    }
}

// 格式化打印根节点
func (p *gJsonParser)printNode(n *gJsonNode, indent string) {
    if len(n.m) > 0 {
        for k, v := range n.m {
            if t, ok := v.(*gJsonNode); ok {
                if len(t.m) > 0 {
                    fmt.Printf("%v%v\t: {\n", indent, k)
                    p.printNode(t, indent + "\t")
                    fmt.Printf("%v}\n", indent)
                } else {
                    fmt.Printf("%v%v\t: [\n", indent, k)
                    p.printNode(t, indent + "\t")
                    fmt.Printf("%v}\n", indent)
                }
            } else {
                fmt.Printf("%v%v\t: %v\n", indent, k, v)
            }
        }
    }
    if len(n.a) > 0 {
        for k, v := range n.a {
            if t, ok := v.(*gJsonNode); ok {
                if len(t.m) > 0 {
                    fmt.Printf("%v%v\t: {\n", indent, k)
                    p.printNode(t, indent + "\t")
                    fmt.Printf("%v}\n", indent)
                } else {
                    fmt.Printf("%v%v\t: [\n", indent, k)
                    p.printNode(t, indent + "\t")
                    fmt.Printf("%v}\n", indent)
                }
            } else {
                fmt.Printf("%v%v : %v\n", indent, k, v)
            }
        }
    }
}

// 创建一个json数据对象
func newJsonNode() *gJsonNode {
    return &gJsonNode {
        m: make(map[string]interface{}),
        a: make([]interface{}, 0),
    }
}

