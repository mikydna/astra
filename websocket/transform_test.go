package websocket

import (
	"testing"
)

func TestDownsample_4x4(t *testing.T) {
	frame := &Frame{
		Index:     0,
		Timestamp: 1234,
		Width:     4,
		Height:    4,
		Data:      []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
	}

	expected := &Frame{
		Index:     0,
		Timestamp: 1234,
		Width:     2,
		Height:    2,
		Data:      []int{1, 3, 9, 11},
	}

	// action
	Downsample(2)(frame)

	// test
	if frame.Index != expected.Index {
		t.Fail()
	}

	if frame.Timestamp != expected.Timestamp {
		t.Fail()
	}

	if frame.Width != expected.Width {
		t.Fail()
	}

	if frame.Height != expected.Height {
		t.Fail()
	}

	for i, val := range frame.Data {
		if expected.Data[i] != val {
			t.Fail()
		}
	}

}
