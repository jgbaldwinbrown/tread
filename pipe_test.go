package tread

import (
	"reflect"
	"testing"
)

func TestPipe(t *testing.T) {
	in := []float64{3.5, 4, 5, 6}
	ir := MakeReader[float64](in)

	var b SliceBuffer[float64]

	pr, pw := Pipe[float64]()
	done := make(chan struct{}, 2)
	go func() {
		Copy[float64](pw, &ir)
		pw.Close()
		done <- struct{}{}
	}()

	go func() {
		Copy[float64](&b, pr)
		pr.Close()
		done <- struct{}{}
	}()

	<-done
	<-done

	out := b.Slice()
	if !reflect.DeepEqual(out, in) {
		t.Errorf("out %v not equal to in %v", out, in)
	}
}
