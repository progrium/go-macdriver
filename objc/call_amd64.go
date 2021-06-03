// Copyright (c) 2012 The 'objc' Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package objc

/*
extern unsigned long GoObjc_CallTargetFrameSetup;
*/
import "C"
import (
	"math"
	"reflect"
	"unsafe"
)

// amd64frame represents the layout of the
// register set once it's pushed onto the stack
// when the Objective-C runtime calls our call target.
type amd64frame struct {
	rdi  uintptr
	rsi  uintptr
	rdx  uintptr
	rcx  uintptr
	r8   uintptr
	r9   uintptr
	xmm0 uintptr
	xmm1 uintptr
	xmm2 uintptr
	xmm3 uintptr
	xmm4 uintptr
	xmm5 uintptr
	xmm6 uintptr
	xmm7 uintptr
}

// amd64frameFetcher implements the logic needed
// to fetch arguments from an amd64frame in the
// correct order.
type amd64frameFetcher struct {
	frame  *amd64frame
	ints   *[6]uintptr
	floats *[8]uintptr
	stack  *[10]uintptr
	ioff   int
	foff   int
	soff   int
}

// frameFetcher returns a new amd64frameFetcher that
// wraps an existing amd64 frame.
func frameFetcher(frame *amd64frame) amd64frameFetcher {
	ints := (*[6]uintptr)(unsafe.Pointer(frame))
	floats := (*[8]uintptr)(unsafe.Pointer(&frame.xmm0))
	stack := (*[10]uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&frame.xmm7)) + 8))
	return amd64frameFetcher{
		ints:   ints,
		floats: floats,
		stack:  stack,
	}
}

// Int returns the next integer argument from the amd64frame
// wrapped by the frame fetcher.
func (ff *amd64frameFetcher) Int() uintptr {
	if ff.ioff < len(ff.ints) {
		val := ff.ints[ff.ioff]
		ff.ioff++
		return val
	}
	return ff.Stack()
}

// Float returns the next floating point argument from the amd64
// frame wrapped by the frame fetcher.
func (ff *amd64frameFetcher) Float() uintptr {
	if ff.foff < len(ff.floats) {
		val := ff.floats[ff.foff]
		ff.foff++
		return val
	}
	return ff.Stack()
}

// Stack returns the next stack argument from amd64
// frame wrapped by the frame fetcher.
func (ff *amd64frameFetcher) Stack() uintptr {
	if ff.soff < len(ff.stack) {
		val := ff.stack[ff.soff]
		ff.soff++
		return val
	}
	panic("call: argument list exhausted")
}

// methodCallTarget returns a pointer to the entry point
// that the Objective-C runtime must call to reach an
// exported method from Go.
func methodCallTarget() unsafe.Pointer {
	return unsafe.Pointer(&C.GoObjc_CallTargetFrameSetup)
}

