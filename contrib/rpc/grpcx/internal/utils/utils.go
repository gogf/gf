// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package utils provides utilities for GRPC.
package utils

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	protoJSONMarshaller = &protojson.MarshalOptions{
		EmitUnpopulated: true,
	}
)

// MarshalPbMessageToJsonString marshals protobuf message to json string.
func MarshalPbMessageToJsonString(msg proto.Message) string {
	return protoJSONMarshaller.Format(msg)
}

func MarshalMessageToJsonStringForTracing(value interface{}, msgType string, maxBytes int) string {
	var messageContent string
	if msg, ok := value.(proto.Message); ok {
		if proto.Size(msg) <= maxBytes {
			messageContent = MarshalPbMessageToJsonString(msg)
		} else {
			messageContent = fmt.Sprintf(
				"[%s Message Too Large For Tracing, Max: %d bytes]",
				msgType,
				maxBytes,
			)
		}
	} else {
		messageContent = fmt.Sprintf("%v", value)
	}
	return messageContent
}
