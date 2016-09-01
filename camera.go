package astra

import (
	"errors"
	"log"
	"time"
)

var (
	ErrCameraClosed              = errors.New("Camera must be openned first")
	ErrDepthStreamAlreadyStarted = errors.New("")
)

var (
	DefaultStreamConf = CameraStreamConf{
		1000 * time.Millisecond,
		100 * time.Millisecond,
	}
)

type CameraStreamConf struct {
	delay time.Duration
	sleep time.Duration
}

type Camera struct {
	addr   string
	conn   *StreamSetConnection
	reader *Reader
	frames chan bool
}

func NewCamera() (*Camera, error) {
	if rc := Initialize(); rc != StatusSuccess {
		return nil, rc.Error()
	}

	return &Camera{
		conn:   new(StreamSetConnection),
		reader: new(Reader),
		frames: make(chan bool),
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

func (c *Camera) DepthStream() (*CameraDepthStream, error) {
	if c.conn == nil || c.reader == nil {
		return nil, ErrCameraClosed
	}

	newDepthStream := new(DepthStream)
	if rc := GetDepthStream(*c.reader, newDepthStream); rc != StatusSuccess {
		return nil, rc.Error()
	}

	return &CameraDepthStream{newDepthStream}, nil
}

func (c *Camera) StartStream(conf CameraStreamConf) {
	time.Sleep(conf.delay)

	for i := 0; i < 10; i++ {
		time.Sleep(conf.sleep)

		Update()

		newFrame := new(ReaderFrame)
		rc := OpenReaderFrame(*c.reader, newFrame)
		defer CloseReaderFrame(newFrame)

		log.Println("Frame ", i, rc.String())
	}

}

func (c *Camera) Terminate() error {
	return Terminate().Error()
}
