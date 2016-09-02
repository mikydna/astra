package edge

import (
	"log"
)

import (
	"github.com/mikydna/astra"
)

type Edge struct {
	cameras map[string]*astra.Camera
	done    chan bool

	Depth chan astra.CameraDepthFrame
}

func NewEdge(addrs []string) (*Edge, error) {
	cameras := make(map[string]*astra.Camera)
	for _, addr := range addrs {
		camera, err := astra.NewCamera()
		if err != nil {
			return nil, err
		}

		if err := camera.Use(addr); err != nil {
			return nil, err
		}

		cameras[addr] = camera
	}

	edge := &Edge{
		cameras: cameras,
		done:    make(chan bool),

		Depth: make(chan astra.CameraDepthFrame),
	}

	return edge, nil
}

func (e *Edge) Start() error {
	inbound := []<-chan astra.CameraDepthFrame{}

	for addr, camera := range e.cameras {
		log.Printf("Starting depth-stream for @%s", addr)

		perCameraDepthStream, err := astra.AcquireCameraDepthStream(camera)
		if err != nil {
			return err
		}

		// pick off frame chan
		inbound = append(inbound, perCameraDepthStream.Frames())

		// trigger camera poll
		go camera.PollStream(astra.DefaultStreamConf)

	}

	// single control loop for all streams
	alive := true
	for alive {
		select {
		case <-e.done:
			alive = false
		default:
			for _, stream := range inbound {
				select {
				case frame := <-stream:
					e.Depth <- frame
				default:
					// nothing heard
				}
			}
		}
	}

	return nil
}

func (e *Edge) Shutdown() {

	e.done <- true

	for addr, camera := range e.cameras {
		if err := camera.Stop(); err != nil {
			log.Printf("Could not showdown: @%s %v", addr, err)
		} else {
			log.Printf("Shutting down camera @%s", addr)
		}
	}
}

func (e *Edge) Found() []string {
	addrs := []string{}
	for addr, _ := range e.cameras {
		addrs = append(addrs, addr)
	}

	return addrs
}
