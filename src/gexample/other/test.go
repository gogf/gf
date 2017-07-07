package main

import (
    "fmt"
)

const (
    GJSON_CHAR_BRACE_LEFT        = rune('{')
    GJSON_CHAR_BRACE_RIGHT       = rune('}')
    GJSON_CHAR_BRACKET_LEFT      = rune('[')
    GJSON_CHAR_BRACKET_RIGHT     = rune(']')
    GJSON_CHAR_QUOTATION         = rune('\\')
    GJSON_CHAR_COMMA             = rune(',')
    GJSON_CHAR_COLON             = rune(':')
    GJSON_CHAR_DOUBLE_QUOTE_MARK = rune('"')
)

const (
    GJSON_TOKEN_BRACE_LEFT        = rune('{')
    GJSON_TOKEN_BRACE_RIGHT       = rune('}')
    GJSON_TOKEN_BRACKET_LEFT      = rune('[')
    GJSON_TOKEN_BRACKET_RIGHT     = rune(']')
    GJSON_TOKEN_COMMA             = rune(',')
    GJSON_TOKEN_COLON             = rune(':')
    GJSON_TOKEN_STRING            = rune('"')
    GJSON_TOKEN_NUMBER            = rune('0')
)

// JSON数据对象
type GJson struct {
    m      map[string]interface{}
    a      []interface{}
    level  int
    next  *GJson
}

type jsonToken struct {
    token     []rune
    tokenType rune
}

type jsonParser struct {
    content []rune
    tokens  []jsonToken
    root    GJson       // json根节点
    pointer *GJson      // 指向当前正在解析的json节点
}

// 遍历json字符串数组，并且判断转义
func (p *jsonParser)getNextChar(c rune, f int) int {
    for i := f + 1; i < len(p.content); i++ {
        if p.content[i] == c {
            if i > 0 && p.content[i - 1] != GJSON_CHAR_QUOTATION {
                return i
            }
        } else {
            switch p.content[i] {
                case GJSON_CHAR_DOUBLE_QUOTE_MARK:
                    r := p.getNextChar(GJSON_CHAR_DOUBLE_QUOTE_MARK, i)
                    if r > 0 {
                        i = r
                    }
            }
        }
    }
    return 0
}

// 判断字符是否为数字
func (p *jsonParser)isCharNumber(c rune) bool {
    if c >= rune('0') && c <= rune('9') {
        return true
    }
    return false
}

// 将json字符串解析为语义token
func (p *jsonParser)parseTokens() {
    for i := 0; i < len(p.content); i++ {
        if p.isCharNumber(p.content[i]) {
            j := i + 1
            for ; j < len(p.content); j++ {
                if !p.isCharNumber(p.content[j]) {
                    break;
                }
            }
            p.tokens = append(p.tokens, jsonToken{token: p.content[i:j], tokenType: GJSON_TOKEN_NUMBER})
            i = j - 1
        } else {
            switch p.content[i] {
                case GJSON_CHAR_DOUBLE_QUOTE_MARK:
                    r := p.getNextChar(GJSON_CHAR_DOUBLE_QUOTE_MARK, i)
                    if r > 0 {
                        p.tokens = append(p.tokens, jsonToken{token: p.content[i:r+1], tokenType: GJSON_TOKEN_STRING})
                        i = r
                    }
                case GJSON_CHAR_COLON:
                    p.tokens = append(p.tokens, jsonToken{token: p.content[i:i+1], tokenType: GJSON_TOKEN_COLON})
                case GJSON_CHAR_COMMA:
                    p.tokens = append(p.tokens, jsonToken{token: p.content[i:i+1], tokenType: GJSON_TOKEN_COMMA})
                case GJSON_CHAR_BRACE_LEFT:
                    p.tokens = append(p.tokens, jsonToken{token: p.content[i:i+1], tokenType: GJSON_TOKEN_BRACE_LEFT})
                case GJSON_CHAR_BRACE_RIGHT:
                    p.tokens = append(p.tokens, jsonToken{token: p.content[i:i+1], tokenType: GJSON_TOKEN_BRACE_RIGHT})
                case GJSON_CHAR_BRACKET_LEFT:
                    p.tokens = append(p.tokens, jsonToken{token: p.content[i:i+1], tokenType: GJSON_TOKEN_BRACKET_LEFT})
                case GJSON_CHAR_BRACKET_RIGHT:
                    p.tokens = append(p.tokens, jsonToken{token: p.content[i:i+1], tokenType: GJSON_TOKEN_BRACKET_RIGHT})
            }
        }

    }
}

