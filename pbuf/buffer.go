package pbuf

import (
	"bytes"

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
