package sarama

import (
	"fmt"
	"testing"
)

var (
	emptyOffsetFetchResponse = []byte{
		0x00, 0x00, 0x00, 0x00}

	emptyOffsetFetchResponseV2 = []byte{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x2A}

	emptyOffsetFetchResponseV3 = []byte{
		0x00, 0x00, 0x00, 0x09,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x2A}
)

func TestEmptyOffsetFetchResponse(t *testing.T) {
	for version := 0; version <= 1; version++ {
		response := OffsetFetchResponse{Version: int16(version)}
		testResponse(t, fmt.Sprintf("empty v%d", version), &response, emptyOffsetFetchResponse)
	}

	responseV2 := OffsetFetchResponse{Version: 2, Err: ErrInvalidRequest}
	testResponse(t, "empty V2", &responseV2, emptyOffsetFetchResponseV2)

	for version := 3; version <= 5; version++ {
		responseV3 := OffsetFetchResponse{Version: int16(version), Err: ErrInvalidRequest, ThrottleTimeMs: 9}
		testResponse(t, fmt.Sprintf("empty v%d", version), &responseV3, emptyOffsetFetchResponseV3)
	}
}

func TestNormalOffsetFetchResponse(t *testing.T) {
	// The response encoded form cannot be checked for it varies due to
	// unpredictable map traversal order.
	// Hence the 'nil' as byte[] parameter in the 'testResponse(..)' calls

	for version := 0; version <= 1; version++ {
		response := OffsetFetchResponse{Version: int16(version)}
		response.AddBlock("t", 0, &OffsetFetchResponseBlock{0, 0, "md", ErrRequestTimedOut})
		response.Blocks["m"] = nil
		testResponse(t, fmt.Sprintf("Normal v%d", version), &response, nil)
	}

	responseV2 := OffsetFetchResponse{Version: 2, Err: ErrInvalidRequest}
	responseV2.AddBlock("t", 0, &OffsetFetchResponseBlock{0, 0, "md", ErrRequestTimedOut})
	responseV2.Blocks["m"] = nil
	testResponse(t, "normal V2", &responseV2, nil)

	for version := 3; version <= 4; version++ {
		responseV3 := OffsetFetchResponse{Version: int16(version), Err: ErrInvalidRequest, ThrottleTimeMs: 9}
		responseV3.AddBlock("t", 0, &OffsetFetchResponseBlock{0, 0, "md", ErrRequestTimedOut})
		responseV3.Blocks["m"] = nil
		testResponse(t, fmt.Sprintf("Normal v%d", version), &responseV3, nil)
	}

	responseV5 := OffsetFetchResponse{Version: 5, Err: ErrInvalidRequest, ThrottleTimeMs: 9}
	responseV5.AddBlock("t", 0, &OffsetFetchResponseBlock{Offset: 10, LeaderEpoch: 100, Metadata: "md", Err: ErrRequestTimedOut})
	responseV5.Blocks["m"] = nil
	testResponse(t, "normal V5", &responseV5, nil)
}
