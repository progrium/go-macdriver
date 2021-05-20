package objc

/*
#cgo LDFLAGS: -lobjc -framework Foundation
extern void objc_autoreleasePoolPop(void* pool);
extern void* objc_autoreleasePoolPush(void);
*/
import "C"

func AutoreleasePool(body func()) {
	pool := C.objc_autoreleasePoolPush()
	defer C.objc_autoreleasePoolPop(pool)
	body()
}
