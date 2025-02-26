package pbuf

import (
	"bytes"

	"github.com/fumiama/orbyte"
)

// Bytes wrap pooled buffer into []byte
// while sharing the same pool.
type Bytes struct {
	buf  *orbyte.Item[bytes.Buffer]
	a, b int
}

// BufferItemToBytes convert between *orbyte.Item[bytes.Buffer]
// and Bytes.
//
// Please notice that Bytes cannnot convert back to
// *orbyte.Item[bytes.Buffer] again.
func BufferItemToBytes(buf *orbyte.Item[bytes.Buffer]) (b Bytes) {
	b.buf = buf
	buf.P(func(buf *bytes.Buffer) {
		b.b = buf.Len()
	})
	return
}

// NewBytes alloc sz bytes.
func (bufferPool BufferPool) NewBytes(sz int) (b Bytes) {
	buf := bufferPool.p.New(sz)
	b.buf = buf
	buf.P(func(buf *bytes.Buffer) {
		b.b = buf.Len()
	})
	return
}

// InvolveBytes involve outside buf into pool.
func (bufferPool BufferPool) InvolveBytes(p ...byte) (b Bytes) {
	buf := bufferPool.p.Involve(len(p), bytes.NewBuffer(p))
	b.buf = buf
	buf.P(func(buf *bytes.Buffer) {
		b.b = buf.Len()
	})
	return
}

// ParseBytes convert outside bytes to Bytes safely
// without adding it into pool.
func (bufferPool BufferPool) ParseBytes(p ...byte) (b Bytes) {
	buf := bufferPool.p.Parse(len(p), bytes.NewBuffer(p))
	b.buf = buf
	buf.P(func(buf *bytes.Buffer) {
		b.b = buf.Len()
	})
	return
}

// HasInit whether this Bytes is made by pool or
// just declared.
func (b Bytes) HasInit() bool {
	return b.buf != nil
}

// Trans please refer to Item.Trans().
func (b Bytes) Trans() []byte {
	buf := b.buf.Trans()
	return buf.Bytes()[b.a:b.b]
}

// Len of slice.
func (b Bytes) Len() int {
	return b.b - b.a
}

// Cap of slice.
func (b Bytes) Cap() (c int) {
	b.buf.P(func(b *bytes.Buffer) {
		c = b.Cap()
	})
	return c
}

// V use the inner value safely
func (b Bytes) V(f func([]byte)) {
	b.buf.P(func(buf *bytes.Buffer) {
		f(buf.Bytes()[b.a:b.b])
	})
}

// Copy please refer to Item.Copy().
func (b Bytes) Copy() (cb Bytes) {
	cb.buf = b.buf.Copy()
	cb.a, cb.b = b.a, b.b
	return
}

// SliceFrom dat[from:] with Ref.
func (b Bytes) SliceFrom(from int) Bytes {
	return Bytes{buf: b.buf, a: b.a + from, b: b.b}
}

// SliceTo dat[:to] with Ref.
func (b Bytes) SliceTo(to int) Bytes {
	return Bytes{buf: b.buf, a: b.a, b: b.a + to}
}

// Slice dat[from:to] with Ref.
func (b Bytes) Slice(from, to int) Bytes {
	return Bytes{buf: b.buf, a: b.a + from, b: b.a + to}
}
