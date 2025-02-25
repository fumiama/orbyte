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
// The value in new item will not be Reset().
//
// Call this function to drop your ownership
// before passing it to another function
// that is not controlled by you.
func (b *Item[T]) Trans() (tb *Item[T]) {
	if b.stat.hasdestroyed() {
		panic("use after destroy")
	}
	tb = b.pool.newempty()
	*tb = *b
	tb.stat = status(atomic.SwapUintptr(
		(*uintptr)(&b.stat), uintptr(destroyedstatus),
	))
	tb.stat.setintrans(true)
	b.destroybystat(status(0))
	return tb
}

// IsTrans whether this item has been marked as trans.
func (b *Item[T]) IsTrans() bool {
	return b.stat.isintrans()
}

// Unwrap use value of the item
func (b *Item[T]) Unwrap() T {
	if b.stat.hasdestroyed() {
		panic("use after destroy")
	}
	return b.val
}

// Pointer use pointer value of the item
func (b *Item[T]) Pointer() *T {
	if b.stat.hasdestroyed() {
		panic("use after destroy")
	}
	return &b.val
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
	rb.stat.setintrans(false)
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

func (b *Item[T]) destroybystat(stat status) {
	switch {
	case stat.hasdestroyed():
		panic("use after destroy")
	case stat.isintrans():
		var v T
		b.val = v
	case stat.isbuffered():
		b.pool.pooler.Reset(&b.val)
	default:
		var v T
		b.val = v
	}
	b.pool.put(b)
}

// Destroy item and put it back to pool.
func (b *Item[T]) Destroy() {
	b.destroybystat(status(atomic.SwapUintptr(
		(*uintptr)(&b.stat), uintptr(destroyedstatus),
	)))
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
		item.destroybystat(item.stat)
	})
	return b
}
