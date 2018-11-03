package main

import (
    "html/template"
    "log"
    "os"
)

type Person string

func (p Person) Label() string {
    return "This is " + string(p)
}

func main() {
    tmpl, err := template.New("").Parse(`{{sum 1 2 3}}`)
    if err != nil {
        log.Fatalf("Parse: %v", err)
    }
    tmpl.Execute(os.Stdout, nil)
}