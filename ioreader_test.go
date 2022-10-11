package tread

import (
	"testing"
	"io"
	"bytes"
)

func TestIoReadingOfReader(t *testing.T) {
	b := []byte{4,5,6,7}
	sr := MakeReader[byte](b)
	var _ io.Reader = &sr
}

func TestIoReader(t *testing.T) {
	b := []byte{4,5,6,7}
	sr := bytes.NewReader(b)
	var _ Reader[byte] = sr
}
