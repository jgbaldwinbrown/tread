package tread

import (
	"testing"
	"strings"
	"io"
)

func TestExec(t *testing.T) {
	expect := "one two\n"
	r := strings.NewReader(expect)
	var b strings.Builder

	r2 := MustExec(r, "cat")
	io.Copy(&b, r2)
	out := b.String()
	if expect != out {
		t.Errorf("expect %v != out %v", expect, out)
	}
}
