package orbyte

import (
	"crypto/rand"
	"runtime"
	"sync"
	"testing"
)

func TestPool(t *testing.T) {
	p := NewPool[[]byte](simplepooler{})
	x := p.New(200)
	x.ManualDestroy()
	out, in := p.CountItems()
	t.Log("out", out, "in", in)
	if out != 0 || in != 1 {
		t.Fatal("unexpected behavior")
	}
	for i := 0; i < 2000; i++ {
		item := p.New(i)
		out, in = p.CountItems()
		if out != 1 || in != 0 {
			t.Fatal("unexpected behavior")
		}
		item.ManualDestroy()
	}
	out, in = p.CountItems()
	t.Log("out", out, "in", in)
	if out != 0 {
		t.Fatal("unexpected behavior")
	}
	wg := sync.WaitGroup{}
	for i := 0; i < 4096; i++ {
		item := p.New(i)
		wg.Add(1)
		go useranddes(item.Copy(), &wg)
		wg.Add(1)
		go userv(item.Trans(), &wg)
	}
	wg.Wait()
	runtime.GC()
	out, in = p.CountItems()
	t.Log("out", out, "in", in)
	if out != 0 {
		t.Fatal("unexpected behavior")
	}
}

func useranddes(item *Item[[]byte], wg *sync.WaitGroup) {
	defer wg.Done()
	item.V(func(b []byte) {
		rand.Read(b)
	})
	item.ManualDestroy()
}

func userv(b []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	rand.Read(b)
}

type simplepooler struct{}

func (simplepooler) New(config any, pooled []byte) []byte {
	if cap(pooled) >= config.(int) {
		return pooled[:config.(int)]
	}
	return make([]byte, config.(int))
}

func (simplepooler) Parse(obj any, pooled []byte) []byte {
	src := obj.([]byte)
	if cap(pooled) >= len(src) {
		copy(pooled[:len(src)], src)
		return pooled[:len(src)]
	}
	return obj.([]byte)
}

func (simplepooler) Reset(item *[]byte) {
	*item = (*item)[:0]
}

func (simplepooler) Copy(dst, src *[]byte) {
	copy(*dst, *src)
}
