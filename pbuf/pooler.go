package pbuf

import (
	"bytes"
	"io"
	"reflect"
)

type bufpooler struct{}

func (bufpooler) New(config any, pooled bytes.Buffer) bytes.Buffer {
	switch c := config.(type) {
	case int:
		pooled.Grow(c)
		return pooled
	case []byte:
		if len(c) > 0 || pooled.Cap() < cap(c) {
			return *bytes.NewBuffer(c)
		}
		return pooled
	case string:
		pooled.Grow(len(c))
		pooled.WriteString(c)
		return pooled
	default:
		panic("config type " + reflect.ValueOf(config).Type().String() + " isn't supported")
	}
}

func (bufpooler) Parse(obj any, pooled bytes.Buffer) bytes.Buffer {
	switch o := obj.(type) {
	case *bytes.Buffer:
		return *o
	case bytes.Buffer:
		return o
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

func (bufpooler) Reset(item *bytes.Buffer) {
	// See https://golang.org/issue/23199
	const maxSize = 1 << 16
	if item.Cap() > maxSize { // drop large buffer
		*item = bytes.Buffer{}
		return
	}
	item.Reset()
}

func (bufpooler) Copy(dst, src *bytes.Buffer) {
	_, err := io.Copy(dst, src)
	if err != nil {
		panic(err)
	}
}
