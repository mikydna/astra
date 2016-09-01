package astra

import (
	"errors"
)

var (
	ErrCameraClosed              = errors.New("Camera must be openned first")
	ErrDepthStreamAlreadyStarted = errors.New("")
)

type CameraDepthStream struct {
	stream *DepthStream
}

func (ds *CameraDepthStream) GetFOV() (float32, float32, error) {
	hfov, vfov, rc := GetDepthStreamFOV(*ds.stream)
	if rc != StatusSuccess {
		return -1, -1, rc.Error()
	}

	return hfov, vfov, nil
}

type Camera struct {
	addr   string
	conn   *StreamSetConnection
	reader *Reader
}

func NewCamera() (*Camera, error) {
	if rc := Initialize(); rc != StatusSuccess {
		return nil, rc.Error()
	}

	return &Camera{
		conn:   new(StreamSetConnection),
		reader: new(Reader),
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

func (c *Camera) StartDepthStream() (*CameraDepthStream, error) {
	if c.conn == nil || c.reader == nil {
		return nil, ErrCameraClosed
	}

	newDepthStream := new(DepthStream)
	if rc := GetDepthStream(*c.reader, newDepthStream); rc != StatusSuccess {
		return nil, rc.Error()
	}

	return &CameraDepthStream{newDepthStream}, nil
}

func (c *Camera) Terminate() error {
	return Terminate().Error()
}
