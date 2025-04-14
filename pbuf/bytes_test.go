package pbuf

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	mrand "math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

// TestBytesSlice sometimes fails at first run because
// GC not collecting all unused items.
func TestBytesSlice(t *testing.T) {
	for i := 10; i < 4096; i++ {
		b := NewBytes(i)
		if b.Len() != i {
			t.Fatal("index", i, "excpet len", i, "but got", b.Len())
		}
		b.V(func(b []byte) {
			rand.Read(b)
		})
		buf := make([]byte, b.Len())
		b.V(func(b []byte) {
			copy(buf, b)
		})
		x := b.SliceFrom(5).SliceTo(i - 5 - 5)
		dat := x.Trans()
		if !bytes.Equal(buf[5:i-5], dat) {
			t.Log("exp:", hex.EncodeToString(buf[5:i-5]))
			t.Log("got:", hex.EncodeToString(dat))
			t.Fatal("index", i, "unexpected")
		}
	}
	runtime.GC()
	runtime.Gosched()
	runtime.GC()
	out, in := bufferPool.CountItems()
	t.Log(out, in)
	if out != 0 {
		t.Fail()
	}
}

func TestBytesInvolve(t *testing.T) {
	buf := make([]byte, 4096)
	rand.Read(buf)
	for i := 0; i < 4096; i++ {
		b := InvolveBytes(buf[:i]...)
		if b.Len() != i {
			t.Fatal("index", i, "excpet len", i, "but got", b.Len())
		}
		b.V(func(b []byte) {
			rand.Read(b)
		})
		if !bytes.Equal(b.Trans(), buf[:i]) {
			t.Fatal("index", i, "unexpected")
		}
	}
	runtime.GC()
	out, in := bufferPool.CountItems()
	t.Log(out, in)
	if out != 0 {
		t.Fail()
	}
}

func TestBytesParse(t *testing.T) {
	buf := make([]byte, 4096)
	rand.Read(buf)
	for i := 0; i < 4096; i++ {
		b := ParseBytes(buf[:i]...)
		if b.Len() != i {
			t.Fatal("index", i, "excpet len", i, "but got", b.Len())
		}
		if !bytes.Equal(b.Trans(), buf[:i]) {
			t.Fatal("index", i, "unexpected")
		}
	}
	runtime.GC()
	out, in := bufferPool.CountItems()
	t.Log(out, in)
	if out != 0 {
		t.Fail()
	}
}

func TestBytesCopy(t *testing.T) {
	buf := make([]byte, 4096)
	rand.Read(buf)
	for i := 10; i < 4096; i++ {
		a := ParseBytes(buf...)
		x := a.Slice(5, i-5)
		b := x.Copy()
		if b.Len() != i-10 {
			t.Fatal("index", i, "excpet len", i, "but got", b.Len())
		}
		b.V(func(b []byte) {
			rand.Read(b)
		})
		if bytes.Equal(b.Trans(), buf[:i]) {
			t.Fatal("index", i, "unexpected")
		}
	}
	runtime.GC()
	runtime.Gosched()
	runtime.GC()
	out, in := bufferPool.CountItems()
	t.Log(out, in)
	if out != 0 {
		t.Fail()
	}
}

func TestBytesTransMultithread(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 2048; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(mrand.Intn(10)))
			buf := NewBytes(65536)
			refer := make([]byte, 65536)
			rand.Read(refer)
			buf.V(func(b []byte) {
				copy(b, refer)
			})
			wg.Add(1)
			go func(buf []byte) {
				defer wg.Done()
				time.Sleep(time.Millisecond * time.Duration(mrand.Intn(10)))
				if !bytes.Equal(refer, buf) {
					panic("unexpected")
				}
			}(buf.Trans())
		}()
	}
	wg.Wait()
}
