package tread

import (
	"golang.org/x/exp/constraints"
	"io"
)

type number interface {
	constraints.Float | constraints.Integer
}

type Range[T number] struct {
	Start T
	End T
	Step T
}

type Ranger[T number] struct {
	Range[T]
	Idx int
}

func (r Range[T]) Get(i int) T {
	return r.Start + (T(i) * r.Step)
}

func (r Range[T]) Contains(t T) bool {
	return r.Start <= t && r.End > t
}

func MakeRanger[T number](in Range[T]) Ranger[T] {
	return Ranger[T] { in, -1 }
}

func NewRanger[T number](in Range[T]) *Ranger[T] {
	r := MakeRanger(in)
	return &r
}

func (r *Ranger[T]) Next() (T, bool) {
	r.Idx++
	out := r.Get(r.Idx)
	return out, r.Contains(out)
}

func (r *Ranger[T]) Read(p []T) (n int, err error) {
	for i, _ := range p {
		val, ok := r.Next()
		if !ok {
			err = io.EOF
			break
		}
		p[i] = val
		n++
	}
	return n, err
}
