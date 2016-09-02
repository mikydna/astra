package astra

import (
	"errors"
	"log"
)

var (
	ErrNoCameraInUse             = errors.New("No camera in use")
	ErrDepthStreamAlreadyStarted = errors.New("")
)

type CameraDepthFrame struct {
	Index         int
	Width, Height uint
	Buffer        []int16
}

type CameraDepthStream struct {
	stream *DepthStream
	out    chan CameraDepthFrame
}

func AcquireCameraDepthStream(c *Camera) (*CameraDepthStream, error) {
	if c.conn == nil || c.reader == nil {
		return nil, ErrNoCameraInUse
	}

	newDepthStream := new(DepthStream)

	if rc := GetDepthStream(*c.reader, newDepthStream); rc != StatusSuccess {
		return nil, rc.Error()
	}

	if rc := StartDepthStream(*newDepthStream); rc != StatusSuccess {
		return nil, rc.Error()
	}

	cameraDepthStream := &CameraDepthStream{
		stream: newDepthStream,
		out:    make(chan CameraDepthFrame, 1),
	}

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

	index, rc := GetDepthFrame(frame, newDepthFrame)
	if rc == StatusSuccess {

		if width, height, buffer, err := processDepthFrame(*newDepthFrame); err != nil {
			log.Println(err) // fix

		} else {
			ds.out <- CameraDepthFrame{index, width, height, buffer}

		}

	} else {
		log.Printf("Skipping Frame: index=%d reasons=%s", index, rc.String())

	}
}

func (ds *CameraDepthStream) Frames() <-chan CameraDepthFrame {
	return ds.out
}

func processDepthFrame(frame DepthFrame) (uint, uint, []int16, error) {
	width, height, rc := GetDepthFrameMetadata(frame)
	if rc != StatusSuccess {
		return 0, 0, nil, rc.Error()
	}

	buffer, rc := GetDepthFrameBuffer(frame)
	if rc != StatusSuccess {
		return 0, 0, nil, rc.Error()
	}

	return width, height, buffer, nil
}
