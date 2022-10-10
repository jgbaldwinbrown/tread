package tread

import (
	"testing"
	"reflect"
)

func TestIterator(t *testing.T) {
	in := []float64{3,4,5}
	out := []float64{}
	r := MakeReader(in)
	for f, ok := Next[float64](&r); ok; f, ok = Next[float64](&r) {
		out = append(out, f)
	}
	if !reflect.DeepEqual(in, out) {
		t.Errorf("in %v != out %v", in, out)
	}
}
