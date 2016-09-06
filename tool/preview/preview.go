package preview

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	// "math/rand"
	"time"
)

import (
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
)

const (
	FPS_1  = 1000 * time.Millisecond
	FPS_10 = 100 * time.Millisecond
	FPS_20 = 50 * time.Millisecond
	FPS_30 = 33 * time.Millisecond
	FPS_60 = 16 * time.Millisecond
)

type Conf struct {
	Mult, Scale   int
	Width, Height int
}

func Launch(conf Conf, frames [][]int) {
	mult := conf.Mult
	scale := conf.Scale
	width := conf.Width / scale
	height := conf.Height / scale

	driver.Main(func(s screen.Screen) {
		size := image.Point{width, height}

		win, err := s.NewWindow(&screen.NewWindowOptions{mult * scale * width, mult * scale * height})
		if err != nil {
			log.Fatal(err)
		}
		defer win.Release()

		buf, err := s.NewBuffer(size)
		if err != nil {
			log.Fatal(err)
		}
		defer buf.Release()

		tex, err := s.NewTexture(size)
		if err != nil {
			log.Fatal(err)
		}
		defer tex.Release()

		go func() {

			ticker := time.NewTicker(FPS_10)

			for i := 0; true; i++ {
				select {
				case <-ticker.C:
					mat := frames[i%len(frames)]

					drawDepth(buf.RGBA(), mat)
					tex.Upload(image.Point{}, buf, buf.Bounds())

					win.Scale(
						image.Rectangle{
							image.Point{0, 0},
							image.Point{mult * scale * width, mult * scale * height},
						},
						tex, tex.Bounds(), screen.Src, nil)

					win.Publish()
				}
			}
		}()

		for {
			e := win.NextEvent()

			switch e := e.(type) {

			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}

			case key.Event:
				if e.Code == key.CodeEscape {
					return
				}

			case paint.Event:
				log.Println("paint")
			}

		}

	})

}

func drawDepth(m *image.RGBA, mat []int) {
	bounds := m.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			x := bounds.Min.X + c
			y := bounds.Min.Y + r
			val := float32(mat[r*width+c]) / float32(10000)

			m.SetRGBA(x, y, color.RGBA{
				uint8(val * 0xff),
				0x00,
				0x00, // uint8((rand.Float32() * 0xff)),
				0xff,
			})

		}
	}

}
