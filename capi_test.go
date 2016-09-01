package astra

import (
	"testing"
)

func TestAstra(t *testing.T) {

	if rc := Initialize(); rc != StatusSuccess {
		t.Fatalf("Astra should return SUCCESS on initialize, heard: %s", rc.String())
	}

	if rc := Terminate(); rc != StatusSuccess {
		t.Fatalf("Astra should return SUCCESS on terminate, heard: %s", rc.String())
	}

}
