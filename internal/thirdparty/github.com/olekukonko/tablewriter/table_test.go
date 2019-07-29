// Copyright 2014 Oleku Konko All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// This module is a Table Writer  API for the Go Programming Language.
// The protocols were written in pure Go and works on windows and unix systems

package tablewriter

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func checkEqual(t *testing.T, got, want interface{}, msgs ...interface{}) {
	if !reflect.DeepEqual(got, want) {
		buf := bytes.Buffer{}
		buf.WriteString("got:\n[%v]\nwant:\n[%v]\n")
		for _, v := range msgs {
			buf.WriteString(v.(string))
		}
		t.Errorf(buf.String(), got, want)
	}
}

func ExampleShort() {
	data := [][]string{
		{"A", "The Good", "500"},
		{"B", "The Very very Bad Man", "288"},
		{"C", "The Ugly", "120"},
		{"D", "The Gopher", "800"},
	}

	table := NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Sign", "Rating"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()

	// Output: +------+-----------------------+--------+
	// | NAME |         SIGN          | RATING |
	// +------+-----------------------+--------+
	// | A    | The Good              |    500 |
	// | B    | The Very very Bad Man |    288 |
	// | C    | The Ugly              |    120 |
	// | D    | The Gopher            |    800 |
	// +------+-----------------------+--------+
}

func ExampleLong() {
	data := [][]string{
		{"Learn East has computers with adapted keyboards with enlarged print etc", "  Some Data  ", " Another Data"},
		{"Instead of lining up the letters all ", "the way across, he splits the keyboard in two", "Like most ergonomic keyboards", "See Data"},
	}

	table := NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Sign", "Rating"})
	table.SetCenterSeparator("*")
	table.SetRowSeparator("=")

	for _, v := range data {
		table.Append(v)
	}
	table.Render()

	// Output: *================================*================================*===============================*==========*
	// |              NAME              |              SIGN              |            RATING             |          |
	// *================================*================================*===============================*==========*
	// | Learn East has computers       |   Some Data                    |  Another Data                 |
	// | with adapted keyboards with    |                                |                               |
	// | enlarged print etc             |                                |                               |
	// | Instead of lining up the       | the way across, he splits the  | Like most ergonomic keyboards | See Data |
	// | letters all                    | keyboard in two                |                               |          |
	// *================================*================================*===============================*==========*
}

func ExampleCSV() {
	table, _ := NewCSV(os.Stdout, "testdata/test.csv", true)
	table.SetCenterSeparator("*")
	table.SetRowSeparator("=")

	table.Render()

	// Output: *============*===========*=========*
	// | FIRST NAME | LAST NAME |   SSN   |
	// *============*===========*=========*
	// | John       | Barry     |  123456 |
	// | Kathy      | Smith     |  687987 |
	// | Bob        | McCornick | 3979870 |
	// *============*===========*=========*
}

// TestNumLines to test the numbers of lines
func TestNumLines(t *testing.T) {
	data := [][]string{
		{"A", "The Good", "500"},
		{"B", "The Very very Bad Man", "288"},
		{"C", "The Ugly", "120"},
		{"D", "The Gopher", "800"},
	}

	buf := &bytes.Buffer{}
	table := NewWriter(buf)
	table.SetHeader([]string{"Name", "Sign", "Rating"})

	for i, v := range data {
		table.Append(v)
		checkEqual(t, table.NumLines(), i+1, "Number of lines failed")
	}

	checkEqual(t, table.NumLines(), len(data), "Number of lines failed")
}

func TestCSVInfo(t *testing.T) {
	buf := &bytes.Buffer{}
	table, err := NewCSV(buf, "testdata/test_info.csv", true)
	if err != nil {
		t.Error(err)
		return
	}
	table.SetAlignment(ALIGN_LEFT)
	table.SetBorder(false)
	table.Render()

	got := buf.String()
	want := `   FIELD   |     TYPE     | NULL | KEY | DEFAULT |     EXTRA       
+----------+--------------+------+-----+---------+----------------+
  user_id  | smallint(5)  | NO   | PRI | NULL    | auto_increment  
  username | varchar(10)  | NO   |     | NULL    |                 
  password | varchar(100) | NO   |     | NULL    |                 
`
	checkEqual(t, got, want, "CSV info failed")
}