//export goMethodCallEntryPoint
func goMethodCallEntryPoint(p uintptr) uintptr {
	frame := (*amd64frame)(unsafe.Pointer(p))
	fetcher := frameFetcher(frame)

	obj := object{ptr: fetcher.Int()}
	sel := stringFromSelector(unsafe.Pointer(fetcher.Int()))

	clsName := object{ptr: getObjectClass(obj).Pointer()}.className()
	clsInfo := classMap[clsName]
	method := clsInfo.MethodForSelector(sel)

	// Check if we have an internal pointer set for this object.
	// If not, make it happen.
	internalPtr := obj.internalPointer()
	if internalPtr == nil {
		// Allocate the Go struct.
		val := reflect.New(clsInfo.typ)
		ptr := unsafe.Pointer(val.Pointer())
		// Add a reference in the classInfo. This ensures we have
		// a reference to our Go struct somewhere in Go land, which
		// makes the garbage collector not collect it under our feet.
		clsInfo.AddRef(ptr)
		// Set the internalPointer so we can easily access
		// the instance's Go struct pointer.
		obj.setInternalPointer(ptr)
		internalPtr = ptr
		// Finally, update the Go struct's embedded objc.Object to
		// point to the actual Objective-C instance.
		structVal := val.Elem()
		if structVal.Kind() == reflect.Struct {
			objectVal := structVal.FieldByName("Object")
			if objectVal.IsValid() {
				objectVal.Set(reflect.ValueOf(obj))
			}
		}
	}
	objVal := reflect.NewAt(clsInfo.typ, internalPtr)

	// Check if the invoked selector is a setter for IBOutlets.
	if _, isSetter := clsInfo.setters[sel]; isSetter {
		valuePtr := fetcher.Int()
		keyName := sel[3:]                    // strip 'set'
		keyName = keyName[0 : len(keyName)-1] // strip ':'
		setIBOutletValue(objVal, keyName, object{ptr: valuePtr})
		return 0
	}

	// Our own internal override for setValue:forKey: in order
	// to support key-value coding. (For IBOutlets on iOS)
	if sel == "setValue:forKey:" && method == nil {
		// We only support Object values, so fetching
		// Ints here is OK.
		valuePtr := fetcher.Int()
		keyPtr := fetcher.Int()

		// We don't export any NSString-based functionality
		// in package objc, except for the String() method
		// on object.  It calls the object's decription method,
		// which for NSStrings returns the string itself (or at
		// least something that has been good enough for now!).
		keyName := object{ptr: keyPtr}.String()

		setIBOutletValue(objVal, keyName, object{ptr: valuePtr})
		return 0
	}

	// The default dealloc implementation. This is called
	// if a class doesn't register its own custom dealloc
	// method.
	if sel == "dealloc" && method == nil {
		clsInfo.RemoveRef(internalPtr)
		obj.SendSuper("dealloc")
		return 0
	}

	methodVal := reflect.ValueOf(method)

	// First argument should point to the Go method's proper receiver.
	// That's stored in the internalPointer, so fetch that.
	args := []reflect.Value{objVal}

	// Take care of the rest of the arguments
	mt := reflect.TypeOf(method)
	for i := 1; i < mt.NumIn(); i++ {
		typ := mt.In(i)

		if typ.Implements(objectInterfaceType) {
			args = append(args, reflect.ValueOf(object{ptr: fetcher.Int()}))
			continue
		} else if typ.Implements(selectorInterfaceType) {
			sel := selector(stringFromSelector(unsafe.Pointer(fetcher.Int())))
			args = append(args, reflect.ValueOf(sel))
			continue
		}

		switch typ.Kind() {
		case reflect.Int:
			args = append(args, reflect.ValueOf(int(fetcher.Int())))
		case reflect.Int8:
			args = append(args, reflect.ValueOf(int8(fetcher.Int())))
		case reflect.Int16:
			args = append(args, reflect.ValueOf(int16(fetcher.Int())))
		case reflect.Int32:
			args = append(args, reflect.ValueOf(int32(fetcher.Int())))
		case reflect.Int64:
			args = append(args, reflect.ValueOf(int64(fetcher.Int())))

		case reflect.Uint8:
			args = append(args, reflect.ValueOf(uint8(fetcher.Int())))
		case reflect.Uint16:
			args = append(args, reflect.ValueOf(uint16(fetcher.Int())))
		case reflect.Uint32:
			args = append(args, reflect.ValueOf(uint32(fetcher.Int())))
		case reflect.Uint64:
			args = append(args, reflect.ValueOf(uint64(fetcher.Int())))
		case reflect.Uintptr:
			args = append(args, reflect.ValueOf(fetcher.Int()))

		case reflect.Float32:
			args = append(args, reflect.ValueOf(math.Float32frombits(uint32(fetcher.Float()))))
		case reflect.Float64:
			args = append(args, reflect.ValueOf(math.Float64frombits(uint64(fetcher.Float()))))

		case reflect.Bool:
			val := fetcher.Int() != 0
			args = append(args, reflect.ValueOf(val))

		case reflect.Ptr:
			ptrAddr := unsafe.Pointer(uintptr(fetcher.Int()))
			args = append(args, reflect.NewAt(typ.Elem(), ptrAddr))

		default:
			panic("call: unhandled arg")
		}
	}

	retVals := methodVal.Call(args)

	// If a custom dealloc method has been registered, we
	// still need to remove the reference to our Go struct
	// pointer in order for the GC to kick in. Do that now.
	if sel == "dealloc" {
		clsInfo.RemoveRef(internalPtr)
	}

	if len(retVals) > 0 {
		val := retVals[0]
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return uintptr(val.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return uintptr(val.Uint())
		case reflect.Bool:
			if val.Bool() {
				return 1
			} else {
				return 0
			}
		case reflect.Float32:
			frame.xmm0 = uintptr(math.Float32bits(float32(val.Float())))
			return 1
		case reflect.Float64:
			frame.xmm0 = uintptr(math.Float64bits(val.Float()))
			return 1
		case reflect.Interface:
			if obj, ok := val.Interface().(Object); ok {
				return obj.Pointer()
			}
			panic("call: bad interface return value")
		default:
			panic("call: unknown return value")
		}
	}

	return 0
}
