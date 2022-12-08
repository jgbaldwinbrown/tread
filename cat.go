package tread

import (
	"io"
)

type multi[T any] struct {
	readers []Reader[T]
}

func multiReaderInternal[T any](rs ...Reader[T]) *multi[T] {
	m := new(multi[T])
	m.readers = rs
	return m
}

func MultiReader[T any](rs ...Reader[T]) Reader[T] {
	return multiReaderInternal(rs...)
}

func (m *multi[T]) Read(p []T) (n int, err error) {
	for len(m.readers) > 0 && n < len(p) {
		onen, oneerr := m.readers[0].Read(p[n:])
		n += onen;
		if (oneerr == io.EOF) {
			m.readers = m.readers[1:]
			continue
		}
		if (oneerr != nil) {
			return n, oneerr
		}
	}
	if len(m.readers) <= 0 {
		return n, io.EOF
	}
	return n, nil
}
