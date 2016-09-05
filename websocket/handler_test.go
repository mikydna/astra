package websocket

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

import (
	"github.com/mikydna/astra"
	"golang.org/x/net/websocket"
)

func generateTestCaptureFrames(n, width, height int) []astra.CameraDepthFrame {
	gen := make([]astra.CameraDepthFrame, n)

	for i := 0; i < n; i++ {
		gen[i].Index = i
		gen[i].Width = width
		gen[i].Height = height

		gen[i].Buffer = make([]int16, width*height)
		for j := 0; j < width*height; j++ {
			gen[i].Buffer[j] = int16((i + 1) * (j + 1))
		}
	}

	return gen
}

func TestBroadcastFrames(t *testing.T) {
	// setup fake data
	n := 10
	sleep := 250 // ms
	testDepthFrames := generateTestCaptureFrames(n, 4, 4)

	// setup edge
	fakeEdge := &astra.Edge{
		Depth: make(chan astra.CameraDepthFrame),
	}

	// setup socket server
	http.Handle("/broadcast", BroadcastFrames(fakeEdge.Depth, Downsample(2)))
	server := httptest.NewServer(nil)
	addr := server.Listener.Addr().String()
	defer server.Close()

	// connect
	url := fmt.Sprintf("ws://%s/broadcast", addr)
	conn, err := websocket.Dial(url, "", "http://localhost/")
	if err != nil {
		t.Fail()
	}

	// send fake frames
	go func() {
		for _, sendFrame := range testDepthFrames {
			time.Sleep(time.Duration(sleep) * time.Millisecond)
			fakeEdge.Depth <- sendFrame
		}
	}()

	// action
	time.AfterFunc(time.Duration(n*sleep+100)*time.Millisecond, t.Fail)
	for i := 0; i < n; i++ {
		received := &Frame{}
		err := websocket.JSON.Receive(conn, received)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("Received frame, index=%d, len=%d", received.Index, len(received.Data))
	}
}
