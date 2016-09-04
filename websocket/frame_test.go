package websocket

import (
	"testing"
)

import (
	"github.com/mikydna/astra"
)

func TestFromDepthFrame(t *testing.T) {
	testDepthFrame := astra.CameraDepthFrame{
		Index:  0,
		Width:  4,
		Height: 4,
		Buffer: []int16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
	}

	websocketFrame, err := FromDepthFrame(testDepthFrame)
	if err != nil {
		t.Fatal(err)
	}

	if websocketFrame.Index != testDepthFrame.Index {
		t.Errorf("Did not set index correctly: %d != %d", websocketFrame.Index, testDepthFrame.Index)
	}

	if websocketFrame.Width != testDepthFrame.Width {
		t.Errorf("Did not set width correctly: %d != %d", websocketFrame.Width, testDepthFrame.Width)
	}

	if websocketFrame.Height != testDepthFrame.Height {
		t.Errorf("Did not set height correctly: %d != %d", websocketFrame.Height, testDepthFrame.Height)
	}

	if len(websocketFrame.Data) != len(testDepthFrame.Buffer) {
		t.Errorf("Did not set data correctly: %d != %d", len(websocketFrame.Data), len(testDepthFrame.Buffer))
	}

	for i, bufferVal := range testDepthFrame.Buffer {
		if int(bufferVal) != websocketFrame.Data[i] {
			t.Errorf("Values are not equal: @=%d, %d != %d", i, bufferVal, websocketFrame.Data[i])
			t.FailNow()
		}
	}

}
