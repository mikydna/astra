package astra

type CameraDepthStream struct {
	stream *DepthStream
}

func (ds *CameraDepthStream) GetFOV() (float32, float32, error) {
	hfov, vfov, rc := GetDepthStreamFOV(*ds.stream)
	if rc != StatusSuccess {
		return -1, -1, rc.Error()
	}

	return hfov, vfov, nil // radians
}

func (ds *CameraDepthStream) Start() error {
	if rc := StartDepthStream(*ds.stream); rc != StatusSuccess {
		return rc.Error()
	}

	return nil
}
