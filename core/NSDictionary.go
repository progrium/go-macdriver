package core

import (
	"unsafe"

	"github.com/progrium/macdriver/objc"
)

/*
#cgo LDFLAGS: -lobjc
#define __OBJC2__ 1
#include <objc/message.h>

int (* send_NSDictionary_count)(id self, SEL _cmd) = (int(*)(id,SEL))objc_msgSend;

int NSDictionary_count(void *id) {
	return send_NSDictionary_count(id, (void *) sel_registerName("count"));
}
*/
import "C"

type NSDictionary struct {
	objc.Object
}

var nsDictionary = objc.Get("NSDictionary")

func NSDictionary_New() NSDictionary {
	return NSDictionary{nsDictionary.Alloc().Init()}
}

func NSDictionary_Init(valueKeys ...interface{}) NSDictionary {
	return NSDictionary{nsDictionary.Alloc().Send("initWithObjectsAndKeys:", valueKeys...)}
}

func (d NSDictionary) ObjectForKey(key objc.Object) objc.Object {
	return d.Send("objectForKey:", key)
}

func (d NSDictionary) Count() int {
	x := C.NSDictionary_count(unsafe.Pointer(d.Pointer()))
	return int(x)
}
