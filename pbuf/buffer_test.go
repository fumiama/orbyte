package pbuf

import (
	"bytes"
	"crypto/rand"
	"io"
	"runtime"
	"testing"
)

func TestBuffer(t *testing.T) {
	testBuffer(NewBuffer(nil), t)
	testBuffer(NewBuffer(make([]byte, 0, 8192)), t)
	testBuffer(ParseBuffer(bytes.NewBuffer(nil)), t)
	testBuffer(ParseBuffer(bytes.NewBuffer(make([]byte, 0, 8192))), t)
	testBuffer(InvolveBuffer(bytes.NewBuffer(nil)), t)
	testBuffer(InvolveBuffer(bytes.NewBuffer(make([]byte, 0, 8192))), t)
}

func testBuffer(buf *OBuffer, t *testing.T) {
	buf.P(func(buf *Buffer) {
		if buf.Len() != 4096 {
			io.CopyN(buf, rand.Reader, 4096)
			if buf.Len() != 4096 {
				t.Fatal("got", buf.Len())
			}
		}
	})
	bufcp := buf.Copy()
	dat := buf.Trans()
	bufcp.P(func(bufcp *Buffer) {
		if bufcp.Len() != 4096 {
			t.Fatal("got", bufcp.Len())
		}
		if !bytes.Equal(bufcp.Bytes(), dat.Bytes()) {
			t.Fatal("unexpected")
		}
	})

	runtime.GC()
	runtime.Gosched()
	runtime.GC()

	out, in := bufferPool.p.CountItems()
	t.Log(out, in)
	if out != 0 {
		t.Fail()
	}
}
