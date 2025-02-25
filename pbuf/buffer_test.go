package pbuf

import (
	"bytes"
	"crypto/rand"
	"io"
	"runtime"
	"testing"

	"github.com/fumiama/orbyte"
)

func TestBuffer(t *testing.T) {
	testBuffer(NewBuffer(nil), t)
	testBuffer(NewBuffer(make([]byte, 0, 8192)), t)
	testBuffer(ParseBuffer(bytes.NewBuffer(nil)), t)
	testBuffer(ParseBuffer(bytes.NewBuffer(make([]byte, 0, 8192))), t)
	testBuffer(InvolveBuffer(bytes.NewBuffer(nil)), t)
	testBuffer(InvolveBuffer(bytes.NewBuffer(make([]byte, 0, 8192))), t)
}

func testBuffer(buf *orbyte.Item[bytes.Buffer], t *testing.T) {
	if buf.Pointer().Len() != 4096 {
		io.CopyN(buf.Pointer(), rand.Reader, 4096)
		if buf.Pointer().Len() != 4096 {
			t.Fatal("got", buf.Pointer().Len())
		}
	}

	bufcp := buf.Copy()
	if bufcp.Pointer().Len() != 4096 {
		t.Fatal("got", bufcp.Pointer().Len())
	}
	if !bytes.Equal(bufcp.Pointer().Bytes(), buf.Pointer().Bytes()) {
		t.Fatal("unexpected")
	}

	bufr := buf.Ref()
	if bufr.Pointer().Len() != 4096 {
		t.Fatal("got", bufr.Pointer().Len())
	}
	if !bytes.Equal(bufr.Pointer().Bytes(), buf.Pointer().Bytes()) {
		t.Fatal("unexpected")
	}
	bufr.Destroy()

	bufcp = bufcp.Trans()
	if bufcp.Pointer().Len() != 4096 {
		t.Fatal("got", bufcp.Pointer().Len())
	}
	if !bytes.Equal(bufcp.Pointer().Bytes(), buf.Pointer().Bytes()) {
		t.Fatal("unexpected")
	}
	bufcp.Destroy()

	runtime.GC()
	runtime.Gosched()
	runtime.GC()

	out, in := bufferPool.p.CountItems()
	t.Log(out, in)
	if out != 0 {
		t.Fail()
	}
}
