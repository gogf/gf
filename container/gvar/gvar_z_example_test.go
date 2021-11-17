// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
package gvar_test

import (
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"time"
)

type Product struct {
	Code          string    `json:"Code"`
	Name          string    `json:"Name"`
	Quantity      int32     `json:"Quantity"`
	Price         float64   `json:"Price"`
	OnSale        bool      `json:"OnSale"`
	PromotionTime time.Time `json:"PromotionTime"`
}

var product1 = Product{
	Code:     "0001",
	Name:     "Golang Book1",
	Quantity: 1000,
	Price:    100.00,
	OnSale:   true,
}

var product2 = Product{
	Code:     "0002",
	Name:     "GoFrame Book2",
	Quantity: 2000,
	Price:    200.00,
	OnSale:   false,
}

// New
func ExampleVar_New() {
	var a1 = gvar.New(product1.Price, false)
	g.Dump(a1)
	var a2 = a1.Set(product2.Price).(float64)
	g.Dump(a1)
	g.Dump(a2)

	var b1 = gvar.New(product1.Price, true)
	g.Dump(b1)
	var b2 = b1.Set(product2.Name).(float64)
	g.Dump(b1)
	g.Dump(b2)

	// Output:
	// "100"
	// "200"
	// 100
	// "100"
	// "GoFrame Book2"
	// 100
}

// Create
func ExampleVar_Create() {
	var a1 = gvar.Create(product1.Price, false)
	g.Dump(&a1)
	var a2 = a1.Set(product2.Price).(float64)
	g.Dump(&a1)
	g.Dump(a2)

	var b1 = gvar.Create(product1.Price, true)
	g.Dump(&b1)
	var b2 = b1.Set(product2.Name).(float64)
	g.Dump(&b1)
	g.Dump(b2)

	// Output:
	// "100"
	// "200"
	// 100
	// "100"
	// "GoFrame Book2"
	// 100
}

// Clone
func ExampleVar_Clone() {
	var a1 = gvar.New(product1, false)
	var a2 = a1.Clone()
	var p2 = a2.Val().(Product)
	p2.PromotionTime, _ = time.Parse("2006-Jan-02", "2021-Nov-11")
	g.Dump(p2)

	// Output:
	// {
	//     Code:          "0001",
	//     Name:          "Golang Book1",
	//     Quantity:      1000,
	//     Price:         100,
	//     OnSale:        true,
	//     PromotionTime: "2021-11-11 00:00:00 +0000 UTC",
	// }
}

// Set
func ExampleVar_Set() {
	var tag1 = gvar.New(product1.Price, false)
	_ = tag1.Set(product2.Price).(float64)
	g.Dump(tag1)

	var tag2 = gvar.New(product1.Price, true)
	_ = tag2.Set(product2.Name).(float64)
	g.Dump(tag2)

	// Output:
	// "200"
	// "GoFrame Book2"
}

// Val
func ExampleVar_Val() {
	var a1 = gvar.New(product1.Price, false)
	g.Dump(a1.Val().(float64) * 100)

	// Output:
	// 10000
}

// Interface
func ExampleVar_Interface() {
	var a1 = gvar.New(product1.Price, false)
	g.Dump(a1.Interface().(float64) * 100)

	// Output:
	// 10000
}

// Bytes
func ExampleVar_Bytes() {
	var a1 = gvar.New(product1.Name, false)
	g.Dump(len(a1.Bytes()))
	a1.Set(a1.Val().(string) + "从入门到精通")
	g.Dump(len(a1.Bytes()))

	// Output:
	// 12
	// 30
}

// String
func ExampleVar_String() {
	var a1 = gvar.New(product1.Name, false)
	g.Dump(len(a1.String()))
	a1.Set(a1.String() + "从入门到精通")
	g.Dump(len(a1.String()))

	// Output:
	// 12
	// 30
}

// Bool
func ExampleVar_Bool() {
	var a1 = gvar.New(product1.OnSale, false)
	g.Dump(a1.Bool())

	// Output:
	// true
}

// Int
func ExampleVar_Int() {
	var a1 = gvar.New(-1000, true)
	g.Dump(a1.Int())

	// Output:
	// -1000
}

