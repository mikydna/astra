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

	depth, err := NewCameraDepthStream(camera)
	if err != nil {
		t.Fatal(err)
	}

	hfov, vfov, err := depth.GetFOV()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Depth Stream FOV: h=%f v=%f", hfov, vfov)

	camera.StartStream(DefaultStreamConf)

}
