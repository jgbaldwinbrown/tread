package tread

import (
	"testing"
	"reflect"
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
