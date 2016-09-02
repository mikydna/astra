package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

import (
	"github.com/facebookgo/httpdown"
	"github.com/mikydna/astra"
	"golang.org/x/net/websocket"
)

var (
	httpConf = &httpdown.HTTP{
		StopTimeout: 5 * time.Second,
		KillTimeout: 1 * time.Second,
	}
)

type ProcessedFrame struct {
	Index  int
	Width  int
	Height int
	Buffer []int16
}

func rle(buffer []int16) []int {
	rle := []int{}
	count := 1
	curr := buffer[0]

	rle = append(rle, int(curr))
	for i := 1; i < len(buffer); i++ {
		if buffer[i] == curr {
			count += 1
		} else {
			rle = append(rle, int(count), int(buffer[i]))

			curr = buffer[i]
			count = 1
		}
	}

	rle = append(rle, int(count))

	return rle
}

func downsample(frame astra.CameraDepthFrame, factor uint8) []int16 {
	result := []int16{}
	for r := 0; r < int(frame.Height); r += 2 {
		for c := 0; c < int(frame.Width); c += 2 {
			offset := r*int(frame.Width) + c
			result = append(result, frame.Buffer[offset])
		}
	}

	return result
}

func BroadcastDepthFrames(edge *astra.Edge) http.HandlerFunc {
	subscribers := []chan []int{}

	// consume edge frames
	go func() {
		for depthFrame := range edge.Depth {
			for _, subscriber := range subscribers {
				select {
				case subscriber <- rle(downsample(depthFrame, 2)):
				default:
				}
			}
		}
	}()

	broadcast := func(conn *websocket.Conn) {
		log.Println("CONN MADE")

		listen := make(chan []int)
		subscribers = append(subscribers, listen)
		at := len(subscribers) - 1

		for frame := range listen {
			temp := struct {
				Frame  int   `json:"frame"`
				Buffer []int `json:"data"`
			}{
				Frame:  -1,
				Buffer: frame,
			}

			b, err := json.Marshal(temp)
			if err != nil {
				break
			}

			if _, err := conn.Write(b); err != nil {
				subscribers = append(subscribers[:at], subscribers[at:]...)
				break
			}
		}

		log.Println("CONN EXIT")
	}

	return websocket.Handler(broadcast).ServeHTTP
}

func main() {
	done := make(chan bool)
	// go func() {
	//  <-time.After(100 * time.Second)
	//  done <- true
	// }()

	// initialize cameras
	cameras := []string{"device/default"}
	edge, err := astra.NewEdge(cameras)
	if err != nil {
		log.Fatal(err)
	}
	defer edge.Shutdown()

	// web server
	mux := http.NewServeMux()
	server := &http.Server{Addr: ":9091", Handler: mux}

	mux.Handle("/", BroadcastDepthFrames(edge))

	graceful, err := httpConf.ListenAndServe(server)
	if err != nil {
		return
	}
	defer graceful.Stop()

	// start
	go edge.Start()

	<-done
}
