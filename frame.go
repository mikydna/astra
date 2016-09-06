package astra

import (
	"errors"
	"time"
)

var (
	ErrInvalidCaptureFrame = errors.New("Invalid capture frame")
)

type Frame struct {
	Timestamp int64 `json:"t"`
	Index     int   `json:"i"`
	Width     int   `json:"w"`
	Height    int   `json:"h"`
	Data      []int `json:"data"`
}

func FromDepthFrame(c CameraDepthFrame) (*Frame, error) {
	if (c.Width * c.Height) != len(c.Buffer) {
		return nil, ErrInvalidCaptureFrame
	}

	data := make([]int, len(c.Buffer))
	for i, val := range c.Buffer {
		data[i] = int(val)
	}

	frame := Frame{
		Index:     c.Index,
		Timestamp: time.Now().Unix(),
		Width:     c.Width,
		Height:    c.Height,
		Data:      data,
	}

	return &frame, nil
}

type FrameProcessor func(*Frame)

func Downsample(f int) FrameProcessor {
	return func(frame *Frame) {
		width := frame.Width / f
		height := frame.Height / f
		data := make([]int, width*height)

		for r := 0; r < height; r++ {
			for c := 0; c < width; c++ {
				data[r*width+c] = int(frame.Data[(r*f*frame.Width)+(c*f)])
			}
		}

		frame.Width = width
		frame.Height = height
		frame.Data = data
	}

}
