package main

type ReaderFrom[T any] interface {
	ReadFrom(Reader[T]) (int, error)
}

type WriterTo[T any] interface {
	WriteTo(Writer[T]) (int, error)
}
