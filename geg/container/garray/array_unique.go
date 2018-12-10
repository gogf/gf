package main

import "fmt"

func main() {
    array := []uint{1, 1, 2, 3, 3, 4, 4, 5, 5, 6, 6}
    for i := 0; i < len(array) - 1; i++ {
        for j := i + 1; j < len(array); j++ {
            if array[i] == array[j] {
                array = append(array[ : j], array[j + 1 : ]...)
            }
        }
    }
    fmt.Println(array)


}