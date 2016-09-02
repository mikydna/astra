package edge

import (
	"fmt"
	"testing"
	"time"
)

func TestEdge(t *testing.T) {
	cameras := []string{"device/default"}

	edge, err := NewEdge(cameras)
	if err != nil {
		t.Fatal(err)
	}
	defer edge.Shutdown()

	go edge.Start()

	timeout := time.After(5 * time.Second)
	alive := true
	for alive {
		select {
		case <-timeout:
			alive = false
		case frame := <-edge.Depth:
			fmt.Printf("FRAME %d %d %d %d\n", frame.Index, frame.Width, frame.Height, len(frame.Buffer))
		}
	}

}
