package websocket

import (
// "log"
)

type FrameProcessor func(*Frame)

func Downsample(f int) FrameProcessor {
	return func(frame *Frame) {
		width := frame.Metadata.Width / f
		height := frame.Metadata.Height / f
		data := make([]int, width*height)

		i := 0
		for r := 0; r <= height; r += f {
			for c := 0; c <= width; c += f {
				data[i] = int(frame.Data[r*frame.Metadata.Width+c])
				i += 1
			}
		}

		frame.Metadata.Width = width
		frame.Metadata.Height = height
		frame.Data = data
	}

}
