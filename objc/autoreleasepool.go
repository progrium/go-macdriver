package objc

/*
#cgo LDFLAGS: -lobjc -framework Foundation
extern void objc_autoreleasePoolPop(void* pool);
extern void* objc_autoreleasePoolPush(void);
*/
import "C"
import "runtime"

func Autorelease(body func()) {
	// ObjC autoreleasepools are thread-local, so we need the body to be locked to
	// the same OS thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	pool := C.objc_autoreleasePoolPush()
	defer C.objc_autoreleasePoolPop(pool)
	body()
}