func TestCSVSeparator(t *testing.T) {
	buf := &bytes.Buffer{}
	table, err := NewCSV(buf, "testdata/test.csv", true)
	if err != nil {
		t.Error(err)
		return
	}
	table.SetRowLine(true)
	table.SetCenterSeparator("+")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")
	table.SetAlignment(ALIGN_LEFT)
	table.Render()

	want := `+------------+-----------+---------+
| FIRST NAME | LAST NAME |   SSN   |
+------------+-----------+---------+
| John       | Barry     | 123456  |
+------------+-----------+---------+
| Kathy      | Smith     | 687987  |
+------------+-----------+---------+
| Bob        | McCornick | 3979870 |
+------------+-----------+---------+
`

	checkEqual(t, buf.String(), want, "CSV info failed")
}

func TestNoBorder(t *testing.T) {
	data := [][]string{
		{"1/1/2014", "Domain name", "2233", "$10.98"},
		{"1/1/2014", "January Hosting", "2233", "$54.95"},
		{"", "    (empty)\n    (empty)", "", ""},
		{"1/4/2014", "February Hosting", "2233", "$51.00"},
		{"1/4/2014", "February Extra Bandwidth", "2233", "$30.00"},
		{"1/4/2014", "    (Discount)", "2233", "-$1.00"},
	}

	var buf bytes.Buffer
	table := NewWriter(&buf)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Date", "Description", "CV2", "Amount"})
	table.SetFooter([]string{"", "", "Total", "$145.93"}) // Add Footer
	table.SetBorder(false)                                // Set Border to false
	table.AppendBulk(data)                                // Add Bulk Data
	table.Render()

	want := `    DATE   |       DESCRIPTION        |  CV2  | AMOUNT   
+----------+--------------------------+-------+---------+
  1/1/2014 | Domain name              |  2233 | $10.98   
  1/1/2014 | January Hosting          |  2233 | $54.95   
           |     (empty)              |       |          
           |     (empty)              |       |          
  1/4/2014 | February Hosting         |  2233 | $51.00   
  1/4/2014 | February Extra Bandwidth |  2233 | $30.00   
  1/4/2014 |     (Discount)           |  2233 | -$1.00   
+----------+--------------------------+-------+---------+
                                        TOTAL | $145.93  
                                      +-------+---------+
`

	checkEqual(t, buf.String(), want, "border table rendering failed")
}

func TestWithBorder(t *testing.T) {
	data := [][]string{
		{"1/1/2014", "Domain name", "2233", "$10.98"},
		{"1/1/2014", "January Hosting", "2233", "$54.95"},
		{"", "    (empty)\n    (empty)", "", ""},
		{"1/4/2014", "February Hosting", "2233", "$51.00"},
		{"1/4/2014", "February Extra Bandwidth", "2233", "$30.00"},
		{"1/4/2014", "    (Discount)", "2233", "-$1.00"},
	}

	var buf bytes.Buffer
	table := NewWriter(&buf)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Date", "Description", "CV2", "Amount"})
	table.SetFooter([]string{"", "", "Total", "$145.93"}) // Add Footer
	table.AppendBulk(data)                                // Add Bulk Data
	table.Render()

	want := `+----------+--------------------------+-------+---------+
|   DATE   |       DESCRIPTION        |  CV2  | AMOUNT  |
+----------+--------------------------+-------+---------+
| 1/1/2014 | Domain name              |  2233 | $10.98  |
| 1/1/2014 | January Hosting          |  2233 | $54.95  |
|          |     (empty)              |       |         |
|          |     (empty)              |       |         |
| 1/4/2014 | February Hosting         |  2233 | $51.00  |
| 1/4/2014 | February Extra Bandwidth |  2233 | $30.00  |
| 1/4/2014 |     (Discount)           |  2233 | -$1.00  |
+----------+--------------------------+-------+---------+
|                                       TOTAL | $145.93 |
+----------+--------------------------+-------+---------+
`

	checkEqual(t, buf.String(), want, "border table rendering failed")
}

