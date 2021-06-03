// Copyright (c) 2012 The 'objc' Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package objc

import (
	"sync"
	"testing"
	"unsafe"
)

type SomeObject struct {
	Object `objc:"SomeObject : NSObject"`

	t *testing.T
}

const (
	float64val float64 = 42.0
	float32val float32 = 54.0
	uint64val  uint64  = 0x432fea3effca34e3
	uint32val  uint32  = 0x3214324e
	uint16val  uint16  = 0x465f
	uint8val   uint8   = 0xe3
	int64val   int64   = -4214321421423
	int32val   int32   = -43414324
	int16val   int16   = -600
	int8val    int8    = -120
	boolval    bool    = true
)

var once sync.Once

func registerTestClass() {
	once.Do(func() {
		c := NewClassFromStruct(SomeObject{})
		c.AddMethod("setTesting:", (*SomeObject).SetTesting)
		c.AddMethod("callWithObject:selector:", (*SomeObject).CallWithObjectAndSelector)
		c.AddMethod("callWithFloat64:", (*SomeObject).CallWithFloat64)
		c.AddMethod("callWithFloat32:", (*SomeObject).CallWithFloat32)
		c.AddMethod("callWithUint64:", (*SomeObject).CallWithUint64)
		c.AddMethod("callWithUint32:", (*SomeObject).CallWithUint32)
		c.AddMethod("callWithUint16:", (*SomeObject).CallWithUint16)
		c.AddMethod("callWithUint8:", (*SomeObject).CallWithUint8)
		c.AddMethod("callWithInt64:", (*SomeObject).CallWithInt64)
		c.AddMethod("callWithInt32:", (*SomeObject).CallWithInt32)
		c.AddMethod("callWithInt16:", (*SomeObject).CallWithInt16)
		c.AddMethod("callWithInt8:", (*SomeObject).CallWithInt8)
		c.AddMethod("callWithBool:", (*SomeObject).CallWithBool)
		c.AddMethod("callReturnObject:", (*SomeObject).CallReturnObject)
		RegisterClass(c)
	})
}

// func (so *SomeObject) SetTesting(t *testing.T) {
func (so *SomeObject) SetTesting(t uintptr) {
	// FIXME
	so.t = (*testing.T)(unsafe.Pointer(t))
}

func (so *SomeObject) CallWithObjectAndSelector(object Object, selector Selector) {
	if selector.Selector() != "callWithObject:selector:" {
		so.t.Errorf("unexpected selector")
	}
	if so.Pointer() != object.Pointer() {
		so.t.Errorf("unexpected object")
	}
}

func (so *SomeObject) CallWithFloat64(val float64) {
	if val != float64val {
		so.t.Errorf("float64: expected %v, got %v", float64val, val)
	}
}

func (so *SomeObject) CallWithFloat32(val float32) {
	if val != float32val {
		so.t.Errorf("float32: expected %v, got %v", float32val, val)
	}
}

func (so *SomeObject) CallWithUint64(val uint64) {
	if val != uint64val {
		so.t.Errorf("uint64: expected %v, got %v", uint64val, val)
	}
}

func (so *SomeObject) CallWithUint32(val uint32) {
	if val != uint32val {
		so.t.Errorf("uint32: expected %v, got %v", uint32val, val)
	}
}

func (so *SomeObject) CallWithUint16(val uint16) {
	if val != uint16val {
		so.t.Errorf("uint16: expected %v, got %v", uint16val, val)
	}
}

func (so *SomeObject) CallWithUint8(val uint8) {
	if val != uint8val {
		so.t.Errorf("uint8: expected %v, got %v", uint8val, val)
	}
}

func (so *SomeObject) CallWithInt64(val int64) {
	if val != int64val {
		so.t.Errorf("int64: expected %v, got %v", int64val, val)
	}
}

func (so *SomeObject) CallWithInt32(val int32) {
	if val != int32val {
		so.t.Errorf("int32: expected %v, got %v", int32val, val)
	}
}

