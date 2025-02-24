package orbyte

import (
	"runtime"
	"sync/atomic"
)

// Item represents a thread-safe user-defined value.
//
// Only the item that has ownership can be a pointer.
// Do not copy neither Item nor *Item by yourself.
// You must always use the given methods.
type Item[T any] struct {
	pool *Pool[T]
	cfg  any

	stat status

	val T
}

// Trans ownership to a new item and
// destroy original item immediately.
//
// Call this function to drop your ownership
// before passing it to another function
// that do not return to you.
func (b *Item[T]) Trans() (tb *Item[T]) {
	if b.stat.hasdestroyed() {
		panic("use after destroy")
	}
	tb = b.pool.newempty()
	*tb = *b
	tb.stat = status(atomic.SwapUintptr(
		(*uintptr)(&b.stat), uintptr(destroyedstatus),
	))
	b.pool.put(b)
	return tb
}

// Unwrap use value of the item
func (b *Item[T]) Unwrap() T {
	if b.stat.hasdestroyed() {
		panic("use after destroy")
	}
	return b.val
}

// Ref gens new item without ownership.
//
// It's a safe reference, thus calling this
// will not destroy the original item
// comparing with Trans().
func (b *Item[T]) Ref() (rb *Item[T]) {
	if b.stat.hasdestroyed() {
		panic("use after destroy")
	}
	rb = b.pool.newempty()
	*rb = *b
	rb.stat.setbuffered(false)
	return
}

// Copy data completely with separated ownership.
func (b *Item[T]) Copy() (cb *Item[T]) {
	if b.stat.hasdestroyed() {
		panic("use after destroy")
	}
	cb = b.pool.New(b.cfg)
	b.pool.pooler.Copy(&cb.val, &b.val)
	return
}

// Destroy item and put it back to pool.
func (b *Item[T]) Destroy() {
	stat := status(atomic.SwapUintptr(
		(*uintptr)(&b.stat), uintptr(destroyedstatus),
	))
	if stat.hasdestroyed() {
		panic("use after destroy")
	}
	if b.stat.isbuffered() {
		b.pool.pooler.Reset(&b.val)
	}
	b.pool.put(b)
}

// setautodestroy item on GC.
//
// Only can call once.
func (b *Item[T]) setautodestroy() *Item[T] {
	runtime.SetFinalizer(b, func(item *Item[T]) {
		// no one is using, no concurrency issue.
		if item.stat.hasdestroyed() {
			panic("unexpected hasdestroyed")
		}
		if item.stat.isbuffered() {
			item.pool.pooler.Reset(&item.val)
		}
		item.stat.setdestroyed(true)
		item.pool.put(item)
	})
	return b
}
