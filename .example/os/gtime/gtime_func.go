package main

import (
	"fmt"

	"github.com/jin502437344/gf/os/gtime"
)

func main() {
	fmt.Println("Date       :", gtime.Date())
	fmt.Println("Datetime   :", gtime.Datetime())
	fmt.Println("Second     :", gtime.Timestamp())
	fmt.Println("Millisecond:", gtime.TimestampMilli())
	fmt.Println("Microsecond:", gtime.TimestampMicro())
	fmt.Println("Nanosecond :", gtime.TimestampNano())
}