func (p *jsonParser)getTokenBorderRightIndex(token rune, from int) int {
    switch token {
        case GJSON_TOKEN_BRACE_LEFT:
            leftCount := 0
            for i := from + 1; i < len(p.tokens); i++ {
                if p.tokens[i].tokenType == GJSON_TOKEN_BRACE_LEFT {
                    leftCount ++
                } else if p.tokens[i].tokenType == GJSON_TOKEN_BRACE_RIGHT {
                    if leftCount < 1 {
                        return i
                    } else {
                        leftCount--
                    }
                }
            }
        case GJSON_CHAR_BRACKET_LEFT:
            leftCount := 0
            for i := from + 1; i < len(p.tokens); i++ {
                if p.tokens[i].tokenType == GJSON_CHAR_BRACKET_LEFT {
                    leftCount ++
                } else if p.tokens[i].tokenType == GJSON_CHAR_BRACKET_RIGHT {
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

//func (p *jsonParser)getPrevToken(from int) jsonToken {
//    for i := from -1; i >= 0; i-- {
//        if p.tokens[i].tokenType != GJSON_CHAR_COMMA {
//            return p.tokens[i]
//        }
//    }
//    return nil
//}

func (p *jsonParser)parseTokenNodeToVar(left int, right int, level int) {
    //fmt.Println("================================")
    //for i := left; i <= right; i++ {
    //    fmt.Println(string(p.tokens[i].token))
    //}
    for i := left; i <= right; i++ {
        //fmt.Println(string(p.tokens[i].token))
        switch p.tokens[i].tokenType {
            case GJSON_TOKEN_BRACE_LEFT:
                fallthrough
            case GJSON_CHAR_BRACKET_LEFT:
                node       := newJsonNode()
                node.level  = level
                // 判断层级关系
                if i > 0 && p.tokens[i-1].tokenType == GJSON_CHAR_COLON {
                    oldptr := p.pointer
                    k      := string(p.tokens[i-2].token)
                    node   := newJsonNode()
                    p.pointer.m[k] = &node
                    p.pointer      = &node
                    j := p.getTokenBorderRightIndex(p.tokens[i].tokenType, i)
                    p.parseTokenNodeToVar(i + 1, j - 1, level)
                    i         = j
                    p.pointer = oldptr
                } else {
                    p.pointer.next = &node
                    p.pointer      = &node
                    j := p.getTokenBorderRightIndex(p.tokens[i].tokenType, i)
                    p.parseTokenNodeToVar(i + 1, j - 1, level + 1)
                    i = j
                }


            case GJSON_TOKEN_STRING:
                fallthrough
            case GJSON_TOKEN_NUMBER:
                if i > 0 && p.tokens[i-1].tokenType == GJSON_CHAR_COLON {
                    k := string(p.tokens[i-2].token)
                    v := string(p.tokens[i].token)
                    p.pointer.m[k] = v
                } else if p.tokens[i+1].tokenType != GJSON_CHAR_COLON {
                    p.pointer.a = append(p.pointer.a, p.tokens[i].token)
                }
        }
    }
}

// 打印出所有的token
func (p *jsonParser)printTokens() {
    for _, v := range p.tokens {
        fmt.Println(string(v.token))
    }
}

// 格式化打印根节点
func (p *jsonParser)printRootNode() {
    node := &p.root
    for {
        fmt.Println("==============")
        fmt.Println(node.level)
        fmt.Println(node.m)
        fmt.Println(node.a)
        for k, v := range node.m {
            fmt.Println(k)
            fmt.Println(v)
        }
        if node.next != nil {
            node = node.next
        } else {
            break;
        }
    }
}

// 创建一个json数据对象
func newJsonNode() GJson {
    return GJson {
        m: make(map[string]interface{}),
        a: make([]interface{}, 0),
        next: nil,
    }
}

func jsonDecode(j *string)  {
    p        := &jsonParser{content:[]rune(*j)}
    p.root    = newJsonNode()
    p.pointer = &p.root
    p.parseTokens()
    //p.printTokens()
    p.parseTokenNodeToVar(0, len(p.tokens) - 1, 1)
    p.printRootNode()
}

func main() {
    //json := `{"name":"中国","age":31,"list":[["a","b","c"],["d","e","f"]],"item":{"title":"make\"he moon","name":"make'he moon","content":"'[}]{[}he moon"}}`
    //json := `{"name"  :  "中国",  "age" : 31, "items":[1,2,3]}`
    //json := `[["a","b","c"],["d","e","f"]]`
    //json := `["a","b","c"]`
    //jsonDecode(&json)
    //fmt.Println()
    //fmt.Println()
    ////v := make(map[string]interface{})
    ////i := 31
    ////j := "john"
    ////v["age"]  = i
    ////v["name"] = make(map[string]interface{})
    ////t := v["name"]
    ////t.(map[string]interface{})["n"] = j
    ////
    ////fmt.Println(v)
    //var s struct{
    //    v interface{}
    //    p interface{}
    //}
    //v  := make(map[string]interface{})
    //s.v = v
    //s.p = &v
    //c  := (*s.p.(*map[string]interface{}))
    //c["name1"] = "john1"
    //
    //t          := make(map[string]interface{})
    //c["/"]      = t
    //s.p         = &t
    //t["name2"]  = "john2"
    //
    //c2         := (*s.p.(*map[string]interface{}))
    //c2["name3"] = "john3"
    //
    ////t2[2] = 100
    //fmt.Println(s.v)


    a := map[string]interface{} {
        "name" : "john",
        "list" : []interface{}{
            1,2,3, "fuck",
        },
        "item" : map[string]string {
            "n1" : "v1",
            "n2" : "v2",
            "n3" : "v3",
        },
    }
    fmt.Println(a["list"][0])

    //
    //var a = []int{1,2,3}
    //var b = []int{4,5,6, 7,8}
    //cc := make([]int, len(a) + 12)
    //a = cc
    //copy(a, b)
    //fmt.Println(a)
    //fmt.Println(b)
}