package pbuf

import (
	"bytes"
	"runtime"

	"github.com/fumiama/orbyte"
)

// Bytes wrap pooled buffer into []byte
// while sharing the same pool.
type Bytes struct {
	buf *orbyte.Item[bytes.Buffer]
	dat []byte
}

// BufferItemToBytes convert between *orbyte.Item[bytes.Buffer]
// and Bytes.
//
// Please notice that Bytes cannnot convert back to
// *orbyte.Item[bytes.Buffer] again.
func BufferItemToBytes(buf *orbyte.Item[bytes.Buffer]) Bytes {
	return Bytes{buf: buf, dat: buf.Pointer().Bytes()}
}

// NewBytes alloc sz bytes.
func (bufferPool BufferPool) NewBytes(sz int) Bytes {
	buf := bufferPool.p.New(sz)
	return Bytes{buf: buf, dat: buf.Pointer().Bytes()[:sz]}
}

// InvolveBytes involve outside buf into pool.
func (bufferPool BufferPool) InvolveBytes(b ...byte) Bytes {
	buf := bufferPool.p.Involve(len(b), bytes.NewBuffer(b))
	return Bytes{buf: buf, dat: buf.Pointer().Bytes()[:len(b)]}
}

// ParseBytes convert outside bytes to Bytes safely
// without adding it into pool.
func (bufferPool BufferPool) ParseBytes(b ...byte) Bytes {
	buf := bufferPool.p.Parse(len(b), bytes.NewBuffer(b))
	return Bytes{buf: buf, dat: buf.Pointer().Bytes()}
}

// HasInit whether this Bytes is made by pool or
// just declared.
func (b Bytes) HasInit() bool {
	return b.buf != nil
}

// Trans please refer to Item.Trans().
func (b Bytes) Trans() (tb Bytes) {
	tb.buf = b.buf.Trans()
	tb.dat = b.dat
	return
}

// Len of slice.
func (b Bytes) Len() int {
	return len(b.dat)
}

// Cap of slice.
func (b Bytes) Cap() int {
	return cap(b.dat)
}

// Bytes is the inner value.
func (b Bytes) Bytes() []byte {
	return b.dat
}

// Ref please refer to Item.Ref().
func (b Bytes) Ref() (rb Bytes) {
	rb.buf = b.buf.Ref()
	rb.dat = b.dat
	return
}

// Copy please refer to Item.Copy().
func (b Bytes) Copy() (cb Bytes) {
	cb.buf = b.buf.Copy()
	cb.dat = cb.buf.Pointer().Bytes()
	return
}

// SliceFrom dat[from:] with Ref.
func (b Bytes) SliceFrom(from int) Bytes {
	if b.buf.IsTrans() {
		return InvolveBytes(b.dat[from:]...)
	}
	nb := b.Ref()
	skip(nb.buf.Pointer(), from)
	nb.dat = b.dat[from:]
	return nb
}

// SliceTo dat[:to] with Ref.
func (b Bytes) SliceTo(to int) Bytes {
	if b.buf.IsTrans() {
		return InvolveBytes(b.dat[:to]...)
	}
	nb := b.Ref()
	nb.buf.Pointer().Truncate(to)
	nb.dat = b.dat[:to]
	return nb
}

// Slice dat[from:to] with Ref.
func (b Bytes) Slice(from, to int) Bytes {
	if b.buf.IsTrans() {
		return InvolveBytes(b.dat[from:to]...)
	}
	nb := b.Ref()
	buf := nb.buf.Pointer()
	skip(buf, from)
	buf.Truncate(to - from)
	nb.dat = b.dat[from:to]
	return nb
}

// KeepAlive marks Bytes value as reachable.
func (b Bytes) KeepAlive() {
	_ = b.buf
	_ = b.dat
	runtime.KeepAlive(b.buf)
	runtime.KeepAlive(b.dat)
}
