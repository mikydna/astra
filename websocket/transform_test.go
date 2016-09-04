package websocket

import (
	"testing"
	"time"
)

func TestDownsample_4x4(t *testing.T) {
	index := 0
	timestamp := time.Now().Unix()

	testFrame := &Frame{
		Metadata: FrameMetadata{
			Index:     index,
			Timestamp: timestamp,
			Width:     4,
			Height:    4,
		},
		Data: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
	}

	Downsample(2)(testFrame)

	if testFrame.Metadata.Index != index {
		t.Fail()
	}

	if testFrame.Metadata.Timestamp != timestamp {
		t.Fail()
	}

	if testFrame.Metadata.Width != 2 {
		t.Fail()
	}

	if testFrame.Metadata.Height != 2 {
		t.Fail()
	}

	expectedData := []int{1, 3, 9, 11}
	for i, val := range testFrame.Data {
		if expectedData[i] != val {
			t.Fail()
		}
	}

}
