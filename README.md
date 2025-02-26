# orbyte
Lightweight &amp; Safe (buffer-writer | general object) pool.

## Quick Start
```go
package main

import (
	"crypto/rand"
	"io"

	"github.com/fumiama/orbyte/pbuf"
)

func main() {
	buf := pbuf.NewBuffer(nil) // Allocate Buffer from pool.
	buf.P(func(buf *pbuf.Buffer) {
		io.CopyN(buf, rand.Reader, 4096) // Do sth.
	})
	// After that, buf will be auto-reused on GC.

	b := pbuf.NewBytes(1024) // Allocate Bytes from pool.
	b.V(func(b []byte) {
		rand.Read(b) // Do sth.
	})
	// After that, b will be auto-reused on GC.

}
```
