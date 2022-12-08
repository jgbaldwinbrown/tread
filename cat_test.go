package tread

import (
	"io"
	"testing"
	"strings"
)

func TestMultiReader(t *testing.T) {
	in1 := strings.NewReader("apple")
	in2 := strings.NewReader("banana")
	expect := "applebanana"
	var b strings.Builder
	m := MultiReader[byte](in1, in2)
	io.Copy(&b, m)
	out := b.String()
	if expect != out {
		t.Errorf("expect %v != out %v", expect, out)
	}
}
