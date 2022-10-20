package tread

import (
	"io"
)

const bufsize int = 4096

type Reader[T any] interface {
	Read(p []T) (n int, err error)
}

type Writer[T any] interface {
	Write(p []T) (n int, err error)
}

type Closer[T any] interface {
	Close() error
}

type ReadCloser[T any] interface {
	Reader[T]
	Closer[T]
}

type WriteCloser[T any] interface {
	Writer[T]
	Closer[T]
}

type ReadWriteCloser[T any] interface {
	Reader[T]
	Writer[T]
	Closer[T]
}

func Copy[T any](dest Writer[T], src Reader[T]) (n int, err error) {
	if wt, ok := src.(WriterTo[T]); ok {
		return wt.WriteTo(dest)
	}
	if rt, ok := dest.(ReaderFrom[T]); ok {
		return rt.ReadFrom(src)
	}
	buf := make([]T, bufsize)
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

func Next[T any](r Reader[T]) (T, bool) {
	buf := [1]T{}
	n, _ := r.Read(buf[:])
	return buf[0], n == 1
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

func (s *SliceBuffer[T]) ReadFrom(r Reader[T]) (n int, err error) {
	for {
		buf := s.slice[len(s.slice):cap(s.slice)]
		nread := 0
		nread, err = r.Read(buf)
		s.slice = s.slice[:len(s.slice)+nread]
		n += nread
		if err != nil {
			if err != io.EOF {
				return n, err
			}
			return n, nil
		}

		if len(s.slice) == cap(s.slice) {
			newslice := make([]T, len(s.slice), cap(s.slice) * 2 + 1)
			copy(newslice, s.slice)
			s.slice = newslice
		}
	}
	return n, nil
}

func (s *SliceBuffer[T]) WriteTo(w Writer[T]) (n int, err error) {
	nwritten := 0
	for s.idx < len(s.slice) {
		nwritten, err = w.Write(s.slice[s.idx:])
		s.idx += nwritten
		n += nwritten
		if err != nil {
			if err != io.EOF {
				return n, err
			}
			return n, nil
		}
	}
	return n, nil
}
