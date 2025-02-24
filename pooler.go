package orbyte

// Pooler connects to a user-defined struct.
type Pooler[T any] interface {
	New(config any, pooled T) T
	Parse(obj any, pooled T) T
	Reset(item *T)
	Copy(dst, src *T)
}
