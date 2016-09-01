package astra

import (
	"errors"
	"log"
)

var (
	ErrCameraClosed              = errors.New("Camera must be openned first")
	ErrDepthStreamAlreadyStarted = errors.New("")
)

type CameraDepthStream struct {
	stream *DepthStream
}

func AcquireCameraDepthStream(c *Camera) (*CameraDepthStream, error) {
	if c.conn == nil || c.reader == nil {
		return nil, ErrCameraClosed
	}

	newDepthStream := new(DepthStream)

	if rc := GetDepthStream(*c.reader, newDepthStream); rc != StatusSuccess {
		return nil, rc.Error()
	}

	if rc := StartDepthStream(*newDepthStream); rc != StatusSuccess {
		return nil, rc.Error()
	}

	cameraDepthStream := &CameraDepthStream{newDepthStream}

	c.HandleFrame(cameraDepthStream) // weird

	return cameraDepthStream, nil
}

func (ds *CameraDepthStream) GetFOV() (float32, float32, error) {
	hfov, vfov, rc := GetDepthStreamFOV(*ds.stream)
	if rc != StatusSuccess {
		return -1, -1, rc.Error()
	}

	return hfov, vfov, nil // radians
}

func (ds *CameraDepthStream) Handle(frame ReaderFrame) {
	newDepthFrame := new(DepthFrame)

	frameIndex, rc := GetDepthFrame(frame, newDepthFrame)
	if rc == StatusSuccess {
		log.Printf("Process Depth Frame: index=%d", frameIndex)

	} else {
		log.Printf("Skipping Frame: index=%d reasons=%s", frameIndex, rc.String())

	}
}
