package tread

import (
	"testing"
	"reflect"
)

func TestBufferedReader(t *testing.T) {
	r := NewBufferedReader[int](NewRanger(Range[int]{0, 10000, 1}))
	expect := make([]int, 10000)
	for i, _ := range expect {
		expect[i] = i
	}

	var out SliceBuffer[int]

	Copy[int](&out, r)

	if !reflect.DeepEqual(out.Slice(), expect) {
		t.Errorf("out.Slice() %v != expect %v", out.Slice(), expect);
	}
}
