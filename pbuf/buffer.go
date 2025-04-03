package pbuf

import (
	"bytes"

	"github.com/fumiama/orbyte"
)

// NewBuffer wraps bytes.NewBuffer into Item.
func (bufferPool BufferPool[USRDAT]) NewBuffer(
	buf []byte,
) *orbyte.Item[UserBuffer[USRDAT]] {
	return bufferPool.p.New(buf)
}

// InvolveBuffer involve external *bytes.Buffer into Item.
func (bufferPool BufferPool[USRDAT]) InvolveBuffer(
	buf *bytes.Buffer,
) *orbyte.Item[UserBuffer[USRDAT]] {
	return bufferPool.p.Involve(buf.Len(), buf)
}

// ParseBuffer convert external *bytes.Buffer into Item.
func (bufferPool BufferPool[USRDAT]) ParseBuffer(
	buf *bytes.Buffer,
) *orbyte.Item[UserBuffer[USRDAT]] {
	return bufferPool.p.Parse(buf.Len(), buf)
}

func (bufferPool BufferPool[USRDAT]) CountItems() (outside int32, inside int32) {
	return bufferPool.p.CountItems()
}
