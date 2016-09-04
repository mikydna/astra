package websocket

type FrameProcessor func(*Frame)

func Downsample(f int) FrameProcessor {
	return func(frame *Frame) {
		width := frame.Width / f
		height := frame.Height / f
		data := make([]int, width*height)

		i := 0
		for r := 0; r <= height; r += f {
			for c := 0; c <= width; c += f {
				data[i] = int(frame.Data[r*frame.Width+c])
				i += 1
			}
		}

		frame.Width = width
		frame.Height = height
		frame.Data = data
	}

}
