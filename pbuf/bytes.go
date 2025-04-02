package pbuf

import (
	"bytes"
	"runtime"

	"github.com/fumiama/orbyte"
)

// UserBytes wrap pooled buffer into []byte
// while sharing the same pool.
type UserBytes[USRDAT any] struct {
	buf  *orbyte.Item[UserBuffer[USRDAT]]
	a, b int
}

// BufferItemToBytes convert between *Buffer
// and Bytes.
func BufferItemToBytes[USRDAT any](
	buf *orbyte.Item[UserBuffer[USRDAT]],
) (b UserBytes[USRDAT]) {
	b.buf = buf
	buf.P(func(buf *UserBuffer[USRDAT]) {
		b.b = buf.Len()
	})
	return
}

// B directly use inner buf data and USRDAT safely.
func (b UserBytes[USRDAT]) B(f func([]byte, *USRDAT)) {
	b.buf.P(func(ub *UserBuffer[USRDAT]) {
		f(ub.Buffer.Bytes(), &ub.DAT)
		runtime.KeepAlive(b.buf)
	})
}

// NewBytes alloc sz bytes.
func (bufferPool BufferPool[USRDAT]) NewBytes(sz int) (b UserBytes[USRDAT]) {
	buf := bufferPool.p.New(sz)
	b.buf = buf
	buf.P(func(buf *UserBuffer[USRDAT]) {
		b.b = buf.Len()
	})
	return
}

// InvolveBytes involve outside buf into pool.
func (bufferPool BufferPool[USRDAT]) InvolveBytes(p ...byte) (b UserBytes[USRDAT]) {
	buf := bufferPool.p.Involve(len(p), bytes.NewBuffer(p))
	b.buf = buf
	buf.P(func(buf *UserBuffer[USRDAT]) {
		b.b = buf.Len()
	})
	return
}

// ParseBytes convert outside bytes to Bytes safely
// without adding it into pool.
func (bufferPool BufferPool[USRDAT]) ParseBytes(p ...byte) (b UserBytes[USRDAT]) {
	buf := bufferPool.p.Parse(len(p), bytes.NewBuffer(p))
	b.buf = buf
	buf.P(func(buf *UserBuffer[USRDAT]) {
		b.b = buf.Len()
	})
	return
}

// HasInit whether this Bytes is made by pool or
// just declared.
func (b UserBytes[USRDAT]) HasInit() bool {
	return b.buf != nil
}

// Trans please refer to Item.Trans().
func (b UserBytes[USRDAT]) Trans() []byte {
	buf := b.buf.Trans()
	return buf.Bytes()[b.a:b.b]
}

// Len of slice.
func (b UserBytes[USRDAT]) Len() int {
	return b.b - b.a
}

// Cap of slice.
func (b UserBytes[USRDAT]) Cap() (c int) {
	b.buf.P(func(b *UserBuffer[USRDAT]) {
		c = b.Cap()
	})
	return c
}

// V use the inner value safely
func (b UserBytes[USRDAT]) V(f func([]byte)) {
	b.buf.P(func(buf *UserBuffer[USRDAT]) {
		f(buf.Bytes()[b.a:b.b])
		runtime.KeepAlive(b.buf)
	})
}

// Copy please refer to Item.Copy().
func (b UserBytes[USRDAT]) Copy() (cb UserBytes[USRDAT]) {
	cb.buf = b.buf.Copy()
	cb.a, cb.b = b.a, b.b
	return
}

// SliceFrom dat[from:] with Ref.
func (b UserBytes[USRDAT]) SliceFrom(from int) UserBytes[USRDAT] {
	return UserBytes[USRDAT]{buf: b.buf, a: b.a + from, b: b.b}
}

// SliceTo dat[:to] with Ref.
func (b UserBytes[USRDAT]) SliceTo(to int) UserBytes[USRDAT] {
	return UserBytes[USRDAT]{buf: b.buf, a: b.a, b: b.a + to}
}

// Slice dat[from:to] with Ref.
func (b UserBytes[USRDAT]) Slice(from, to int) UserBytes[USRDAT] {
	return UserBytes[USRDAT]{buf: b.buf, a: b.a + from, b: b.a + to}
}
