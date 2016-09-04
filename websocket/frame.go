package websocket

import (
	"errors"
	"time"
)

import (
	"github.com/mikydna/astra"
)

var (
	ErrInvalidCaptureFrame = errors.New("Invalid capture frame")
)

type FrameMetadata struct {
	Timestamp int64 `json:"t"`
	Index     int   `json:"i"`
	Width     int   `json:"w"`
	Height    int   `json:"h"`
}

type Frame struct {
	Metadata FrameMetadata `json:"metadata"`
	Data     []int         `json:"data"`
}

func FromDepthFrame(c astra.CameraDepthFrame) (*Frame, error) {
	if (c.Width * c.Height) != len(c.Buffer) {
		return nil, ErrInvalidCaptureFrame
	}

	data := make([]int, len(c.Buffer))
	for i, val := range c.Buffer {
		data[i] = int(val)
	}

	frame := Frame{
		Metadata: FrameMetadata{
			Timestamp: time.Now().Unix(),
			Index:     c.Index,
			Width:     c.Width,
			Height:    c.Height,
		},
		Data: data,
	}

	return &frame, nil
}
