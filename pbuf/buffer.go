package pbuf

import (
	"bytes"

	"github.com/fumiama/orbyte"
)

// NewBuffer wraps bytes.NewBuffer into Item.
func (bufferPool BufferPool[USRDAT]) NewBuffer(
	buf []byte,
) *orbyte.Item[UserBuffer[USRDAT]] {
	return bufferPool.New(buf)
}

// InvolveBuffer involve external *bytes.Buffer into Item.
func (bufferPool BufferPool[USRDAT]) InvolveBuffer(
	buf *bytes.Buffer,
) *orbyte.Item[UserBuffer[USRDAT]] {
	return bufferPool.Involve(buf.Len(), buf)
}

// ParseBuffer convert external *bytes.Buffer into Item.
func (bufferPool BufferPool[USRDAT]) ParseBuffer(
	buf *bytes.Buffer,
) *orbyte.Item[UserBuffer[USRDAT]] {
	return bufferPool.Parse(buf.Len(), buf)
}
