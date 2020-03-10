package main

import (
	"fmt"
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/encoding/gparser"
	"github.com/gogf/gf/os/gtime"
)

type EntityTest struct {
	ID                        int64
	DeviceID                  string
	DeviceType                uint8
	ProtoVer                  uint8
	DataType                  uint8
	RecordTime                *gtime.Time
	CreateTime                *gtime.Time
	RemoteIP                  string
	Voltage                   int16
	Battery                   uint8
	CellularVersion           uint8
	StateID                   uint16
	OperatorID                uint16
	RegionID                  uint16
	BaseStationID             uint32
	BaseStationSignalStrength int8
	ContainerStatus           int8
	PositionX                 int16
	PositionY                 int16
	PositionZ                 int16
	ContainerID               string
	NetworkRegisterTime       int16
	Temperature               int8
	Humidity                  uint8
	SleepMode                 uint8
	CellularChangeStatus      int8
	GPSSearchTime             int16
	Latitude                  float64
	Longitude                 float64
	Altitude                  int16
	Speed                     uint8
	DateStr                   string
	TimeStr                   string
	FrameID                   uint16
}

func main() {
	array := []string{"DeviceID", "8134567890ABCDEF", "DeviceType", "16", "ProtoVer", "11", "DataType", "2", "RecordTime", "2020-02-27 15:23:23", "CreateTime", "2020-03-09 19:24:22", "RemoteIP", "127.0.0.1:60471", "Voltage", "3610", "Battery", "90", "CellularVersion", "4", "StateID", "401", "OperatorID", "345", "RegionID", "17716", "BaseStationID", "571479331", "BaseStationSignalStrength", "45", "ContainerStatus", "1", "PositionX", "40", "PositionY", "-25", "PositionZ", "16", "ContainerID", "aAvVDEFPQAA", "DateStr", "270220", "TimeStr", "232323", "NetworkRegisterTime", "60", "Temperature", "25", "Humidity", "70", "SleepMode", "1", "CellularChangeStatus", "15", "GPSSearchTime", "32767", "Latitude", "116.23532", "Longitude", "39.21526", "Altitude", "2800", "Speed", "66", "FrameID", "4660"}
	v := gvar.New(array)
	o := &EntityTest{}
	v.Struct(o)
	//b, _ := json.Marshal(o)
	fmt.Println(gparser.MustToJsonIndentString(o))
}
