package tread

import (
	"io"
)

const bufsiz = 4096

type Iter[T any] interface {
	Next() (T, bool)
}

type BufferedReader[T any] struct {
	r Reader[T]
	buf []T
	idx int
	closed bool
	err error
}

func MakeBufferedReader[T any](r Reader[T]) BufferedReader[T] {
	return BufferedReader[T] {
		r: r,
		buf: make([]T, 0, bufsiz),
		idx: 0,
		closed: false,
	}
}

func NewBufferedReader[T any](r Reader[T]) *BufferedReader[T] {
	b := MakeBufferedReader(r)
	return &b
}

func (b *BufferedReader[T]) grow() (n int, err error) {
	b.idx = 0
	b.buf = b.buf[:cap(b.buf)]
	n, err = b.r.Read(b.buf)

	if err == io.EOF {
		b.closed = true
	}

	b.buf = b.buf[:n]
	return n, err
}

func (b *BufferedReader[T]) growIfNeeded() (n int, err error) {
	if b.idx >= len(b.buf) && !b.closed {
		n, err = b.grow()
	}
	return n, err
}

func (b *BufferedReader[T]) Read(p []T) (n int, err error) {
	if b.closed && b.idx >= len(b.buf) {
		return 0, io.EOF
	}

	_, b.err = b.growIfNeeded()

	n = copy(p, b.buf[b.idx:])
	b.idx += n

	if b.err != nil && b.err != io.EOF {
		return n, err
	}

	if b.closed && b.idx >= len(b.buf) {
		return n, io.EOF
	}

	return n, nil
}

func (b *BufferedReader[T]) Next() (T, bool) {
	if b.closed && b.idx >= len(b.buf) {
		var t T
		return t, false
	}

	for b.idx >= len(b.buf) && !b.closed {
		_, b.err = b.grow()
		if b.err != nil {
			b.closed = true
		}
	}

	if b.closed && b.idx >= len(b.buf) {
		var t T
		return t, false
	}

	out := b.buf[b.idx]
	b.idx++
	return out, true
}
