package preview

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

import (
	"github.com/mikydna/astra"
)

func TestPreview(t *testing.T) {
	prerecorded, _ := LoadJSON("/Users/andy/Desktop/capture/astra-*.json")
	t.Log(len(prerecorded))

	Launch(Conf{1, 640, 320}, prerecorded)
}

func TestMain(m *testing.M) {

	go func() {
		prerecorded, _ := LoadJSON("/Users/andy/Desktop/capture/astra-*.json")
		// m.Log(len(prerecorded))

		Launch(Conf{1, 640, 320}, prerecorded)

		m.Run()
		<-time.After(15 * time.Second)
		os.Exit(0)

	}()

}
