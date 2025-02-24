# orbyte
Lightweight &amp; Safe (buffer-writer | general object) pool.

## Quick Start
```go
import (
    "crypto/rand"

    "github.com/fumiama/orbyte/pbuf"
)

func main() {
    b := pbuf.NewBytes(1024) // Allocate Bytes from pool.
    rand.Read(b.Bytes())     // Do sth.
    b.Destroy() // Optional, can be auto-destroyed on GC.
}
```
