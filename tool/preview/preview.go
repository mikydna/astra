package preview

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"time"
)

import (
	"github.com/mikydna/astra"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	// "golang.org/x/mobile/event/paint"
)

const (
	FPS_1  = 1000 * time.Millisecond
	FPS_10 = 100 * time.Millisecond
	FPS_20 = 50 * time.Millisecond
	FPS_30 = 33 * time.Millisecond
	FPS_60 = 16 * time.Millisecond
)

type State uint8

const (
	NotStarted State = iota
	Inactive
	Active
)

type PreviewConf struct {
	Scale         int
	Width, Height int
}

type Preview struct {
	State  State
	Screen screen.Screen
	done   chan bool
}

func NewPreview(scr screen.Screen) *Preview {
	return &Preview{
		State:  NotStarted,
		Screen: scr,
		done:   make(chan bool),
	}
}

func (p *Preview) Stop() {
	if p.State == Active {
		p.State = Inactive
		p.done <- true
	}
}

func (p *Preview) Show(in <-chan astra.CameraDepthFrame) {
	p.State = Active

	win, err := p.Screen.NewWindow(&screen.NewWindowOptions{1280, 960})
	if err != nil {
		log.Fatal(err)
	}
	defer win.Release()

	// ****************** //

	go func() {
		for frame := range in {
			size := image.Point{frame.Width, frame.Height}

			buf, err := p.Screen.NewBuffer(size)
			if err != nil {
				log.Fatal(err)
			}
			defer buf.Release()

			tex, err := p.Screen.NewTexture(size)
			if err != nil {
				log.Fatal(err)
			}
			defer tex.Release()

			drawDepth(buf.RGBA(), frame.Buffer)

			tex.Upload(image.Point{}, buf, buf.Bounds())

			win.Scale(
				image.Rectangle{
					image.Point{0, 0},
					image.Point{1280, 960},
				},
				tex, tex.Bounds(), screen.Src, nil)

			win.Publish()
		}
	}()

	// ****************** //

	events := make(chan interface{})
	defer close(events)
	go func() {
		for {
			events <- win.NextEvent()
		}
	}()

	alive := true
	for alive {
		select {
		case e := <-events:
			switch e := e.(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					p.Stop()
				}
			case key.Event:
				if e.Code == key.CodeEscape {
					p.Stop()
				}
			}

		case <-p.done:
			alive = false
		}

	}
}

func Default(f func(preview *Preview)) {
	driver.Main(func(scr screen.Screen) {
		f(NewPreview(scr))
	})
}

func drawDepth(m *image.RGBA, frame []int16) {
	bounds := m.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			x := bounds.Min.X + c
			y := bounds.Min.Y + r
			val := (float32(frame[r*width+c]) / float32(10000))

			m.SetRGBA(x, y, color.RGBA{
				uint8(val * 0xff),
				0x00,
				0x00,
				0xff,
			})
		}
	}
}
