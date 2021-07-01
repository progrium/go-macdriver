package foundation

import (
	"unsafe"

	objc "github.com/progrium/macdriver/objc2"
)

/*
#cgo LDFLAGS: -lobjc -framework Foundation
#define __OBJC2__ 1
#include <objc/message.h>
#include <stdlib.h>


typedef void* fn_type_NSString_init(void*, SEL);

void* NSString_init(void *id) {
	return ((fn_type_NSString_init*) objc_msgSend) (id, sel_registerName("init"));
}

typedef void* fn_type_NSString_initWithBytes_length_encoding_(void*, SEL, void*, unsigned long, unsigned long);

void* NSString_initWithBytes_length_encoding_(void *id, void* bytes, unsigned long len, unsigned long encoding) {
	return ((fn_type_NSString_initWithBytes_length_encoding_*) objc_msgSend) (id, sel_registerName("initWithBytes:length:encoding:"), bytes, len, encoding);
}

typedef unsigned long fn_type_NSString_length(void*, SEL);

unsigned long NSString_length(void *id) {
	return ((fn_type_NSString_length*) objc_msgSend) (id, sel_registerName("length"));
}

typedef unsigned short fn_type_NSString_characterAtIndex_(void*, SEL, unsigned long);

unsigned short NSString_characterAtIndex_(void *id, unsigned long index) {
	return ((fn_type_NSString_characterAtIndex_*) objc_msgSend) (id, sel_registerName("characterAtIndex:"), index);
}

*/
import "C"

func NSString_init(
	self objc.Object,
) (

	r0 NSString,

) {
	ret := C.NSString_init(
		unsafe.Pointer(self.Id()),
	)

	r0 = NSString_fromPointer(ret)

	return
}

func NSString_initWithBytes_length_encoding_(
	self objc.Object,
	bytes unsafe.Pointer,
	len objc.NSUInteger,
	encoding objc.NSStringEncoding,
) (

	r0 NSString,

) {
	ret := C.NSString_initWithBytes_length_encoding_(
		unsafe.Pointer(self.Id()),
		bytes,
		C.ulong(len),
		C.ulong(encoding),
	)

	r0 = NSString_fromPointer(ret)

	return
}

func NSString_length(
	self objc.Object,
) (

	r0 objc.NSUInteger,

) {
	ret := C.NSString_length(
		unsafe.Pointer(self.Id()),
	)

	r0 = objc.NSUInteger(ret)

	return
}

func NSString_characterAtIndex_(
	self objc.Object,
	index objc.NSUInteger,
) (

	r0 objc.Unichar,

) {
	ret := C.NSString_characterAtIndex_(
		unsafe.Pointer(self.Id()),
		C.ulong(index),
	)

	r0 = objc.Unichar(ret)

	return
}

type NSString_class_methods struct {
	objc.Class
}

var NSString_class = NSString_class_methods{
	Class: objc.GetClass("NSString"),
}

func (cls NSString_class_methods) Alloc() NSString {
	return NSString_fromId(cls.Class.Alloc())
}

type NSString struct {
	objc.Object
}

func NSString_fromId(id objc.Id) NSString {
	return NSString{id}
}

func NSString_fromPointer(ptr unsafe.Pointer) NSString {
	return NSString_fromId(objc.Id(ptr))
}

func (x NSString) Init() NSString {
	return NSString_init(
		x,
	)
}

func (x NSString) InitWithBytes_length_encoding_(

	bytes unsafe.Pointer,

	len objc.NSUInteger,

	encoding objc.NSStringEncoding,

) NSString {
	return NSString_initWithBytes_length_encoding_(
		x,

		bytes,

		len,

		encoding,
	)
}

func (x NSString) Length() objc.NSUInteger {
	return NSString_length(
		x,
	)
}

func (x NSString) CharacterAtIndex_(

	index objc.NSUInteger,

) objc.Unichar {
	return NSString_characterAtIndex_(
		x,

		index,
	)
}
