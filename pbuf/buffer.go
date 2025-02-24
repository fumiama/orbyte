package pbuf

import (
	"bytes"

	"github.com/fumiama/orbyte"
)

// NewBuffer wraps bytes.NewBuffer
func (bufferPool BufferPool) NewBuffer(buf []byte) *orbyte.Item[bytes.Buffer] {
	return bufferPool.p.New(buf)
}
