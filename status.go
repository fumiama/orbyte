package orbyte

import "sync/atomic"

const (
	statusisbuffered = 1 << iota
	statusdestroyed
	statusisintrans
)

type status uintptr

var destroyedstatus status

func init() {
	destroyedstatus.setdestroyed(true)
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

func (c *status) isintrans() bool {
	return c.loadbool(statusisintrans)
}

func (c *status) setintrans(v bool) {
	c.setbool(v, statusisintrans)
}
