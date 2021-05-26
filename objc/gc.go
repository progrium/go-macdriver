package objc

import "runtime"

type gcReleased struct{ Object }

// GCReleased returns a wrapped Object which will automatically call Release
// when it is freed by the Go garbage collector.
// Afterwards, the caller should treat the input as if it had called Release,
// and use the returned Object instead.
func GCReleased(obj Object) Object {
	r := &gcReleased{obj}
	runtime.SetFinalizer(r, func(r *gcReleased) {
		r.Release()
	})
	return r
}
