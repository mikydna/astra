package websocket

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
