package stations

import (
	"bytes"
	"testing"
)

func TestInit(t *testing.T) {
	hash := New()
	hash.Set([]byte("test"), 1, 1, 1, 1)
	res, err := hash.Get([]byte("test"))

	if err != nil {
		t.Error("Unexpected error: ", err)
	}
	if !bytes.Equal(res.Name, []byte("test")) {
		t.Error("Name should be 'test': ", string(res.Name))
	}
}
