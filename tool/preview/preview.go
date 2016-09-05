package preview

import (
	// "fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"time"
)

import (
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
)

func Launch() {
	width := 640
	height := 480

	driver.Main(func(s screen.Screen) {
		size := image.Point{width, height}

		win, err := s.NewWindow(&screen.NewWindowOptions{width, height})
		if err != nil {
			log.Fatal(err)
		}
		defer win.Release()

		buf, err := s.NewBuffer(size)
		if err != nil {
			log.Fatal(err)
		}
		defer buf.Release()

		// tex, err := s.NewTexture(size)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer tex.Release()

		go func() {
			ticker := time.NewTicker(100 * time.Millisecond)

			for i := 0; true; i++ {
				select {
				case <-ticker.C:
					drawDepth(buf.RGBA(), []int{})
					win.Upload(image.Point{}, buf, buf.Bounds())
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

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			m.SetRGBA(x, y, color.RGBA{
				uint8((rand.Float32() * 0xff)),
				uint8((rand.Float32() * 0xff)),
				uint8((rand.Float32() * 0xff)),
				0xff})
		}
	}

}
