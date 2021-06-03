package objc

/*
#cgo LDFLAGS: -lobjc
#define __OBJC2__ 1
#include <objc/message.h>
#include <stdlib.h>

// + (NSInvocation *)invocationWithMethodSignature:(NSMethodSignature *)sig;
id (* send_NSInvocation_invocationWithMethodSignature)(Class cls, SEL _cmd, id sig) = (id(*)(Class,SEL,id))objc_msgSend;

// - (NSMethodSignature *)methodSignatureForSelector:(SEL)aSelector;
id (* send_NSObject_methodSignatureForSelector)(id self, SEL _cmd, SEL aSelector) = (id(*)(id,SEL,SEL))objc_msgSend;

void (* send_NSInvocation_setTarget)(id self, SEL _cmd, id target) = (void(*)(id,SEL,id))objc_msgSend;
void (* send_NSInvocation_setSelector)(id self, SEL _cmd, SEL sel) = (void(*)(id,SEL,SEL))objc_msgSend;

// - (void)setArgument:(void *)argumentLocation atIndex:(NSInteger)idx;
void (* send_setArgument)(id self, SEL _cmd, void *arg, int idx) = (void(*)(id,SEL,void*,int))objc_msgSend;
void setArgument(void *invocation, void *arg, int idx) {
	send_setArgument(invocation, sel_registerName("setArgument:atIndex:"), arg, idx);
}

// - (void)invoke;
void (* send_invoke)(id self, SEL _cmd) = (void(*)(id,SEL))objc_msgSend;

// - (void)getReturnValue:(void *)retLoc;
void (* send_getReturnValue)(id self, SEL _cmd, void *retLoc) = (void(*)(id,SEL,void*))objc_msgSend;

void *methodSignatureForSelector(void *target, char *selName) {
	return send_NSObject_methodSignatureForSelector(target, sel_registerName("methodSignatureForSelector:"), sel_registerName(selName));
}

int (* send_numberOfArguments)(id self, SEL _cmd) = (int(*)(id,SEL))objc_msgSend;

int numberOfArguments(void *target) {
	return send_numberOfArguments(target, sel_registerName("numberOfArguments"));
}

const char * (* send_methodReturnType)(id self, SEL _cmd) = (const char * (*)(id,SEL))objc_msgSend;

const char * methodReturnType(void *target) {
	return send_methodReturnType(target, sel_registerName("methodReturnType"));
}

id (* send_init)(id self, SEL _cmd) = (id(*)(id,SEL))objc_msgSend;

void *init(void *target) {
	return send_init(target, sel_registerName("init"));
}

void (* send_dealloc)(id self, SEL _cmd) = (void(*)(id,SEL))objc_msgSend;

void dealloc(void *target) {
	return send_dealloc(target, sel_registerName("dealloc"));
}

const char * (* send_getArgumentTypeAtIndex)(id self, SEL _cmd, int i) = (const char * (*)(id,SEL,int))objc_msgSend;

const char * getArgumentTypeAtIndex(void *target, int idx) {
	return send_getArgumentTypeAtIndex(target, sel_registerName("getArgumentTypeAtIndex:"), idx);
}

void *nsInvocation(void *target, char *selName) {
	Class cls = objc_getClass("NSInvocation");
	SEL sel = sel_registerName(selName);
	id sig = send_NSObject_methodSignatureForSelector(target, sel_registerName("methodSignatureForSelector:"), sel);
	id inv = send_NSInvocation_invocationWithMethodSignature(cls, sel_registerName("invocationWithMethodSignature:"), sig);
	send_NSInvocation_setTarget(inv, sel_registerName("setTarget:"), target);
	send_NSInvocation_setSelector(inv, sel_registerName("setSelector:"), sel);
	return inv;
}

void invoke(void *invocation, void *retLoc) {
	send_invoke(invocation, sel_registerName("invoke"));
	if (retLoc != NULL) {
		send_getReturnValue(invocation, sel_registerName("getReturnValue:"), retLoc);
	}
}
*/
import "C"
import (
	"fmt"
	"log"
	"math"
	"reflect"
	"unsafe"
)

func numberOfArguments(id uintptr) int {
	return int(C.numberOfArguments(unsafe.Pointer(id)))
}

func methodReturnType(id uintptr) string {
	return C.GoString(C.methodReturnType(unsafe.Pointer(id)))
}

func getArgumentTypeAtIndex(id uintptr, idx int) string {
	return C.GoString(C.getArgumentTypeAtIndex(unsafe.Pointer(id), C.int(idx)))
}

