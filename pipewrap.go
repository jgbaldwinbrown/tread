package tread

func PipeWrite[T any](f func(Writer[T])) Reader[T] {
	r, w := Pipe[T]()
	go func() {
		f(w)
		w.Close()
	}()
	return r
}

func PipeRead[T any](f func(Reader[T])) Writer[T] {
	r, w := Pipe[T]()
	go func() {
		f(r)
		r.Close()
	}()
	return w
}

func PipeJoin[T any](r Reader[T], w Writer[T]) (done <-chan struct{}) {
	adone := make(chan struct{})
	go func() {
		Copy(w, r)
		adone <- struct{}{}
	}()
	return adone
}
