package pbuf

import (
	"crypto/rand"
	"runtime"
	"testing"
)

func TestBytes(t *testing.T) {
	for i := 0; i < 4096; i++ {
		b := NewBytes(i)
		rand.Read(b.Bytes())
		b.Destroy()
	}
	runtime.GC()
	out, in := bufferPool.p.CountItems()
	t.Log(out, in)
	if out != 0 || in != 1 {
		t.Fail()
	}
}
