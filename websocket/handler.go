package websocket

import (
	"log"
	"net/http"
)

import (
	"github.com/mikydna/astra"
	"golang.org/x/net/websocket"
)

func BroadcastFrames(in <-chan astra.CameraDepthFrame, procs ...astra.FrameProcessor) http.Handler {
	demux := []chan astra.Frame{}

	go func() {
		for capture := range in {
			frame, err := astra.FromDepthFrame(capture)
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
		publish := make(chan astra.Frame)
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