// Int8
func ExampleVar_Int8() {
	var a1 = gvar.New(-100, true)
	g.Dump(a1.Int8())

	// Output:
	// -100
}

// Int16
func ExampleVar_Int16() {
	var a1 = gvar.New(-10000, true)
	g.Dump(a1.Int16())

	// Output:
	// -10000
}

// Int32
func ExampleVar_Int32() {
	var a1 = gvar.New(-10000, true)
	g.Dump(a1.Int32())

	// Output:
	// -10000
}

// Int64
func ExampleVar_Int64() {
	var a1 = gvar.New(-100000000, true)
	g.Dump(a1.Int64())

	// Output:
	// -100000000
}

// Uint
func ExampleVar_Uint() {
	var a1 = gvar.New(1000, true)
	g.Dump(a1.Uint())

	// Output:
	// 1000
}

// Uint8
func ExampleVar_Uint8() {
	var a1 = gvar.New(100, true)
	g.Dump(a1.Uint8())

	// Output:
	// 100
}

// Uint16
func ExampleVar_Uint16() {
	var a1 = gvar.New(10000, true)
	g.Dump(a1.Uint16())

	// Output:
	// 10000
}

// Uint32
func ExampleVar_Uint32() {
	var a1 = gvar.New(100000, true)
	g.Dump(a1.Uint32())

	// Output:
	// 100000
}

// Uint64
func ExampleVar_Uint64() {
	var a1 = gvar.New(10000000, true)
	g.Dump(a1.Uint64())

	// Output:
	// 10000000
}

// Float32
func ExampleVar_Float32() {
	var price = gvar.New(product1.Price)
	g.Dump(price.Float32())

	// Output:
	// 100
}

// Float64
func ExampleVar_Float64() {
	var price = gvar.New(product1.Price)
	var quantity = gvar.New(product1.Quantity)
	var money = price.Float64() * quantity.Float64()
	g.Dump(money)

	// Output:
	// 100000
}

// Time
func ExampleVar_Time() {
	product1.PromotionTime, _ = time.Parse("2006-Jan-02", "2021-Nov-11")
	var a1 = gvar.New(product1.PromotionTime)
	g.Dump(a1.Time())
	g.Dump(a1.Time().Unix())

	// Output:
	// "2021-11-11 00:00:00 +0000 UTC"
	// 1636588800
}

// Duration
func ExampleVar_Duration() {
	var a1 = gvar.New("300s")
	g.Dump(a1.Duration())

	// Output:
	// 5m0s
}

// GTime
func ExampleVar_GTime() {
	product1.PromotionTime, _ = time.Parse("2006-Jan-02", "2021-Nov-11")
	var a1 = gvar.New(product1.PromotionTime)
	g.Dump(a1.GTime())
	g.Dump(a1.GTime().Unix())

	// Output:
	// "2021-11-11 00:00:00"
	// 1636588800
}

// MarshalJSON
func ExampleVar_MarshalJSON() {
	var a1 = gvar.New(product1)
	var json, _ = a1.MarshalJSON()
	g.Dump(json)

	// Output:
	// "{\"Code\":\"0001\",\"Name\":\"Golang Book1\",\"Quantity\":1000,\"Price\":100,\"OnSale\":true,\"PromotionTime\":\"2021-11-11T00:00:00Z\"}"
}

// UnmarshalJSON
func ExampleVar_UnmarshalJSON() {
	json := []byte(`{
	     "Code":          "0003",
	     "Name":          "Golang Book3",
	     "Quantity":      3000,
	     "Price":         300,
	     "OnSale":        true
}`)
	var a1 = gvar.New(&Product{})
	_ = a1.UnmarshalJSON(json)
	g.Dump(a1)

	// Output:
	// "{\"Code\":\"0003\",\"Name\":\"Golang Book3\",\"OnSale\":true,\"Price\":300,\"Quantity\":3000}"
}

// UnmarshalValue
func ExampleVar_UnmarshalValue() {
	var a1 = gvar.New(&Product{})
	_ = a1.UnmarshalValue(product2)
	g.Dump(a1)

	// Output:
	// "{\"Code\":\"0002\",\"Name\":\"GoFrame Book2\",\"Quantity\":2000,\"Price\":200,\"OnSale\":false,\"PromotionTime\":\"0001-01-01T00:00:00Z\"}"
}
