package pbuf

import (
	"bytes"
	"io"
	"reflect"
	"unsafe"
)

// UserBuffer with customizable user data structure inside.
type UserBuffer[USRDAT any] struct {
	bytes.Buffer
	DAT USRDAT
}

type bufpooler[USRDAT any] struct{}

func (bufpooler[USRDAT]) New(config any, pooled UserBuffer[USRDAT]) UserBuffer[USRDAT] {
	switch c := config.(type) {
	case int:
		pooled.Grow(c)
		*(*[]byte)(unsafe.Pointer(&pooled)) = pooled.Bytes()[:c]
		if c != pooled.Len() {
			panic("unexpected bad buffer Grow")
		}
		return pooled
	case []byte:
		if len(c) > 0 || pooled.Cap() < cap(c) {
			buf := bytes.NewBuffer(c)
			if len(c) != buf.Len() {
				panic("unexpected bad bytes.NewBuffer")
			}
			return UserBuffer[USRDAT]{Buffer: *buf}
		}
		return pooled
	case string:
		pooled.WriteString(c)
		return pooled
	default:
		panic("config type " + reflect.ValueOf(config).Type().String() + " isn't supported")
	}
}

func (bufpooler[USRDAT]) Parse(obj any, pooled UserBuffer[USRDAT]) UserBuffer[USRDAT] {
	switch o := obj.(type) {
	case *bytes.Buffer:
		return UserBuffer[USRDAT]{Buffer: *o}
	case bytes.Buffer:
		return UserBuffer[USRDAT]{Buffer: o}
	case []byte:
		pooled.Write(o)
		return pooled
	case string:
		pooled.WriteString(o)
		return pooled
	case io.Reader:
		_, err := io.Copy(&pooled, o)
		if err != nil {
			panic(err)
		}
		return pooled
	default:
		panic("object type " + reflect.ValueOf(obj).Type().String() + " isn't supported")
	}
}

func (bufpooler[USRDAT]) Reset(item *UserBuffer[USRDAT]) {
	// See https://golang.org/issue/23199
	const maxSize = 1 << 16
	if item.Cap() > maxSize { // drop large buffer
		*item = UserBuffer[USRDAT]{}
		return
	}
	item.Reset()
}

func (bufpooler[USRDAT]) Copy(dst, src *UserBuffer[USRDAT]) {
	dst.Reset()
	srccp := *src
	_, err := io.Copy(dst, &srccp)
	if err != nil {
		panic(err)
	}
}