func (so *SomeObject) CallWithInt16(val int16) {
	if val != int16val {
		so.t.Errorf("int8: expected %v, got %v", int16val, val)
	}
}

func (so *SomeObject) CallWithInt8(val int8) {
	if val != int8val {
		so.t.Errorf("int8: expected %v, got %v", int8val, val)
	}
}

func (so *SomeObject) CallWithBool(val bool) {
	if val != boolval {
		so.t.Errorf("bool: expected %v, got %v", boolval, val)
	}
}

func (so *SomeObject) CallReturnObject(o Object) Object {
	return o
}

func TestSelectorObjectPassing(t *testing.T) {
	registerTestClass()
	so := GetClass("SomeObject").Send("alloc").Send("init")
	so.Send("setTesting:", t)
	so.Send("callWithObject:selector:", so, GetSelector("callWithObject:selector:"))
}

func TestFloat64Passing(t *testing.T) {
	registerTestClass()
	so := GetClass("SomeObject").Send("alloc").Send("init")
	so.Send("setTesting:", t)
	so.Send("callWithFloat64:", float64val)
}

func TestFloat32Passing(t *testing.T) {
	registerTestClass()
	so := GetClass("SomeObject").Send("alloc").Send("init")
	so.Send("setTesting:", t)
	so.Send("callWithFloat32:", float32val)
}

func TestUint64Passing(t *testing.T) {
	registerTestClass()
	so := GetClass("SomeObject").Send("alloc").Send("init")
	so.Send("setTesting:", t)
	so.Send("callWithUint64:", uint64val)
}

func TestUint32Passing(t *testing.T) {
	registerTestClass()
	so := GetClass("SomeObject").Send("alloc").Send("init")
	so.Send("setTesting:", t)
	so.Send("callWithUint32:", uint32val)
}

func TestUint16Passing(t *testing.T) {
	registerTestClass()
	so := GetClass("SomeObject").Send("alloc").Send("init")
	so.Send("setTesting:", t)
	so.Send("callWithUint16:", uint16val)
}

func TestUint8Passing(t *testing.T) {
	registerTestClass()
	so := GetClass("SomeObject").Send("alloc").Send("init")
	so.Send("setTesting:", t)
	so.Send("callWithUint8:", uint8val)
}

func TestInt64Passing(t *testing.T) {
	registerTestClass()
	so := GetClass("SomeObject").Send("alloc").Send("init")
	so.Send("setTesting:", t)
	so.Send("callWithInt64:", int64val)
}

func TestInt32Passing(t *testing.T) {
	registerTestClass()
	so := GetClass("SomeObject").Send("alloc").Send("init")
	so.Send("setTesting:", t)
	so.Send("callWithInt32:", int32val)
}

func TestInt16Passing(t *testing.T) {
	registerTestClass()
	so := GetClass("SomeObject").Send("alloc").Send("init")
	so.Send("setTesting:", t)
	so.Send("callWithInt16:", int16val)
}

func TestInt8Passing(t *testing.T) {
	registerTestClass()
	so := GetClass("SomeObject").Send("alloc").Send("init")
	so.Send("setTesting:", t)
	so.Send("callWithInt8:", int8val)
}

func TestBoolPassing(t *testing.T) {
	registerTestClass()
	so := GetClass("SomeObject").Send("alloc").Send("init")
	so.Send("setTesting:", t)
	so.Send("callWithBool:", boolval)
}

func TestObjectReturn(t *testing.T) {
	registerTestClass()
	so := GetClass("SomeObject").Send("alloc").Send("init")
	out := so.Send("callReturnObject:", so)
	if out != so {
		t.Errorf("expected to return the input object")
	}
}

func TestObjectReturn2(t *testing.T) {
	registerTestClass()
	so := GetClass("SomeObject").Send("alloc").Send("init")
	out := send2(so.Pointer(), "callReturnObject:", so)
	if out != so.Pointer() {
		t.Errorf("expected to return the input object")
	}
}
