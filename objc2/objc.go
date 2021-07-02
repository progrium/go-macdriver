package objc2

import (
	"unsafe"
)

/*
#cgo LDFLAGS: -lobjc
#define __OBJC2__ 1
#include <objc/message.h>
#include <stdlib.h>

void *getClass(char *name) {
	return (void *) objc_getClass(name);
}

typedef void* fn_type_NSObject_alloc(void*, SEL);

void* NSObject_alloc(void *id) {
	return ((fn_type_NSObject_alloc*) objc_msgSend) (id, sel_registerName("alloc"));
}
*/
import "C"

type Id uintptr

func (x Id) Id() Id { return x }

type Object interface {
	Id() Id
}

type Class struct {
	Object
}

func (c Class) Alloc() Id {
	return Id(C.NSObject_alloc(unsafe.Pointer(c.Id())))
}

func GetClass(name string) Class {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	id := Id(C.getClass(cName))
	return Class{id}
}

// https://developer.apple.com/documentation/objectivec/nsuinteger?language=objc
// typedef unsigned long NSUInteger;
type NSUInteger uint
type NSStringEncoding NSUInteger
type Unichar uint16
