package astra

import (
	"log"
	"time"
)

var (
	DefaultStreamConf = CameraStreamConf{
		1000 * time.Millisecond,
		// ^^^ fix: initialization crutch
		50 * time.Millisecond,
	}
)

type FrameHandler interface {
	Handle(frame ReaderFrame)
}

type CameraStreamConf struct {
	delay    time.Duration
	interval time.Duration
}

type CameraState uint8

const (
	Initialized CameraState = iota
	Active
	Streaming
	Stopped
)

type Camera struct {
	state    CameraState
	addr     string
	conn     *StreamSetConnection
	reader   *Reader
	handlers []FrameHandler
	frames   chan ReaderFrame
	done     chan bool
}

func NewCamera() (*Camera, error) {
	if rc := Initialize(); rc != StatusSuccess {
		return nil, rc.Error()
	}

	return &Camera{
		state:    Initialized,
		conn:     nil,
		reader:   nil,
		frames:   make(chan ReaderFrame),
		handlers: []FrameHandler{},
		done:     make(chan bool),
	}, nil
}

func (c *Camera) Use(deviceAddr string) error {

	c.conn = new(StreamSetConnection)
	if rc := OpenStream(deviceAddr, c.conn); rc != StatusSuccess {
		c.conn = nil
		return rc.Error()
	}

	c.reader = new(Reader)
	if rc := CreateReader(*c.conn, c.reader); rc != StatusSuccess {
		c.conn = nil
		c.reader = nil
		return rc.Error()
	}

	c.state = Active

	return nil
}

func (c *Camera) HandleFrame(h FrameHandler) {
	c.handlers = append(c.handlers, h)
}

func (c *Camera) PollStream(conf CameraStreamConf) {
	// todo: this can only be called once?
	time.Sleep(conf.delay)

	alive := true
	ticker := time.NewTicker(conf.interval)
	c.state = Streaming

	for alive {

		select {
		case <-ticker.C:

			if rc := Update(); rc == StatusSuccess {
				newFrame := new(ReaderFrame)
				rc := OpenReaderFrame(*c.reader, newFrame)

				if rc == StatusSuccess {
					for _, handler := range c.handlers {
						handler.Handle(*newFrame)
					}

					CloseReaderFrame(newFrame)
				}

			} else {
				log.Println("Update failed? ", rc)

			}

		case <-c.done:
			ticker.Stop()
			alive = false
			c.state = Active
		}
	}
}

func (c *Camera) Stop() error {

	// has to block; poll will most likely be executed by a goroutine
	// must stop stream thread before destroying readers and conn
	// fix: consider tracking state (or use a waitgroup?)
	// - if not, camera terminate/stop can panic
	if c.state == Streaming {
		c.done <- true
	}

	if rc := DestroyReader(c.reader); rc != StatusSuccess {
		c.reader = nil
		return rc.Error()
	}

	if rc := CloseStream(c.conn); rc != StatusSuccess {
		c.conn = nil
		return rc.Error()
	}

	if rc := Terminate(); rc != StatusSuccess {
		return rc.Error()
	}

	c.state = Stopped

	return nil
}
