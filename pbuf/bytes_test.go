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
		rand.Read(b.Bytes())
		buf := make([]byte, b.Len())
		copy(buf, b.Bytes())
		// test normal slice
		x := b.SliceFrom(5).SliceTo(i - 5 - 5)
		if !bytes.Equal(buf[5:i-5], x.Bytes()) {
			t.Log("exp:", hex.EncodeToString(buf[5:i-5]))
			t.Log("got:", hex.EncodeToString(x.Bytes()))
			t.Fatal("index", i, "unexpected")
		}
		x.Destroy()
		// test trans slice
		b = b.Trans().SliceFrom(5).SliceTo(i - 5 - 5)
		if !bytes.Equal(buf[5:i-5], b.Bytes()) {
			t.Log("exp:", hex.EncodeToString(buf[5:i-5]))
			t.Log("got:", hex.EncodeToString(b.Bytes()))
			t.Fatal("index", i, "unexpected")
		}
		b.Destroy()
	}
	runtime.GC()
	runtime.Gosched()
	runtime.GC()
	out, in := bufferPool.p.CountItems()
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
		rand.Read(b.Bytes())
		if !bytes.Equal(b.Bytes(), buf[:i]) {
			t.Fatal("index", i, "unexpected")
		}
		b.Destroy()
	}
	runtime.GC()
	out, in := bufferPool.p.CountItems()
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
		if !bytes.Equal(b.Bytes(), buf[:i]) {
			t.Fatal("index", i, "unexpected")
		}
		b.Destroy()
	}
	runtime.GC()
	out, in := bufferPool.p.CountItems()
	t.Log(out, in)
	if out != 0 {
		t.Fail()
	}
}

func TestBytesCopy(t *testing.T) {
	buf := make([]byte, 4096)
	rand.Read(buf)
	for i := 10; i < 4096; i++ {
		b := ParseBytes(buf...).Slice(5, i-5).Copy()
		if b.Len() != i-10 {
			t.Fatal("index", i, "excpet len", i, "but got", b.Len())
		}
		rand.Read(b.Bytes())
		t.Log("org:", hex.EncodeToString(buf[:i]))
		t.Log("new:", hex.EncodeToString(b.Bytes()))
		if bytes.Equal(b.Bytes(), buf[:i]) {
			t.Fatal("index", i, "unexpected")
		}
		b.Destroy()
	}
	runtime.GC()
	out, in := bufferPool.p.CountItems()
	t.Log(out, in)
	if out != 0 {
		t.Fail()
	}
}

func TestBytesTransMultithread(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 4096; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(mrand.Intn(10)))
			buf := NewBytes(65536)
			refer := make([]byte, 65536)
			rand.Read(refer)
			copy(buf.Bytes(), refer)
			wg.Add(1)
			go func(buf Bytes) {
				defer wg.Done()
				time.Sleep(time.Millisecond * time.Duration(mrand.Intn(10)))
				if !bytes.Equal(refer, buf.Bytes()) {
					panic("unexpected")
				}
				buf.Destroy()
			}(buf.Trans())
		}()
	}
	wg.Wait()
}
