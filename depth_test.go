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

	t.Logf("Depth Stream FOV: h=%f v=%f", hfov, vfov)

	go camera.StartStream(DefaultStreamConf)

	<-time.After(5 * time.Second)

}
