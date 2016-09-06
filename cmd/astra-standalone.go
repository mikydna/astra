package main

import (
	"log"
	"net/http"
	"time"
)

import (
	"github.com/facebookgo/httpdown"
	"github.com/mikydna/astra"
	"github.com/mikydna/astra/tool/preview"
	"github.com/mikydna/astra/websocket"
)

var (
	httpConf = &httpdown.HTTP{
		StopTimeout: 5 * time.Second,
		KillTimeout: 1 * time.Second,
	}
)

func main() {
	done := make(chan bool)

	// initialize cameras
	cameras := []string{"device/default"}
	edge, err := astra.NewEdge(cameras)
	if err != nil {
		log.Fatal(err)
	}
	defer edge.Shutdown()

	go edge.Start()

	// web server
	if false {
		mux := http.NewServeMux()
		mux.Handle("/depth", websocket.BroadcastFrames(edge.Depth, astra.Downsample(4)))
		server := &http.Server{Addr: ":9091", Handler: mux}

		graceful, err := httpConf.ListenAndServe(server)
		if err != nil {
			return
		}
		defer graceful.Stop()
	}

	// preview
	if true {
		preview.Default(func(p *preview.Preview) {
			go p.Show(edge.Depth)
			defer p.Stop()

			ticker := time.NewTicker(1 * time.Second)

			alive := true
			for alive {
				select {
				case <-ticker.C:
					alive = p.State == preview.Active

				case <-done:
					alive = false
				}
			}
		})
	}

}
