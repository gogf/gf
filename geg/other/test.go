package main
import (
    "fmt"
    "os"
)

func Add(path string, name ... string) {
    fmt.Println(name)
}

func main() {
    for _, e := range os.Environ() {

        fmt.Println(e)

    }

}