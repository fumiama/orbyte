package pbuf

import (
	"bytes"
	"io"
	"unsafe"

	"github.com/fumiama/orbyte"
)

// NewBuffer wraps bytes.NewBuffer into Item.
func (bufferPool BufferPool) NewBuffer(buf []byte) *orbyte.Item[bytes.Buffer] {
	return bufferPool.p.New(buf)
}

// InvolveBuffer involve external *bytes.Buffer into Item.
func (bufferPool BufferPool) InvolveBuffer(buf *bytes.Buffer) *orbyte.Item[bytes.Buffer] {
	return bufferPool.p.Involve(buf.Len(), buf)
}

// ParseBuffer convert external *bytes.Buffer into Item.
func (bufferPool BufferPool) ParseBuffer(buf *bytes.Buffer) *orbyte.Item[bytes.Buffer] {
	return bufferPool.p.Parse(buf.Len(), buf)
}

// A Buffer is a variable-sized buffer of bytes with Read and Write methods.
// The zero value for Buffer is an empty buffer ready to use.
type buffer struct {
	buf      []byte // contents are the bytes buf[off : len(buf)]
	off      int    // read at &buf[off], write at &buf[len(buf)]
	lastRead readOp // last read operation, so that Unread* can work correctly.
}

func skip(w *bytes.Buffer, n int) (int, error) {
	if n == 0 {
		return 0, nil
	}
	b := (*buffer)(unsafe.Pointer(w))
	b.lastRead = opInvalid
	if len(b.buf) <= b.off {
		// Buffer is empty, reset to recover space.
		w.Reset()
		if n == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
	n = minnum(n, len(b.buf[b.off:]))
	b.off += n
	if n > 0 {
		b.lastRead = opRead
	}
	return n, nil
}

// The readOp constants describe the last action performed on
// the buffer, so that UnreadRune and UnreadByte can check for
// invalid usage. opReadRuneX constants are chosen such that
// converted to int they correspond to the rune size that was read.
type readOp int8

// Don't use iota for these, as the values need to correspond with the
// names and comments, which is easier to see when being explicit.
const (
	opRead      readOp = -1 // Any other read operation.
	opInvalid   readOp = 0  // Non-read operation.
	opReadRune1 readOp = 1  // Read rune of size 1.
	opReadRune2 readOp = 2  // Read rune of size 2.
	opReadRune3 readOp = 3  // Read rune of size 3.
	opReadRune4 readOp = 4  // Read rune of size 4.
)

// minnum 返回两数最小值，该函数将被内联
func minnum[T int | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64](a, b T) T {
	if a > b {
		return b
	}
	return a
}
