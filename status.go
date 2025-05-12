package orbyte

import (
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
)

const (
	statusisbuffered = 1 << iota
	statusdestroyed
	statusinsyncop
	statushasignored
)

type status uintptr

var destroyedstatus status

func init() {
	destroyedstatus.setdestroyed(true)
}

func getGoroutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(string(buf[:n]))[1]
	id, err := strconv.ParseInt(idField, 10, 64)
	if err != nil {
		panic(err)
	}
	return id
}
func (c status) mask(v bool, typ uintptr) (news status) {
	news = c
	if v {
		news |= status(typ)
	} else {
		news &= ^status(typ)
	}
	return
}

func (c *status) setbool(v bool, typ uintptr) {
	olds := atomic.LoadUintptr((*uintptr)(c))
	oldv := olds&typ != 0
	if oldv == v {
		return
	}
	news := status(olds).mask(v, typ)
	for !atomic.CompareAndSwapUintptr((*uintptr)(c), olds, uintptr(news)) {
		olds = atomic.LoadUintptr((*uintptr)(c))
		news = status(olds).mask(v, typ)
	}
}

// setboolunique return false on non-unique set
func (c *status) setboolunique(v bool, typ uintptr) bool {
	olds := atomic.LoadUintptr((*uintptr)(c))
	oldv := olds&typ != 0
	if oldv == v {
		return false
	}
	news := status(olds).mask(v, typ)
	for !atomic.CompareAndSwapUintptr((*uintptr)(c), olds, uintptr(news)) {
		olds = atomic.LoadUintptr((*uintptr)(c))
		oldv = olds&typ != 0
		if oldv == v {
			return false
		}
		news = status(olds).mask(v, typ)
	}
	return true
}

func (c *status) loadbool(typ uintptr) bool {
	return atomic.LoadUintptr((*uintptr)(c))&typ != 0
}

func (c *status) isbuffered() bool {
	return c.loadbool(statusisbuffered)
}

func (c *status) setbuffered(v bool) {
	c.setbool(v, statusisbuffered)
}

func (c *status) hasdestroyed() bool {
	return c.loadbool(statusdestroyed)
}

func (c *status) setdestroyed(v bool) {
	c.setbool(v, statusdestroyed)
}

func (c *status) setinsyncop(v bool) bool {
	return c.setboolunique(v, statusinsyncop)
}

func (c *status) hasignored() bool {
	return c.loadbool(statushasignored)
}

func (c *status) setignored(v bool) {
	c.setbool(v, statushasignored)
}
