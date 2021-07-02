package foundation

// #include <stdlib.h>
import "C"

import (
	"unsafe"

	"github.com/progrium/macdriver/objc2"
)

const (
	NSUTF8StringEncoding = 4
)

func NSString_FromString(s string) NSString {
	b := []byte(s)
	c := C.CBytes(b)
	defer C.free(unsafe.Pointer(c))
	ret := NSString_class.Alloc().InitWithBytes_length_encoding_(c, objc2.NSUInteger(len(b)), NSUTF8StringEncoding)
	return ret
}
