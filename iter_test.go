package tread

import (
	"testing"
	"reflect"
	"io"
)

func TestIterator(t *testing.T) {
	in := []float64{3,4,5}
	out := []float64{}
	r := MakeReader(in)
	iter := NewIterator[float64](&r)
	for iter.Next() {
		out = append(out, iter.Value())
	}
	if !reflect.DeepEqual(in, out) {
		t.Errorf("in %v != out %v", in, out)
	}
}

func TestIterRead(t *testing.T) {
	in := []float64{3,4,5}
	out := make([]float64, 4)
	r := MakeReader(in)
	iter := NewIterator[float64](&r)
	n, err := iter.Read(out)
	if !reflect.DeepEqual(in, out[:3]) || n != 3 || err != io.EOF {
		t.Errorf("in %v != out %v", in, out)
	}
}
