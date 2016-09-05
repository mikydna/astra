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

	// tool/preview stub
	preview.Launch()

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

	mux.Handle("/depth", websocket.BroadcastFrames(edge.Depth, websocket.Downsample(8)))

	graceful, err := httpConf.ListenAndServe(server)
	if err != nil {
		return
	}
	defer graceful.Stop()

	// start
	go edge.Start()

	<-done
}
