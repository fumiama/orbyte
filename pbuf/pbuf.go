// Package pbuf is a lightweight pooled buffer.
package pbuf

import (
	"bytes"

	"github.com/fumiama/orbyte"
)

var bufferPool = NewBufferPool[struct{}]()

type (
	Pool    = BufferPool[struct{}]
	OBuffer = orbyte.Item[Buffer]
	Buffer  = UserBuffer[struct{}]
	Bytes   = UserBytes[struct{}]
)

type BufferPool[USRDAT any] struct {
	*orbyte.Pool[UserBuffer[USRDAT]]
}

// GetDefaultBufferPool exposes inner bufferPool value.
func GetDefaultBufferPool() Pool {
	return bufferPool
}

func NewBufferPool[USRDAT any]() BufferPool[USRDAT] {
	return BufferPool[USRDAT]{
		orbyte.NewPool[UserBuffer[USRDAT]](bufpooler[USRDAT]{}),
	}
}

// NewBuffer wraps bytes.NewBuffer into Item.
func NewBuffer(buf []byte) *OBuffer {
	return bufferPool.NewBuffer(buf)
}

// InvolveBuffer involve external *bytes.Buffer into Item.
func InvolveBuffer(buf *bytes.Buffer) *OBuffer {
	return bufferPool.InvolveBuffer(buf)
}

// ParseBuffer convert external *bytes.Buffer into Item.
func ParseBuffer(buf *bytes.Buffer) *OBuffer {
	return bufferPool.ParseBuffer(buf)
}

// NewBytes alloc sz bytes.
func NewBytes(sz int) Bytes {
	return bufferPool.NewBytes(sz)
}

// NewBytes alloc sz bytes without involving.
func NewLargeBytes(sz int) Bytes {
	return bufferPool.NewLargeBytes(sz)
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

// CountItems see Pool.CountItems
func CountItems() (outside int32, inside int32) {
	return bufferPool.CountItems()
}

// SetNoPutBack see Pool.SetNoPutBack
func SetNoPutBack(on bool) {
	bufferPool.SetNoPutBack(on)
}

// SetSyncItem see Pool.SetSyncItem
func SetSyncItem(on bool) {
	bufferPool.SetSyncItem(on)
}

// LimitInput see Pool.LimitInput
func LimitInput(n int32) {
	bufferPool.LimitInput(n)
}

// LimitInput see Pool.LimitOutput
func LimitOutput(n int32) {
	bufferPool.LimitOutput(n)
}
