// Package pbuf is a lightweight pooled buffer.
package pbuf

import (
	"bytes"

	"github.com/fumiama/orbyte"
)

var bufferPool = NewBufferPool()

type BufferPool struct {
	p *orbyte.Pool[bytes.Buffer]
}

func NewBufferPool() BufferPool {
	return BufferPool{p: orbyte.NewPool[bytes.Buffer](bufpooler{})}
}

// NewBuffer wraps bytes.NewBuffer into Item.
func NewBuffer(buf []byte) *orbyte.Item[bytes.Buffer] {
	return bufferPool.NewBuffer(buf)
}

// InvolveBuffer involve external *bytes.Buffer into Item.
func InvolveBuffer(buf *bytes.Buffer) *orbyte.Item[bytes.Buffer] {
	return bufferPool.InvolveBuffer(buf)
}

// ParseBuffer convert external *bytes.Buffer into Item.
func ParseBuffer(buf *bytes.Buffer) *orbyte.Item[bytes.Buffer] {
	return bufferPool.ParseBuffer(buf)
}

// NewBytes alloc sz bytes.
func NewBytes(sz int) Bytes {
	return bufferPool.NewBytes(sz)
}

// InvolveBytes involve outside buf into pool.
func InvolveBytes(b ...byte) Bytes {
	return bufferPool.InvolveBytes(b...)
}

// ParseBytes convert outside bytes to Bytes safely
// without adding it into pool.
func ParseBytes(b ...byte) Bytes {
	return bufferPool.ParseBytes(b...)
}