func TestPrintingInMarkdown(t *testing.T) {
	data := [][]string{
		{"1/1/2014", "Domain name", "2233", "$10.98"},
		{"1/1/2014", "January Hosting", "2233", "$54.95"},
		{"1/4/2014", "February Hosting", "2233", "$51.00"},
		{"1/4/2014", "February Extra Bandwidth", "2233", "$30.00"},
	}

	var buf bytes.Buffer
	table := NewWriter(&buf)
	table.SetHeader([]string{"Date", "Description", "CV2", "Amount"})
	table.AppendBulk(data) // Add Bulk Data
	table.SetBorders(Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.Render()

	want := `|   DATE   |       DESCRIPTION        | CV2  | AMOUNT |
|----------|--------------------------|------|--------|
| 1/1/2014 | Domain name              | 2233 | $10.98 |
| 1/1/2014 | January Hosting          | 2233 | $54.95 |
| 1/4/2014 | February Hosting         | 2233 | $51.00 |
| 1/4/2014 | February Extra Bandwidth | 2233 | $30.00 |
`
	checkEqual(t, buf.String(), want, "border table rendering failed")
}

func TestPrintHeading(t *testing.T) {
	var buf bytes.Buffer
	table := NewWriter(&buf)
	table.SetHeader([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c"})
	table.printHeading()
	want := `| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | A | B | C |
+---+---+---+---+---+---+---+---+---+---+---+---+
`
	checkEqual(t, buf.String(), want, "header rendering failed")
}

func TestPrintHeadingWithoutAutoFormat(t *testing.T) {
	var buf bytes.Buffer
	table := NewWriter(&buf)
	table.SetHeader([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c"})
	table.SetAutoFormatHeaders(false)
	table.printHeading()
	want := `| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | a | b | c |
+---+---+---+---+---+---+---+---+---+---+---+---+
`
	checkEqual(t, buf.String(), want, "header rendering failed")
}

func TestPrintFooter(t *testing.T) {
	var buf bytes.Buffer
	table := NewWriter(&buf)
	table.SetHeader([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c"})
	table.SetFooter([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c"})
	table.printFooter()
	want := `| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | A | B | C |
+---+---+---+---+---+---+---+---+---+---+---+---+
`
	checkEqual(t, buf.String(), want, "footer rendering failed")
}

func TestPrintFooterWithoutAutoFormat(t *testing.T) {
	var buf bytes.Buffer
	table := NewWriter(&buf)
	table.SetAutoFormatHeaders(false)
	table.SetHeader([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c"})
	table.SetFooter([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c"})
	table.printFooter()
	want := `| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | a | b | c |
+---+---+---+---+---+---+---+---+---+---+---+---+
`
	checkEqual(t, buf.String(), want, "footer rendering failed")
}

func TestPrintShortCaption(t *testing.T) {
	var buf bytes.Buffer
	data := [][]string{
		{"A", "The Good", "500"},
		{"B", "The Very very Bad Man", "288"},
		{"C", "The Ugly", "120"},
		{"D", "The Gopher", "800"},
	}

	table := NewWriter(&buf)
	table.SetHeader([]string{"Name", "Sign", "Rating"})
	table.SetCaption(true, "Short caption.")

	for _, v := range data {
		table.Append(v)
	}
	table.Render()

	want := `+------+-----------------------+--------+
| NAME |         SIGN          | RATING |
+------+-----------------------+--------+
| A    | The Good              |    500 |
| B    | The Very very Bad Man |    288 |
| C    | The Ugly              |    120 |
| D    | The Gopher            |    800 |
+------+-----------------------+--------+
Short caption.
`
	checkEqual(t, buf.String(), want, "long caption for short example rendering failed")
}

func TestPrintLongCaptionWithShortExample(t *testing.T) {
	var buf bytes.Buffer
	data := [][]string{
		{"A", "The Good", "500"},
		{"B", "The Very very Bad Man", "288"},
		{"C", "The Ugly", "120"},
		{"D", "The Gopher", "800"},
	}

	table := NewWriter(&buf)
	table.SetHeader([]string{"Name", "Sign", "Rating"})
	table.SetCaption(true, "This is a very long caption. The text should wrap. If not, we have a problem that needs to be solved.")

	for _, v := range data {
		table.Append(v)
	}
	table.Render()

	want := `+------+-----------------------+--------+
| NAME |         SIGN          | RATING |
+------+-----------------------+--------+
| A    | The Good              |    500 |
| B    | The Very very Bad Man |    288 |
| C    | The Ugly              |    120 |
| D    | The Gopher            |    800 |
+------+-----------------------+--------+
This is a very long caption. The text
should wrap. If not, we have a problem
that needs to be solved.
`
	checkEqual(t, buf.String(), want, "long caption for short example rendering failed")
}

func TestPrintCaptionWithFooter(t *testing.T) {
	data := [][]string{
		{"1/1/2014", "Domain name", "2233", "$10.98"},
		{"1/1/2014", "January Hosting", "2233", "$54.95"},
		{"1/4/2014", "February Hosting", "2233", "$51.00"},
		{"1/4/2014", "February Extra Bandwidth", "2233", "$30.00"},
	}

	var buf bytes.Buffer
	table := NewWriter(&buf)
	table.SetHeader([]string{"Date", "Description", "CV2", "Amount"})
	table.SetFooter([]string{"", "", "Total", "$146.93"})                                                  // Add Footer
	table.SetCaption(true, "This is a very long caption. The text should wrap to the width of the table.") // Add caption
	table.SetBorder(false)                                                                                 // Set Border to false
	table.AppendBulk(data)                                                                                 // Add Bulk Data
	table.Render()

	want := `    DATE   |       DESCRIPTION        |  CV2  | AMOUNT   
+----------+--------------------------+-------+---------+
  1/1/2014 | Domain name              |  2233 | $10.98   
  1/1/2014 | January Hosting          |  2233 | $54.95   
  1/4/2014 | February Hosting         |  2233 | $51.00   
  1/4/2014 | February Extra Bandwidth |  2233 | $30.00   
+----------+--------------------------+-------+---------+
                                        TOTAL | $146.93  
                                      +-------+---------+
This is a very long caption. The text should wrap to the
width of the table.
`
	checkEqual(t, buf.String(), want, "border table rendering failed")
}

func TestPrintLongCaptionWithLongExample(t *testing.T) {
	var buf bytes.Buffer
	data := [][]string{
		{"Learn East has computers with adapted keyboards with enlarged print etc", "Some Data", "Another Data"},
		{"Instead of lining up the letters all", "the way across, he splits the keyboard in two", "Like most ergonomic keyboards"},
	}

	table := NewWriter(&buf)
	table.SetCaption(true, "This is a very long caption. The text should wrap. If not, we have a problem that needs to be solved.")
	table.SetHeader([]string{"Name", "Sign", "Rating"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()

	want := `+--------------------------------+--------------------------------+-------------------------------+
|              NAME              |              SIGN              |            RATING             |
+--------------------------------+--------------------------------+-------------------------------+
| Learn East has computers       | Some Data                      | Another Data                  |
| with adapted keyboards with    |                                |                               |
| enlarged print etc             |                                |                               |
| Instead of lining up the       | the way across, he splits the  | Like most ergonomic keyboards |
| letters all                    | keyboard in two                |                               |
+--------------------------------+--------------------------------+-------------------------------+
This is a very long caption. The text should wrap. If not, we have a problem that needs to be
solved.
`
	checkEqual(t, buf.String(), want, "long caption for long example rendering failed")
}

func Example_autowrap() {
	var multiline = `A multiline
string with some lines being really long.`

	const (
		testRow = iota
		testHeader
		testFooter
		testFooter2
	)
	for mode := testRow; mode <= testFooter2; mode++ {
		for _, autoFmt := range []bool{false, true} {
			if mode == testRow && autoFmt {
				// Nothing special to test, skip
				continue
			}
			for _, autoWrap := range []bool{false, true} {
				for _, reflow := range []bool{false, true} {
					if !autoWrap && reflow {
						// Invalid configuration, skip
						continue
					}
					fmt.Println("mode", mode, "autoFmt", autoFmt, "autoWrap", autoWrap, "reflow", reflow)
					t := NewWriter(os.Stdout)
					t.SetAutoFormatHeaders(autoFmt)
					t.SetAutoWrapText(autoWrap)
					t.SetReflowDuringAutoWrap(reflow)
					if mode == testHeader {
						t.SetHeader([]string{"woo", multiline})
					} else {
						t.SetHeader([]string{"woo", "waa"})
					}
					if mode == testRow {
						t.Append([]string{"woo", multiline})
					} else {
						t.Append([]string{"woo", "waa"})
					}
					if mode == testFooter {
						t.SetFooter([]string{"woo", multiline})
					} else if mode == testFooter2 {
						t.SetFooter([]string{"", multiline})
					} else {
						t.SetFooter([]string{"woo", "waa"})
					}
					t.Render()
				}
			}
		}
		fmt.Println()
	}

	// Output:
	// mode 0 autoFmt false autoWrap false reflow false
	// +-----+-------------------------------------------+
	// | woo |                    waa                    |
	// +-----+-------------------------------------------+
	// | woo | A multiline                               |
	// |     | string with some lines being really long. |
	// +-----+-------------------------------------------+
	// | woo |                    waa                    |
	// +-----+-------------------------------------------+
	// mode 0 autoFmt false autoWrap true reflow false
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	// | woo | A multiline                    |
	// |     |                                |
	// |     | string with some lines being   |
	// |     | really long.                   |
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	// mode 0 autoFmt false autoWrap true reflow true
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	// | woo | A multiline string with some   |
	// |     | lines being really long.       |
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	//
	// mode 1 autoFmt false autoWrap false reflow false
	// +-----+-------------------------------------------+
	// | woo |                A multiline                |
	// |     | string with some lines being really long. |
	// +-----+-------------------------------------------+
	// | woo | waa                                       |
	// +-----+-------------------------------------------+
	// | woo |                    waa                    |
	// +-----+-------------------------------------------+
	// mode 1 autoFmt false autoWrap true reflow false
	// +-----+--------------------------------+
	// | woo |          A multiline           |
	// |     |                                |
	// |     |  string with some lines being  |
	// |     |          really long.          |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	// mode 1 autoFmt false autoWrap true reflow true
	// +-----+--------------------------------+
	// | woo |  A multiline string with some  |
	// |     |    lines being really long.    |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	// mode 1 autoFmt true autoWrap false reflow false
	// +-----+-------------------------------------------+
	// | WOO |                A MULTILINE                |
	// |     | STRING WITH SOME LINES BEING REALLY LONG  |
	// +-----+-------------------------------------------+
	// | woo | waa                                       |
	// +-----+-------------------------------------------+
	// | WOO |                    WAA                    |
	// +-----+-------------------------------------------+
	// mode 1 autoFmt true autoWrap true reflow false
	// +-----+--------------------------------+
	// | WOO |          A MULTILINE           |
	// |     |                                |
	// |     |  STRING WITH SOME LINES BEING  |
	// |     |          REALLY LONG           |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// | WOO |              WAA               |
	// +-----+--------------------------------+
	// mode 1 autoFmt true autoWrap true reflow true
	// +-----+--------------------------------+
	// | WOO |  A MULTILINE STRING WITH SOME  |
	// |     |    LINES BEING REALLY LONG     |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// | WOO |              WAA               |
	// +-----+--------------------------------+
	//
	// mode 2 autoFmt false autoWrap false reflow false
	// +-----+-------------------------------------------+
	// | woo |                    waa                    |
	// +-----+-------------------------------------------+
	// | woo | waa                                       |
	// +-----+-------------------------------------------+
	// | woo |                A multiline                |
	// |     | string with some lines being really long. |
	// +-----+-------------------------------------------+
	// mode 2 autoFmt false autoWrap true reflow false
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// | woo |          A multiline           |
	// |     |                                |
	// |     |  string with some lines being  |
	// |     |          really long.          |
	// +-----+--------------------------------+
	// mode 2 autoFmt false autoWrap true reflow true
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// | woo |  A multiline string with some  |
	// |     |    lines being really long.    |
	// +-----+--------------------------------+
	// mode 2 autoFmt true autoWrap false reflow false
	// +-----+-------------------------------------------+
	// | WOO |                    WAA                    |
	// +-----+-------------------------------------------+
	// | woo | waa                                       |
	// +-----+-------------------------------------------+
	// | WOO |                A MULTILINE                |
	// |     | STRING WITH SOME LINES BEING REALLY LONG  |
	// +-----+-------------------------------------------+
	// mode 2 autoFmt true autoWrap true reflow false
	// +-----+--------------------------------+
	// | WOO |              WAA               |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// | WOO |          A MULTILINE           |
	// |     |                                |
	// |     |  STRING WITH SOME LINES BEING  |
	// |     |          REALLY LONG           |
	// +-----+--------------------------------+
	// mode 2 autoFmt true autoWrap true reflow true
	// +-----+--------------------------------+
	// | WOO |              WAA               |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// | WOO |  A MULTILINE STRING WITH SOME  |
	// |     |    LINES BEING REALLY LONG     |
	// +-----+--------------------------------+
	//
	// mode 3 autoFmt false autoWrap false reflow false
	// +-----+-------------------------------------------+
	// | woo |                    waa                    |
	// +-----+-------------------------------------------+
	// | woo | waa                                       |
	// +-----+-------------------------------------------+
	// |                      A multiline                |
	// |       string with some lines being really long. |
	// +-----+-------------------------------------------+
	// mode 3 autoFmt false autoWrap true reflow false
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// |                A multiline           |
	// |                                      |
	// |        string with some lines being  |
	// |                really long.          |
	// +-----+--------------------------------+
	// mode 3 autoFmt false autoWrap true reflow true
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// |        A multiline string with some  |
	// |          lines being really long.    |
	// +-----+--------------------------------+
	// mode 3 autoFmt true autoWrap false reflow false
	// +-----+-------------------------------------------+
	// | WOO |                    WAA                    |
	// +-----+-------------------------------------------+
	// | woo | waa                                       |
	// +-----+-------------------------------------------+
	// |                      A MULTILINE                |
	// |       STRING WITH SOME LINES BEING REALLY LONG  |
	// +-----+-------------------------------------------+
	// mode 3 autoFmt true autoWrap true reflow false
	// +-----+--------------------------------+
	// | WOO |              WAA               |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// |                A MULTILINE           |
	// |                                      |
	// |        STRING WITH SOME LINES BEING  |
	// |                REALLY LONG           |
	// +-----+--------------------------------+
	// mode 3 autoFmt true autoWrap true reflow true
	// +-----+--------------------------------+
	// | WOO |              WAA               |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// |        A MULTILINE STRING WITH SOME  |
	// |          LINES BEING REALLY LONG     |
	// +-----+--------------------------------+
}

func TestPrintLine(t *testing.T) {
	header := make([]string, 12)
	val := " "
	want := ""
	for i := range header {
		header[i] = val
		want = fmt.Sprintf("%s+-%s-", want, strings.Replace(val, " ", "-", -1))
		val = val + " "
	}
	want = want + "+"
	var buf bytes.Buffer
	table := NewWriter(&buf)
	table.SetHeader(header)
	table.printLine(false)
	checkEqual(t, buf.String(), want, "line rendering failed")
}

func TestAnsiStrip(t *testing.T) {
	header := make([]string, 12)
	val := " "
	want := ""
	for i := range header {
		header[i] = "\033[43;30m" + val + "\033[00m"
		want = fmt.Sprintf("%s+-%s-", want, strings.Replace(val, " ", "-", -1))
		val = val + " "
	}
	want = want + "+"
	var buf bytes.Buffer
	table := NewWriter(&buf)
	table.SetHeader(header)
	table.printLine(false)
	checkEqual(t, buf.String(), want, "line rendering failed")
}

func NewCustomizedTable(out io.Writer) *Table {
	table := NewWriter(out)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetBorder(false)
	table.SetAlignment(ALIGN_LEFT)
	table.SetHeader([]string{})
	return table
}

func TestSubclass(t *testing.T) {
	buf := new(bytes.Buffer)
	table := NewCustomizedTable(buf)

	data := [][]string{
		{"A", "The Good", "500"},
		{"B", "The Very very Bad Man", "288"},
		{"C", "The Ugly", "120"},
		{"D", "The Gopher", "800"},
	}

	for _, v := range data {
		table.Append(v)
	}
	table.Render()

	want := `  A  The Good               500  
  B  The Very very Bad Man  288  
  C  The Ugly               120  
  D  The Gopher             800  
`
	checkEqual(t, buf.String(), want, "test subclass failed")
}

func TestAutoMergeRows(t *testing.T) {
	data := [][]string{
		{"A", "The Good", "500"},
		{"A", "The Very very Bad Man", "288"},
		{"B", "The Very very Bad Man", "120"},
		{"B", "The Very very Bad Man", "200"},
	}
	var buf bytes.Buffer
	table := NewWriter(&buf)
	table.SetHeader([]string{"Name", "Sign", "Rating"})

	for _, v := range data {
		table.Append(v)
	}
	table.SetAutoMergeCells(true)
	table.Render()
	want := `+------+-----------------------+--------+
| NAME |         SIGN          | RATING |
+------+-----------------------+--------+
| A    | The Good              |    500 |
|      | The Very very Bad Man |    288 |
| B    |                       |    120 |
|      |                       |    200 |
+------+-----------------------+--------+
`
	got := buf.String()
	if got != want {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}

	buf.Reset()
	table = NewWriter(&buf)
	table.SetHeader([]string{"Name", "Sign", "Rating"})

	for _, v := range data {
		table.Append(v)
	}
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.Render()
	want = `+------+-----------------------+--------+
| NAME |         SIGN          | RATING |
+------+-----------------------+--------+
| A    | The Good              |    500 |
+      +-----------------------+--------+
|      | The Very very Bad Man |    288 |
+------+                       +--------+
| B    |                       |    120 |
+      +                       +--------+
|      |                       |    200 |
+------+-----------------------+--------+
`
	checkEqual(t, buf.String(), want)

	buf.Reset()
	table = NewWriter(&buf)
	table.SetHeader([]string{"Name", "Sign", "Rating"})

	dataWithlongText := [][]string{
		{"A", "The Good", "500"},
		{"A", "The Very very very very very Bad Man", "288"},
		{"B", "The Very very very very very Bad Man", "120"},
		{"C", "The Very very Bad Man", "200"},
	}
	table.AppendBulk(dataWithlongText)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.Render()
	want = `+------+--------------------------------+--------+
| NAME |              SIGN              | RATING |
+------+--------------------------------+--------+
| A    | The Good                       |    500 |
+------+--------------------------------+--------+
| A    | The Very very very very very   |    288 |
|      | Bad Man                        |        |
+------+                                +--------+
| B    |                                |    120 |
|      |                                |        |
+------+--------------------------------+--------+
| C    | The Very very Bad Man          |    200 |
+------+--------------------------------+--------+
`
	checkEqual(t, buf.String(), want)
}

func TestClearRows(t *testing.T) {
	data := [][]string{
		{"1/1/2014", "Domain name", "2233", "$10.98"},
	}

	var buf bytes.Buffer
	table := NewWriter(&buf)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Date", "Description", "CV2", "Amount"})
	table.SetFooter([]string{"", "", "Total", "$145.93"}) // Add Footer
	table.AppendBulk(data)                                // Add Bulk Data
	table.Render()

	originalWant := `+----------+-------------+-------+---------+
|   DATE   | DESCRIPTION |  CV2  | AMOUNT  |
+----------+-------------+-------+---------+
| 1/1/2014 | Domain name |  2233 | $10.98  |
+----------+-------------+-------+---------+
|                          TOTAL | $145.93 |
+----------+-------------+-------+---------+
`
	want := originalWant

	checkEqual(t, buf.String(), want, "table clear rows failed")

	buf.Reset()
	table.ClearRows()
	table.Render()

	want = `+----------+-------------+-------+---------+
|   DATE   | DESCRIPTION |  CV2  | AMOUNT  |
+----------+-------------+-------+---------+
+----------+-------------+-------+---------+
|                          TOTAL | $145.93 |
+----------+-------------+-------+---------+
`

	checkEqual(t, buf.String(), want, "table clear rows failed")

	buf.Reset()
	table.AppendBulk(data) // Add Bulk Data
	table.Render()

	want = `+----------+-------------+-------+---------+
|   DATE   | DESCRIPTION |  CV2  | AMOUNT  |
+----------+-------------+-------+---------+
| 1/1/2014 | Domain name |  2233 | $10.98  |
+----------+-------------+-------+---------+
|                          TOTAL | $145.93 |
+----------+-------------+-------+---------+
`

	checkEqual(t, buf.String(), want, "table clear rows failed")
}

func TestClearFooters(t *testing.T) {
	data := [][]string{
		{"1/1/2014", "Domain name", "2233", "$10.98"},
	}

	var buf bytes.Buffer
	table := NewWriter(&buf)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Date", "Description", "CV2", "Amount"})
	table.SetFooter([]string{"", "", "Total", "$145.93"}) // Add Footer
	table.AppendBulk(data)                                // Add Bulk Data
	table.Render()

	buf.Reset()
	table.ClearFooter()
	table.Render()

	want := `+----------+-------------+-------+---------+
|   DATE   | DESCRIPTION |  CV2  | AMOUNT  |
+----------+-------------+-------+---------+
| 1/1/2014 | Domain name |  2233 | $10.98  |
+----------+-------------+-------+---------+
`

	checkEqual(t, buf.String(), want)
}

func TestMoreDataColumnsThanHeaders(t *testing.T) {
	var (
		buf    = &bytes.Buffer{}
		table  = NewWriter(buf)
		header = []string{"A", "B", "C"}
		data   = [][]string{
			{"a", "b", "c", "d"},
			{"1", "2", "3", "4"},
		}
		want = `+---+---+---+---+
| A | B | C |   |
+---+---+---+---+
| a | b | c | d |
| 1 | 2 | 3 | 4 |
+---+---+---+---+
`
	)
	table.SetHeader(header)
	// table.SetFooter(ctx.tableCtx.footer)
	table.AppendBulk(data)
	table.Render()

	checkEqual(t, buf.String(), want)
}

func TestMoreFooterColumnsThanHeaders(t *testing.T) {
	var (
		buf    = &bytes.Buffer{}
		table  = NewWriter(buf)
		header = []string{"A", "B", "C"}
		data   = [][]string{
			{"a", "b", "c", "d"},
			{"1", "2", "3", "4"},
		}
		footer = []string{"a", "b", "c", "d", "e"}
		want   = `+---+---+---+---+---+
| A | B | C |   |   |
+---+---+---+---+---+
| a | b | c | d |
| 1 | 2 | 3 | 4 |
+---+---+---+---+---+
| A | B | C | D | E |
+---+---+---+---+---+
`
	)
	table.SetHeader(header)
	table.SetFooter(footer)
	table.AppendBulk(data)
	table.Render()

	checkEqual(t, buf.String(), want)
}

func TestSetColMinWidth(t *testing.T) {
	var (
		buf    = &bytes.Buffer{}
		table  = NewWriter(buf)
		header = []string{"AAA", "BBB", "CCC"}
		data   = [][]string{
			{"a", "b", "c"},
			{"1", "2", "3"},
		}
		footer = []string{"a", "b", "cccc"}
		want   = `+-----+-----+-------+
| AAA | BBB |  CCC  |
+-----+-----+-------+
| a   | b   | c     |
|   1 |   2 |     3 |
+-----+-----+-------+
|  A  |  B  | CCCC  |
+-----+-----+-------+
`
	)
	table.SetHeader(header)
	table.SetFooter(footer)
	table.AppendBulk(data)
	table.SetColMinWidth(2, 5)
	table.Render()

	checkEqual(t, buf.String(), want)
}

func TestWrapString(t *testing.T) {
	want := []string{"ああああああああああああああああああああああああ", "あああああああ"}
	got, _ := WrapString("ああああああああああああああああああああああああ あああああああ", 55)
	checkEqual(t, got, want)
}

func TestCustomAlign(t *testing.T) {
	var (
		buf    = &bytes.Buffer{}
		table  = NewWriter(buf)
		header = []string{"AAA", "BBB", "CCC"}
		data   = [][]string{
			{"a", "b", "c"},
			{"1", "2", "3"},
		}
		footer = []string{"a", "b", "cccc"}
		want   = `+-----+-----+-------+
| AAA | BBB |  CCC  |
+-----+-----+-------+
| a   |  b  |     c |
| 1   |  2  |     3 |
+-----+-----+-------+
|  A  |  B  | CCCC  |
+-----+-----+-------+
`
	)
	table.SetHeader(header)
	table.SetFooter(footer)
	table.AppendBulk(data)
	table.SetColMinWidth(2, 5)
	table.SetColumnAlignment([]int{ALIGN_LEFT, ALIGN_CENTER, ALIGN_RIGHT})
	table.Render()

	checkEqual(t, buf.String(), want)
}

func TestTitle(t *testing.T) {
	ts := []struct {
		text string
		want string
	}{
		{"", ""},
		{"foo", "FOO"},
		{"Foo", "FOO"},
		{"foO", "FOO"},
		{".foo", "FOO"},
		{"foo.", "FOO"},
		{".foo.", "FOO"},
		{".foo.bar.", "FOO BAR"},
		{"_foo", "FOO"},
		{"foo_", "FOO"},
		{"_foo_", "FOO"},
		{"_foo_bar_", "FOO BAR"},
		{" foo", "FOO"},
		{"foo ", "FOO"},
		{" foo ", "FOO"},
		{" foo bar ", "FOO BAR"},
		{"0.1", "0.1"},
		{"FOO 0.1", "FOO 0.1"},
		{".1 0.1", ".1 0.1"},
		{"1. 0.1", "1. 0.1"},
		{"1. 0.", "1. 0."},
		{".1. 0.", ".1. 0."},
		{".$ . $.", "$ . $"},
		{".$. $.", "$  $"},
	}
	for _, tt := range ts {
		got := Title(tt.text)
		if got != tt.want {
			t.Errorf("want %q, bot got %q", tt.want, got)
		}
	}
}
