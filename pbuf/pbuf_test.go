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
		b.Trans().SliceFrom(0).SliceTo(i).Destroy()
	}
	runtime.GC()
	out, in := bufferPool.p.CountItems()
	t.Log(out, in)
	if out != 0 {
		t.Fail()
	}
}
