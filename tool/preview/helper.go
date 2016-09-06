package preview

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

import (
	"github.com/mikydna/astra"
)

func LoadJSON(glob string) ([][]int, error) {
	files, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}

	frames := [][]int{}
	for _, file := range files {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		frame := astra.CameraDepthFrame{}
		if err := json.Unmarshal(b, &frame); err != nil {
			return nil, err
		}

		data := make([]int, len(frame.Buffer))
		for j, val := range frame.Buffer {
			data[j] = int(val)
		}

		frames = append(frames, data)
	}

	return frames, nil
}
