package main

import (
    "github.com/olekukonko/tablewriter"
    "os"
)


func main() {
    data := [][]string{
        []string{"1/1/2014", "Domain name", "2233", "$10.98"},
        []string{"1/1/2014", "January Hosting", "2233", "$54.95"},
        []string{"1/4/2014", "February Hosting", "2233", "$51.00"},
        []string{"1/4/2014", "February Extra Bandwidth1111111111111111111111111111111111", "2233", "$30.00"},
    }

    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Date", "Description", "CV2", "Amount"})
    table.SetFooter([]string{"", "", "Total", "$146.93"}) // Add Footer
    table.SetCenterSeparator("|")
    table.SetAutoMergeCells(true)
    table.SetRowLine(true)
    table.SetBorder(false)
    table.AppendBulk(data)
    table.Render()
}
