package main

import (
    "fmt"
    //"strconv"
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


type jsonParser struct {
    content []rune
    tokens  [][]rune
    size    int
}

func (p *jsonParser)getNextToken(token rune, from int) int {
    for i := from + 1; i < p.size; i++ {
        if p.content[i] == token {
            if i > 0 && p.content[i - 1] != GJSON_CHAR_QUOTATION {
                return i
            }
        } else {
            switch p.content[i] {
                case GJSON_CHAR_DOUBLE_QUOTE_MARK:
                    r := p.getNextToken(GJSON_CHAR_DOUBLE_QUOTE_MARK, i)
                    if r > 0 {
                        i = r
                    }
            }
        }
    }
    return 0
}

func (p *jsonParser)parseTokens()  {
    for i := 1; i < p.size; i++ {
        switch p.content[i] {
            case GJSON_CHAR_DOUBLE_QUOTE_MARK:
                r := p.getNextToken(GJSON_CHAR_DOUBLE_QUOTE_MARK, i)
                if r > 0 {
                    fmt.Println(string(p.content[i:r + 1]))
                    p.tokens = append(p.tokens, p.content[i:r + 1])
                    i = r + 1
                }

            case GJSON_CHAR_COLON:
                fallthrough
            case GJSON_CHAR_COMMA:
            case GJSON_CHAR_BRACE_LEFT:
            case GJSON_CHAR_BRACE_RIGHT:
            case GJSON_CHAR_BRACKET_LEFT:
            case GJSON_CHAR_BRACKET_RIGHT:
                fmt.Println(string(p.content[i:i]))
                p.tokens = append(p.tokens, p.content[i:i])
        }
    }
    // fmt.Println(p.tokens)
}

func (p *jsonParser)parse()  {
    for i := 1; i < p.size; i++ {
        switch p.content[i] {
        case GJSON_CHAR_DOUBLE_QUOTE_MARK:
            r := p.getNextToken(GJSON_CHAR_DOUBLE_QUOTE_MARK, i)
            if r > 0 {
                fmt.Println(string(p.content[i:r + 1]))
                i = r + 1
            }

        case GJSON_CHAR_BRACE_LEFT:
            r := p.getNextToken(GJSON_CHAR_BRACE_RIGHT, i)
            if r > 0 {
                fmt.Println(string(p.content[i:r + 1]))
                i = r + 1
            }
        }
    }

}

func jsonDecode(j *string)  {
    json   := []rune(*j)
    size   := len(json)
    parser := &jsonParser{
        content     : json,
        size        : size,
    }
    parser.parseTokens()

}

func main() {
    json := `{"name":"中国","age":31,"list":[["a","b","c"],["d","e","f"]],"item":{"title":"make\"he moon","name":"make'he moon","content":"'[}]{[}he moon"}}`
    jsonDecode(&json)
}