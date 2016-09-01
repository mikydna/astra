package astra

import (
	"testing"
)

func TestCameraDepth(t *testing.T) {
	camera, err := NewCamera()
	if err != nil {
		t.Fatal(err)
	}
	defer camera.Terminate()

	if err := camera.Use("device/default"); err != nil {
		t.Fatal(err)
	}

	depth, err := camera.DepthStream()
	if err != nil {
		t.Fatal(err)
	}

	hfov, vfov, err := depth.GetFOV()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("DepthStream FOV: h=%f v=%f", hfov, vfov)

	camera.StartStream(DefaultStreamConf)

}
