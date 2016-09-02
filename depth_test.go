package astra

import (
	"testing"
	"time"
)

func TestCameraDepth(t *testing.T) {
	camera, err := NewCamera()
	if err != nil {
		t.Fatal(err)
	}
	defer camera.Stop()

	if err := camera.Use("device/default"); err != nil {
		t.Fatal(err)
	}

	depth, err := AcquireCameraDepthStream(camera)
	if err != nil {
		t.Fatal(err)
	}

	hfov, vfov, err := depth.GetFOV()
	if err != nil {
		t.Fatal(err)
	}

	if hfov < 0 && hfov >= 2 {
		t.Errorf("Unexpected hfov value: hfov=%f", hfov)
	}

	if vfov < 0 && vfov >= 2 {
		t.Errorf("Unexpected vfov value: vfov=%f", vfov)
	}

	go camera.PollStream(DefaultStreamConf)

	timeout := time.After(10 * time.Second)
	alive := true
	heard := 0
	for alive {
		select {
		case <-depth.Frames():
			heard += 1
			alive = heard < 10

		case <-timeout:
			alive = false
			t.Errorf("Did not hear at least 10 events: heard %d", heard)
			t.Fail()
		}
	}

}
