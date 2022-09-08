package tread

import (
	"io"
)

type Reader[T any] interface {
	Read(p []T) (n int, err error)
}

type Writer[T any] interface {
	Write(p []T) (n int, err error)
}

func Copy[T any](dest Writer[T], src Reader[T]) (n int, err error) {
	buf := make([]T, 4096)
	for err == nil {
		var readn int
		readn, err = src.Read(buf)
		dest.Write(buf[:readn])
		n += readn
	}
	if err == io.EOF {
		err = nil
	}
	return
}

func Transform[T any, U any](dest Writer[T], src Reader[U], tr func(U) T) (n int, err error) {
	tread := WrapReader(src, tr)
	return Copy(dest, Reader[T](&tread))
}

type WrappedReader[T, U any] struct {
	r Reader[T]
	tbuf []T
	f func(T) U
}

func WrapReader[T, U any](r Reader[T], f func(T) U) WrappedReader[T,U] {
	return WrappedReader[T,U]{r, []T{}, f}
}

func resize[T any](buf []T, n int) []T {
	if cap(buf) < n {
		newbuf := make([]T, len(buf), n)
		copy(newbuf, buf)
		buf = newbuf
	}
	return buf[:n]
}

func (r *WrappedReader[T,U]) Read(p []U) (n int, err error) {
	r.tbuf = resize(r.tbuf, len(p))
	n, err = r.r.Read(r.tbuf)
	for i, _ := range r.tbuf[:n] {
		p[i] = r.f(r.tbuf[i])
	}
	return n, err
}

type WrappedWriter[T, U any] struct {
	w Writer[U]
	ubuf []U
	f func(T) U
}

func WrapWriter[T, U any](w Writer[U], f func(T) U) WrappedWriter[T, U] {
	return WrappedWriter[T,U]{w, []U{}, f}
}

func (w *WrappedWriter[T, U]) Write(p []T) (n int, err error) {
	w.ubuf = resize(w.ubuf, len(p))
	for i, t := range p {
		w.ubuf[i] = w.f(t)
	}
	return w.w.Write(w.ubuf)
}

type SliceReader[T any] struct {
	slice []T
	idx int
}

func MakeReader[T any](s []T) SliceReader[T] {
	return SliceReader[T]{
		slice: s,
		idx: 0,
	}
}

func (s *SliceReader[T]) Read(p []T) (n int, err error) {
	n = len(p)
	remaining := len(s.slice) - s.idx
	if remaining <= n {
		err = io.EOF
		n = remaining
		if n < 0 {
			n = 0
		}
	}
	copy(p, s.slice[s.idx:s.idx+n])
	s.idx += n
	return
}

type SliceBuffer[T any] struct {
	slice []T
	idx int
}

func (s *SliceBuffer[T]) Read(p []T) (n int, err error) {
	n = len(p)
	remaining := len(s.slice) - s.idx
	if remaining <= n {
		err = io.EOF
		n = remaining
		if n < 0 {
			n = 0
		}
	}
	copy(p, s.slice[s.idx:s.idx+n])
	s.idx += n
	return
}

func (s *SliceBuffer[T]) Write(p []T) (n int, err error) {
	s.slice = append(s.slice, p...)
	return len(p), nil
}

func (s *SliceBuffer[T]) Slice() []T {
	return s.slice
}

type Iterator[T any] struct {
	r Reader[T]
	buf []T
}

func NewIterator[T any](r Reader[T]) *Iterator[T] {
	i := new(Iterator[T])
	i.buf = make([]T, 1)
	i.r = r
	return i
}

func (i *Iterator[T]) Next() bool {
	n, err := i.r.Read(i.buf)
	return !(n <= 0 && err != nil)
}

func (i *Iterator[T]) Value() T {
	return i.buf[0]
}
