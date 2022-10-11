package tread

import (
	"io"
)

type Chan[T any] chan T

func (c Chan[T]) Read(p []T) (n int, err error) {
	for t := range c {
		p[n] = t
		n++
	}
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}

func (c Chan[T]) Write(p []T) (n int, err error) {
	for _, t := range p {
		c <- t
		n++
	}
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}

func ReadChan[T any]( c chan T) Reader[T] {
	return Chan[T](c)
}

func WriteChan[T any]( c chan T) Writer[T] {
	return Chan[T](c)
}
