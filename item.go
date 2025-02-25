package orbyte

import (
	"runtime"
	"strconv"
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
	ref  *Item[T]
	refc int32 // refc -1 means transferred / destroyed

	val T
}

func (b *Item[T]) incref() {
	atomic.AddInt32(&b.refc, 1)
}

func (b *Item[T]) decref() {
	atomic.AddInt32(&b.refc, -1)
}

// Trans ownership to a new item and
// destroy original item immediately.
//
// The value in new item will not be Reset().
//
// Call this function to drop your ownership
// before passing it to another function
// that is not controlled by you.
//
// Avoid to call this function after calling Ref().
func (b *Item[T]) Trans() (tb *Item[T]) {
	if b.stat.hasdestroyed() {
		panic("use after destroy")
	}
	if b.ref != nil {
		panic("cannot trans ref")
	}
	tb = b.pool.newempty()
	*tb = *b
	tb.stat = status(atomic.SwapUintptr(
		(*uintptr)(&b.stat), uintptr(destroyedstatus),
	))
	tb.refc = 0
	tb.stat.setintrans(true)
	b.destroybystat(status(0))
	return tb
}

// IsTrans whether this item has been marked as trans.
func (b *Item[T]) IsTrans() bool {
	return b.stat.isintrans()
}

// IsRef whether this item is a reference.
func (b *Item[T]) IsRef() bool {
	return b.ref != nil
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
	rb.ref = b
	rb.refc = 0
	b.incref()
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
	if !atomic.CompareAndSwapInt32(&b.refc, 0, -1) {
		if b.refc < 0 {
			panic("use imm. after destroy")
		}
		panic("cannot destroy: " + strconv.Itoa(int(b.refc)) + " refs remained")
	}
	if b.ref != nil {
		defer b.ref.decref()
	}
	switch {
	case stat.hasdestroyed():
		panic("use after put back to pool")
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

// ManualDestroy item and put it back to pool.
//
// Calling this method without setting pool.SetManualDestroy(true)
// can probably cause panic.
func (b *Item[T]) ManualDestroy() {
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
