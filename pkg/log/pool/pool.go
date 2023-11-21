package pool

import "sync"

type Pool[T any] struct {
	pool sync.Pool
}

func NewPool[T any](fn func() T) *Pool[T] {
	pool := &Pool[T]{
		pool: sync.Pool{
			New: func() interface{} {
				return fn()
			},
		},
	}
	return pool
}

func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

func (pl *Pool[T]) Put(x T) {
	pl.pool.Put(x)
}
