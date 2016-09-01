package astra

import (
	"log"
	"time"
)

var (
	DefaultStreamConf = CameraStreamConf{
		1000 * time.Millisecond,
		100 * time.Millisecond,
	}
)

type FrameHandler interface {
	Handle(frame ReaderFrame)
}

type CameraStreamConf struct {
	delay time.Duration
	sleep time.Duration
}

type Camera struct {
	addr     string
	conn     *StreamSetConnection
	reader   *Reader
	frames   chan ReaderFrame
	handlers []FrameHandler
}

func NewCamera() (*Camera, error) {
	if rc := Initialize(); rc != StatusSuccess {
		return nil, rc.Error()
	}

	return &Camera{
		conn:     new(StreamSetConnection),
		reader:   new(Reader),
		frames:   make(chan ReaderFrame),
		handlers: []FrameHandler{},
	}, nil
}

func (c *Camera) Use(deviceAddr string) error {
	if rc := OpenStream(deviceAddr, c.conn); rc != StatusSuccess {
		c.conn = nil
		return rc.Error()
	}

	if rc := CreateReader(*c.conn, c.reader); rc != StatusSuccess {
		c.conn = nil
		c.reader = nil
		return rc.Error()
	}

	return nil
}

func (c *Camera) HandleFrame(h FrameHandler) {
	c.handlers = append(c.handlers, h)
}

func (c *Camera) StartStream(conf CameraStreamConf) {
	time.Sleep(conf.delay)

	for i := 0; i < 10; i++ {
		time.Sleep(conf.sleep)

		Update()

		newFrame := new(ReaderFrame)
		rc := OpenReaderFrame(*c.reader, newFrame)
		if rc != StatusSuccess {
			continue
		}

		log.Println("Frame ", i, rc.String())
		for _, handler := range c.handlers {
			handler.Handle(*newFrame)
		}

		CloseReaderFrame(newFrame)
	}

}

func (c *Camera) Terminate() error {
	return Terminate().Error()
}
