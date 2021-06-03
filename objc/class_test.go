package objc

import (
	"log"
	"testing"
	"unsafe"
)

type DynamicMethodStruct struct {
	Object `objc:"DynamicMethodStruct : NSObject"`
	value  uint64
}

func (gs *DynamicMethodStruct) SetValue() {
	gs.value = 0xbadc0c0a
}

// - (NSMethodSignature *)methodSignatureForSelector:(SEL)aSelector;
func (gs *DynamicMethodStruct) MethodSignatureForSelector(selector uintptr) uintptr {
	// sel := ObjectPtr(selector)
	name := stringFromSelector(unsafe.Pointer(selector))
	if name == "foo:" {
		return newSignature("v@:@@")
	}
	log.Printf("methodSig for %v", name)
	// gs.Send("doesNotRecognizeSelector:", selector)
	return 0
}

// - (void)forwardInvocation:(NSInvocation *)anInvocation;
func (gs *DynamicMethodStruct) ForwardInvocation(anInvocation uintptr) {
	inv := ObjectPtr(anInvocation)
	var arg0 int
	inv.Send("getArgument:atIndex:", uintptr(unsafe.Pointer(&arg0)), 2)
	log.Printf("forwardInvocation: %v", inv)
	log.Printf("arg0: %v", arg0)
}

func (gs *DynamicMethodStruct) Value() uint64 {
	return gs.value
}

func TestDynamicMethod(t *testing.T) {
	t.Skipf("FIXME")

	c := NewClassFromStruct(DynamicMethodStruct{})
	c.AddMethod("setValue", (*DynamicMethodStruct).SetValue)
	c.AddMethod("hasValue", (*DynamicMethodStruct).Value)
	// c.AddMethod("methodSignatureForSelector:", (*DynamicMethodStruct).MethodSignatureForSelector)
	// c.AddMethod("forwardInvocation:", (*DynamicMethodStruct).ForwardInvocation)
	RegisterClass(c)

	o := GetClass("DynamicMethodStruct").Send("alloc").Send("init")
	// out := o.Send("methodSignatureForSelector:", uintptr(RegisterSelector("foo:bar:")))
	// t.Errorf("out: %v", out.Pointer())
	o.Send("foo:", 42)
	t.Fail()
	// o.Send("blah:blah:")
}
