// Copyright (c) 2012 The 'objc' Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package objc

/*
#cgo LDFLAGS: -lobjc
#define __OBJC2__ 1
#include <objc/message.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"log"
	"reflect"
	"unsafe"
)

// setIBOutletValue attempts to assign the Objective-C object represented
// by 'value' to the field named 'name' in the Go struct represented by 'obj'.
//
// The function sends the 'retain' message to the Objective-C object represented
// by 'value' if the assignment was successful.
//
// If the assignment operation fails, this function will raise a runtime panic.
func setIBOutletValue(obj reflect.Value, name string, value Object) {
	// Find an IBOutlet with the name keyName.
	val := obj.Elem()
	typ := val.Type()
	fieldIdx := -1
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Tag.Get("objc") == "IBOutlet" && field.Type.Implements(objectInterfaceType) {
			if field.Name == name {
				fieldIdx = i
				break
			}
		}
	}

	// If we couldn't find a matching field, simply panic.
	// We should only run into this is there is a bug in the package.
	if fieldIdx == -1 {
		panic("objc: bad setter for IBOutlet field '" + name + "'")
	}

	fieldVal := val.Field(fieldIdx)
	fieldVal.Set(reflect.ValueOf(value))
	value.Retain()
}

//export goMethodSignature
func goMethodSignature(self, cmd, selPtr unsafe.Pointer) unsafe.Pointer {
	sel := stringFromSelector(selPtr)
	log.Printf("resolving %v", sel)

	obj := object{ptr: uintptr(self)}

	clsName := object{ptr: getObjectClass(obj).Pointer()}.className()
	clsInfo := classMap[clsName]
	method := clsInfo.MethodForSelector(sel)
	if method == nil {
		// FIXME delegate to super class?
		return nil
	}
	return unsafe.Pointer(newSignature(funcTypeInfo(method)))
}

//export goInvoke
func goInvoke(self, cmd, invocation unsafe.Pointer) {
	inv := ObjectPtr(uintptr(invocation))
	// Or use "self" here?
	obj := object{ptr: inv.Get("target").Pointer()}
	sel := stringFromSelector(unsafe.Pointer(inv.Get("selector").Pointer()))

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
		var valuePtr uintptr
		inv.Send("getArgument:atIndex:", unsafe.Pointer(&valuePtr), 2)
		keyName := sel[3:]                    // strip 'set'
		keyName = keyName[0 : len(keyName)-1] // strip ':'
		setIBOutletValue(objVal, keyName, object{ptr: valuePtr})
		// TODO setReturnValue:?
		return
	}

	// Our own internal override for setValue:forKey: in order
	// to support key-value coding. (For IBOutlets on iOS)
	if sel == "setValue:forKey:" && method == nil {
		// We only support Object values, so fetching
		// Ints here is OK.
		var valuePtr, keyPtr uintptr
		inv.Send("getArgument:atIndex:", unsafe.Pointer(&valuePtr), 2)
		inv.Send("getArgument:atIndex:", unsafe.Pointer(&keyPtr), 3)

		// We don't export any NSString-based functionality
		// in package objc, except for the String() method
		// on object.  It calls the object's decription method,
		// which for NSStrings returns the string itself (or at
		// least something that has been good enough for now!).
		keyName := object{ptr: keyPtr}.String()

		setIBOutletValue(objVal, keyName, object{ptr: valuePtr})
		// TODO setReturnValue:?
		return
	}

	// The default dealloc implementation. This is called
	// if a class doesn't register its own custom dealloc
	// method.
	if sel == "dealloc" && method == nil {
		clsInfo.RemoveRef(internalPtr)
		obj.SendSuper("dealloc")
		// TODO setReturnValue:?
		return
	}

	methodVal := reflect.ValueOf(method)

	// First argument should point to the Go method's proper receiver.
	// That's stored in the internalPointer, so fetch that.
	args := []reflect.Value{objVal}

	// Take care of the rest of the arguments
	mt := reflect.TypeOf(method)
	for i := 1; i < mt.NumIn(); i++ {
		fetchArg := func(p unsafe.Pointer) {
			inv.Send("getArgument:atIndex:", p, i+1)
		}
		typ := mt.In(i)

		if typ.Implements(objectInterfaceType) {
			var obj object
			fetchArg(unsafe.Pointer(&obj.ptr))
			args = append(args, reflect.ValueOf(obj))
			continue
		} else if typ.Implements(selectorInterfaceType) {
			var v uintptr
			fetchArg(unsafe.Pointer(&v))
			sel := selector(stringFromSelector(unsafe.Pointer(v)))
			args = append(args, reflect.ValueOf(sel))
			continue
		}

		switch typ.Kind() {
		case reflect.Int:
			var v int
			fetchArg(unsafe.Pointer(&v))
			args = append(args, reflect.ValueOf(v))
		case reflect.Int8:
			var v int8
			fetchArg(unsafe.Pointer(&v))
			args = append(args, reflect.ValueOf(v))
		case reflect.Int16:
			var v int16
			fetchArg(unsafe.Pointer(&v))
			args = append(args, reflect.ValueOf(v))
		case reflect.Int32:
			var v int32
			fetchArg(unsafe.Pointer(&v))
			args = append(args, reflect.ValueOf(v))
		case reflect.Int64:
			var v int64
			fetchArg(unsafe.Pointer(&v))
			args = append(args, reflect.ValueOf(v))

		case reflect.Uint8:
			var v uint8
			fetchArg(unsafe.Pointer(&v))
			args = append(args, reflect.ValueOf(v))
		case reflect.Uint16:
			var v uint16
			fetchArg(unsafe.Pointer(&v))
			args = append(args, reflect.ValueOf(v))
		case reflect.Uint32:
			var v uint32
			fetchArg(unsafe.Pointer(&v))
			args = append(args, reflect.ValueOf(v))
		case reflect.Uint64:
			var v uint64
			fetchArg(unsafe.Pointer(&v))
			args = append(args, reflect.ValueOf(v))
		case reflect.Uintptr:
			var v uintptr
			fetchArg(unsafe.Pointer(&v))
			args = append(args, reflect.ValueOf(v))

		case reflect.Float32:
			var v float32
			fetchArg(unsafe.Pointer(&v))
			args = append(args, reflect.ValueOf(v))
		case reflect.Float64:
			var v float64
			fetchArg(unsafe.Pointer(&v))
			args = append(args, reflect.ValueOf(v))

		case reflect.Bool:
			var v int // FIXME what size should we decode bool into?
			fetchArg(unsafe.Pointer(&v))
			args = append(args, reflect.ValueOf(v != 0))

		case reflect.Ptr:
			var v uintptr
			fetchArg(unsafe.Pointer(&v))
			args = append(args, reflect.NewAt(typ.Elem(), unsafe.Pointer(v)))

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
		case reflect.Int:
			v := int(val.Int())
			inv.Send("setReturnValue:", uintptr(unsafe.Pointer(&v)))
		case reflect.Int8:
			v := int8(val.Int())
			inv.Send("setReturnValue:", uintptr(unsafe.Pointer(&v)))
		case reflect.Int16:
			v := int16(val.Int())
			inv.Send("setReturnValue:", uintptr(unsafe.Pointer(&v)))
		case reflect.Int32:
			v := int32(val.Int())
			inv.Send("setReturnValue:", uintptr(unsafe.Pointer(&v)))
		case reflect.Int64:
			v := int64(val.Int())
			inv.Send("setReturnValue:", uintptr(unsafe.Pointer(&v)))
		case reflect.Uint:
			v := uint(val.Int())
			inv.Send("setReturnValue:", uintptr(unsafe.Pointer(&v)))
		case reflect.Uint8:
			v := uint8(val.Int())
			inv.Send("setReturnValue:", uintptr(unsafe.Pointer(&v)))
		case reflect.Uint16:
			v := uint16(val.Int())
			inv.Send("setReturnValue:", uintptr(unsafe.Pointer(&v)))
		case reflect.Uint32:
			v := uint32(val.Int())
			inv.Send("setReturnValue:", uintptr(unsafe.Pointer(&v)))
		case reflect.Uint64:
			v := uint64(val.Int())
			inv.Send("setReturnValue:", uintptr(unsafe.Pointer(&v)))
		case reflect.Uintptr:
			v := uintptr(val.Int())
			inv.Send("setReturnValue:", uintptr(unsafe.Pointer(&v)))
		case reflect.Bool:
			v := val.Bool()
			inv.Send("setReturnValue:", uintptr(unsafe.Pointer(&v)))
		case reflect.Float32:
			v := float32(val.Float())
			inv.Send("setReturnValue:", uintptr(unsafe.Pointer(&v)))
		case reflect.Float64:
			v := val.Float()
			inv.Send("setReturnValue:", uintptr(unsafe.Pointer(&v)))
		case reflect.Interface:
			obj, ok := val.Interface().(Object)
			if !ok {
				panic(fmt.Errorf("call: bad interface return value: %v", val.Type()))
			}
			v := uintptr(obj.Pointer())
			inv.Send("setReturnValue:", uintptr(unsafe.Pointer(&v)))
		default:
			panic("call: unknown return value")
		}
	}
}
