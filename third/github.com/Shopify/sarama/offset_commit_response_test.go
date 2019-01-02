package sarama

import (
	"fmt"
	"testing"
)

var (
	emptyOffsetCommitResponse = []byte{
		0x00, 0x00, 0x00, 0x00}
)

func TestEmptyOffsetCommitResponse(t *testing.T) {
	response := OffsetCommitResponse{}
	testResponse(t, "empty", &response, emptyOffsetCommitResponse)
}

func TestNormalOffsetCommitResponse(t *testing.T) {
	response := OffsetCommitResponse{}
	response.AddError("t", 0, ErrNotLeaderForPartition)
	response.Errors["m"] = make(map[int32]KError)
	// The response encoded form cannot be checked for it varies due to
	// unpredictable map traversal order.
	testResponse(t, "normal", &response, nil)
}

func TestOffsetCommitResponseWithThrottleTime(t *testing.T) {
	for version := 3; version <= 4; version++ {
		response := OffsetCommitResponse{
			Version:        int16(version),
			ThrottleTimeMs: 123,
		}
		response.AddError("t", 0, ErrNotLeaderForPartition)
		response.Errors["m"] = make(map[int32]KError)
		// The response encoded form cannot be checked for it varies due to
		// unpredictable map traversal order.
		testResponse(t, fmt.Sprintf("v%d with throttle time", version), &response, nil)
	}
}
