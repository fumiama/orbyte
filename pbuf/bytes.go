package pbuf

import (
	"bytes"

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
	x := buf.Unwrap()
	return Bytes{buf: buf, dat: x.Bytes()}
}

// NewBytes alloc sz bytes.
func (bufferPool BufferPool) NewBytes(sz int) Bytes {
	buf := bufferPool.p.New(sz)
	x := buf.Unwrap()
	return Bytes{buf: buf, dat: x.Bytes()}
}

// InvolveBytes involve outside buf into pool.
func (bufferPool BufferPool) InvolveBytes(b ...byte) Bytes {
	buf := bufferPool.p.Involve(len(b), bytes.NewBuffer(b))
	x := buf.Unwrap()
	return Bytes{buf: buf, dat: x.Bytes()}
}

// ParseBytes convert outside bytes to Bytes safely
// without adding it into pool.
func (bufferPool BufferPool) ParseBytes(b ...byte) Bytes {
	buf := bufferPool.p.Parse(len(b), bytes.NewBuffer(b))
	x := buf.Unwrap()
	return Bytes{buf: buf, dat: x.Bytes()}
}

// Trans please refer to Item.Trans().
func (b Bytes) Trans() (tb Bytes) {
	tb.buf = b.buf.Trans()
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
	return
}

// Copy please refer to Item.Copy().
func (b Bytes) Copy() (cb Bytes) {
	cb.buf = b.buf.Copy()
	return
}

// SliceFrom dat[from:] with Ref.
func (b Bytes) SliceFrom(from int) Bytes {
	nb := b.Ref()
	nb.dat = b.dat[from:]
	return nb
}

// SliceTo dat[:to] with Ref.
func (b Bytes) SliceTo(to int) Bytes {
	nb := b.Ref()
	nb.dat = b.dat[:to]
	return nb
}

// Slice dat[from:to] with Ref.
func (b Bytes) Slice(from, to int) Bytes {
	nb := b.Ref()
	nb.dat = b.dat[from:to]
	return nb
}

// Destroy please refer to Item.Destroy().
func (b Bytes) Destroy() {
	b.buf.Destroy()
}
