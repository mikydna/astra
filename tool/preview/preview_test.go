package preview

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

import (
	"github.com/mikydna/astra"
)

func init() {
	runtime.LockOSThread()
}

func fromDir(glob string) ([][]int, error) {
	files, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}

	log.Println(files)

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

func TestMain(m *testing.M) {
	frames, err := fromDir("/Users/andy/Desktop/capture/astra-*.json")
	if err != nil {
		return
	}

	go func() {
		<-time.After(15 * time.Second)
		os.Exit(m.Run())
	}()

	Launch(Conf{2, 1, 640, 480}, frames)
}
