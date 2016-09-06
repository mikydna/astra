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

	// web server
	mux := http.NewServeMux()

	mux.Handle("/depth", websocket.BroadcastFrames(edge.Depth, astra.Downsample(4)))

	server := &http.Server{Addr: ":9091", Handler: mux}

	graceful, err := httpConf.ListenAndServe(server)
	if err != nil {
		return
	}
	defer graceful.Stop()

	// start
	go edge.Start()

	// tools
	prerecorded, _ := preview.LoadJSON("/Users/andy/Desktop/capture/astra-*.json")
	preview.Launch(preview.Conf{1, 640, 480}, prerecorded)

	<-done
}
