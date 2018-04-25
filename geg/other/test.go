package main

import (
    "fmt"
)

var arabicNumber2RomanSymbolMap = map[int]string{
       1 : "I",
       5 : "V",
      10 : "X",
      50 : "L",
     100 : "C",
     500 : "D",
    1000 : "M",
}

var romanSymbol2ArabicNumberMap = map[int]string{
    4000 : "MMMCMC",
    1000 : "M",
     900 : "CM",
     500 : "D",
     400 : "CD",
     100 : "C",
      90 : "XC",
      50 : "L",
      40 : "XL",
      10 : "X",
       9 : "IX",
       5 : "V",
       4 : "IV",
       1 : "I",
}

var array = []int{4000, 1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}


// Arabic numbers to Roman symbols
func romanSymbols2ArabicNumber(s string) int {

}

// Arabic numbers to Roman symbols
func arabicNumberToRomanSymbols(i int) string {
    r := ""
    for i > 0 {
        for _, v := range array {
            if i >= v {
                r += romans[v]
                i -= v
                break
            }
        }
    }
    return r
}

func main() {
    fmt.Println(intToRomans(1944))
    //for i := 0; i < 10; i++ {
    //    fmt.Println(intToRomans(i))
    //}

}