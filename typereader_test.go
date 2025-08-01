package tread

import (
	"reflect"
	"testing"
	"fmt"
)

func TestCopy(t *testing.T) {
	in := []float64{3.5, 4, 5, 6}

	r := MakeReader(in)
	var buf SliceBuffer[float64]
	n, err := Copy[float64](&buf, &r)
	if err != nil {
		t.Error(err)
	}
	if n != 4 {
		t.Errorf("Copied %v elements instead of 4", n)
	}
	out := buf.Slice()
	if !reflect.DeepEqual(out, in) {
		t.Errorf("out %v not equal to in %v", out, in)
	}
}

func TestTransform(t *testing.T) {
	in := []float64{3.5, 4, 5, 6}
	expect := []string{"3.5", "4", "5", "6"}

	r := MakeReader(in)
	var buf SliceBuffer[string]
	f := func(n float64) string { return fmt.Sprint(n) }
	n, err := Transform[string, float64](&buf, &r, f)
	if err != nil {
		t.Error(err)
	}
	if n != 4 {
		t.Errorf("Transformed %v elements instead of 4", n)
	}
	out := buf.Slice()
	if !reflect.DeepEqual(out, expect) {
		t.Errorf("out %v not equal to expect %v", out, expect)
	}
}