func methodSignatureForSelector(target uintptr, selName string) uintptr {
	cSelName := C.CString(selName)
	sig := uintptr(C.methodSignatureForSelector(unsafe.Pointer(target), cSelName))
	C.free(unsafe.Pointer(cSelName))
	return sig
}

func newInvocation(target uintptr, selName string) uintptr {
	return uintptr(C.nsInvocation(unsafe.Pointer(target), C.CString(selName)))
}

func setArgumentAtIndex(call uintptr, arg unsafe.Pointer, index int) {
	C.setArgument(unsafe.Pointer(call), arg, C.int(index))
}

// TODO variant for calls with no return?
// or just split into two calls?
func invoke(call uintptr, dest unsafe.Pointer) {
	C.invoke(unsafe.Pointer(call), dest)
}

func cstring(ptr uintptr) string {
	return C.GoString((*C.char)(unsafe.Pointer(ptr)))
}

func send(retDest unsafe.Pointer, target uintptr, selName string, args ...unsafe.Pointer) {
	inv := newInvocation(target, selName)
	for i, arg := range args {
		setArgumentAtIndex(inv, arg, i+2)
	}
	invoke(inv, retDest)
}

func send2(target uintptr, selName string, args ...interface{}) uintptr {
	switch selName {
	case "init":
		// TODO check args length is 0?
		// the method signature for init appears to return NULL which causes problems
		// maybe there's a better way, but we can hard-code the msgSend wrapper for
		// this case
		return uintptr(C.init(unsafe.Pointer(target)))
	case "dealloc":
		C.dealloc(unsafe.Pointer(target))
		return 0
	}
	sig := methodSignatureForSelector(target, selName)
	retType := methodReturnType(sig)
	log.Printf("%d %s -> %s", target, selName, retType)
	inv := newInvocation(target, selName)
	for i, arg := range args {
		setArgumentAtIndex(inv, wrapArg(arg), i+2)
	}
	// FIXME clean up some of the special cases here based on normalizing the
	// return type detection.
	switch retType {
	case "":
		invoke(inv, nil)
		return 0
	case encId, encSelector, encBool:
		var v uintptr
		invoke(inv, unsafe.Pointer(&v))
		return v
	case encVoid, "V" + encVoid: // "V" = one-way, used by "release"
		invoke(inv, nil)
		return 0
	case "r*": // const char *
		var v uintptr
		invoke(inv, unsafe.Pointer(&v))
		return v
	case encFloat:
		var v float32
		invoke(inv, unsafe.Pointer(&v))
		return uintptr(math.Float32bits(v))
	case encDouble:
		var v float64
		invoke(inv, unsafe.Pointer(&v))
		return uintptr(math.Float64bits(v))
	default:
		panic(fmt.Errorf("unhandled: %v", retType))
	}
}

func wrapArg(v interface{}) unsafe.Pointer {
	switch v := v.(type) {
	case Object:
		return wrapArg(v.Pointer())
	case selector:
		return wrapArg(RegisterSelector(v.Selector()))
	}

	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Int:
		v := int(val.Int())
		return unsafe.Pointer(&v)
	case reflect.Int8:
		v := int8(val.Int())
		return unsafe.Pointer(&v)
	case reflect.Int16:
		v := int16(val.Int())
		return unsafe.Pointer(&v)
	case reflect.Int32:
		v := int32(val.Int())
		return unsafe.Pointer(&v)
	case reflect.Int64:
		v := int64(val.Int())
		return unsafe.Pointer(&v)

	case reflect.Uint8:
		v := uint8(val.Uint())
		return unsafe.Pointer(&v)
	case reflect.Uint16:
		v := uint16(val.Uint())
		return unsafe.Pointer(&v)
	case reflect.Uint32:
		v := uint32(val.Uint())
		return unsafe.Pointer(&v)
	case reflect.Uint64:
		v := uint64(val.Uint())
		return unsafe.Pointer(&v)
	case reflect.Uintptr:
		v := uintptr(val.Uint())
		return unsafe.Pointer(&v)

	case reflect.Float32:
		v := float32(val.Float())
		return unsafe.Pointer(&v)
	case reflect.Float64:
		v := float64(val.Float())
		return unsafe.Pointer(&v)

	case reflect.Bool:
		v := val.Bool() // FIXME what size should we decode bool into?
		return unsafe.Pointer(&v)

	case reflect.Ptr, reflect.UnsafePointer:
		return wrapArg(val.Pointer())

	default:
		panic(fmt.Errorf("call: unhandled arg: %T %#v", v, v))
	}
}
