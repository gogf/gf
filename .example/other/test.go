package main

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func main() {
	data := [][]string{
		[]string{"1/1/2014", "Domain name", "2233", "$10.98"},
		[]string{"1/1/2014", "January Hosting", "2233", "$54.95"},
		[]string{"1/4/2014", "February Hosting", "2233", "$51.00"},
		[]string{"1/4/2014", "February Extra Bandwidth", "2233", "$30.00"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetColumnSeparator("")
	table.AppendBulk(data)
	table.Render()
}
