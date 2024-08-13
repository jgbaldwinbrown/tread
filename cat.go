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
	if len(m.readers) <= 0 {
		return 0, io.EOF
	}

	n, err = m.readers[0].Read(p)

	if (err == io.EOF) {
		m.readers = m.readers[1:]
		if len(m.readers) <= 0 {
			return n, io.EOF
		}
		return n, nil
	}
	if (err != nil) {
		return n, err
	}

	return n, nil
}
