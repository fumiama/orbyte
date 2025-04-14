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
	stat status
	// align 64

	cfg any
	// align 64

	val T
}

// Trans disable inner val being reset by
// destroy and return a safe copy of val.
//
// This method is not thread-safe.
// Only call once on one item.
// The item will be destroyed after calling Trans().
//
// Use it to drop your ownership
// before passing val (not its pointer)
// to another function that is not controlled by you.
func (b *Item[T]) Trans() T {
	if b.stat.hasdestroyed() {
		panic("use after destroy")
	}
	if b.pool.issync {
		b.stat.setinsyncop(true)
	}
	val := b.val
	atomic.StoreUintptr(
		(*uintptr)(&b.stat), uintptr(destroyedstatus),
	)
	runtime.KeepAlive(b)
	b.destroybystat(0)
	return val
}

// HasInvolved whether this item is buffered
// and will be Reset on putting back.
func (b *Item[T]) HasInvolved() bool {
	if b.pool.issync {
		b.stat.setinsyncop(true)
		defer b.stat.setinsyncop(false)
	}
	return b.stat.isbuffered()
}

// V use value of the item.
//
// This operation is safe in function f.
func (b *Item[T]) V(f func(T)) {
	if b.stat.hasdestroyed() {
		panic("use after destroy")
	}
	if b.pool.issync {
		b.stat.setinsyncop(true)
		defer b.stat.setinsyncop(false)
	}
	f(b.val)
	runtime.KeepAlive(b)
}

// P use pointer value of the item.
//
// This operation is safe in function f.
func (b *Item[T]) P(f func(*T)) {
	if b.stat.hasdestroyed() {
		panic("use after destroy")
	}
	if b.pool.issync {
		b.stat.setinsyncop(true)
		defer b.stat.setinsyncop(false)
	}
	f(&b.val)
	runtime.KeepAlive(b)
}

// Copy data completely with separated ownership.
func (b *Item[T]) Copy() (cb *Item[T]) {
	if b.stat.hasdestroyed() {
		panic("use after destroy")
	}
	if b.pool.issync {
		b.stat.setinsyncop(true)
		defer b.stat.setinsyncop(false)
	}
	cb = b.pool.New(b.cfg)
	b.pool.pooler.Copy(&cb.val, &b.val)
	return
}

func (b *Item[T]) destroybystat(stat status) {
	switch {
	case stat.hasdestroyed():
		panic("destroy after destroy")
	case stat.isbuffered():
		b.pool.pooler.Reset(&b.val)
	default:
		var v T
		b.val = v
	}
	b.pool.put(b)
}

// ManualDestroy item and put it back to pool.
//
// Calling this method only when you're sure that
// no one will use it, or it will cause a panic.
func (b *Item[T]) ManualDestroy() {
	if b.pool.issync {
		b.stat.setinsyncop(true)
	}
	runtime.SetFinalizer(b, nil)
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
		item.destroybystat(item.stat)
	})
	return b
}
