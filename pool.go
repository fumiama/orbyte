// Package orbyte is a lightweight & safe (buffer-writer | general object) pool.
package orbyte

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// Pool lightweight general pool.
type Pool[T any] struct {
	pooler   Pooler[T]
	pool     sync.Pool
	countin  int32
	countout int32
	noputbak bool
	manudstr bool
}

// NewPool make a new pool from custom pooler.
func NewPool[T any](pooler Pooler[T]) *Pool[T] {
	p := new(Pool[T])
	p.pooler = pooler
	p.pool.New = func() any {
		return &Item[T]{pool: p}
	}
	return p
}

// SetNoPutBack make it panic on every use-after-destroy.
//
// Enable this to detect coding errors.
func (pool *Pool[T]) SetNoPutBack(on bool) {
	pool.noputbak = on
}

// SetManualDestroy mark that user must manually
// run Item.Destroy().
func (pool *Pool[T]) SetManualDestroy(on bool) {
	pool.manudstr = on
}

func (pool *Pool[T]) incin() {
	atomic.AddInt32(&pool.countin, 1)
}

func (pool *Pool[T]) decin() {
	atomic.AddInt32(&pool.countin, -1)
}

func (pool *Pool[T]) incout() {
	atomic.AddInt32(&pool.countout, 1)
}

func (pool *Pool[T]) decout() {
	atomic.AddInt32(&pool.countout, -1)
}

func (pool *Pool[T]) newempty() *Item[T] {
	item := pool.pool.Get().(*Item[T])
	if item.stat.hasdestroyed() { // is recycled
		pool.decin()
	}
	item.stat = status(0)
	pool.incout()
	if !pool.manudstr {
		item.setautodestroy()
	}
	return item
}

func (pool *Pool[T]) put(item *Item[T]) {
	if !pool.manudstr {
		runtime.SetFinalizer(item, nil)
	}

	item.stat.setdestroyed(true)
	item.cfg = nil

	if pool.noputbak {
		return
	}
	pool.pool.Put(item)

	pool.decout()
	pool.incin()
}

// New call this to generate an item.
func (pool *Pool[T]) New(config any) *Item[T] {
	item := pool.newempty()
	item.cfg = config
	item.stat.setbuffered(true)
	item.val = pool.pooler.New(config, item.val)
	return item
}

// InvolveItem[T any] involve external object into pool.
//
// After that, you must only use the object through Item.
func (pool *Pool[T]) Involve(config, obj any) *Item[T] {
	item := pool.newempty()
	item.cfg = config
	item.stat.setbuffered(true)
	item.val = pool.pooler.Parse(obj, item.val)
	return item
}

// ParseItem[T any] safely convert obj into pool item without copy.
//
// You can still use the original object elsewhere.
func (pool *Pool[T]) Parse(config, obj any) *Item[T] {
	item := pool.newempty()
	item.cfg = config
	item.val = pool.pooler.Parse(obj, item.val)
	return item
}

// CountItems returns total item count outside and inside.
func (pool *Pool[T]) CountItems() (outside, inside int32) {
	return atomic.LoadInt32(&pool.countout), atomic.LoadInt32(&pool.countin)
}
