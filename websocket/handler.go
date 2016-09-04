package websocket

import (
	"log"
	"net/http"
)

import (
	"github.com/mikydna/astra"
	"golang.org/x/net/websocket"
)

func BroadcastFrames(edge *astra.Edge, procs ...FrameProcessor) http.Handler {
	demux := []chan Frame{}

	go func() {
		for capture := range edge.Depth {
			frame, err := FromDepthFrame(capture)
			if err != nil {
				log.Println(err)
				continue
			}

			for _, f := range procs {
				f(frame)
			}

			for i, out := range demux {
				select {
				case out <- *frame:
				default:
					log.Printf("dropping frame for out %d", i)
				}
			}
		}
	}()

	socketHandler := func(conn *websocket.Conn) {
		publish := make(chan Frame)
		demux = append(demux, publish)
		at := len(demux) - 1

		defer func() {
			demux = append(demux[:at], demux[at:]...)
		}()

		for frame := range publish {
			if err := websocket.JSON.Send(conn, frame); err != nil {
				break
			}
		}
	}

	return websocket.Handler(socketHandler)
}
