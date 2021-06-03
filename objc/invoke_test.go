package objc

import (
	"testing"
	"unsafe"
)

func TestInvocation(t *testing.T) {
	rect := NSRect{
		NSPoint{
			X: 0,
			Y: 0,
		},
		NSSize{
			Width:  500,
			Height: 500,
		},
	}

	var objID uintptr
	send(unsafe.Pointer(&objID), GetClass("NSValue").Pointer(), "valueWithRect:", unsafe.Pointer(&rect))

	var out uintptr
	send(unsafe.Pointer(&out), objID, "description")
	actual := ObjectPtr(out).String()
	expected := "NSRect: {{0, 0}, {500, 500}}"
	if actual != expected {
		t.Errorf("expected %v, but got: %v", expected, actual)
	}

	sigID := methodSignatureForSelector(objID, "description")
	var argType uintptr
	idx := 1
	send(unsafe.Pointer(&argType), sigID, "getArgumentTypeAtIndex:", unsafe.Pointer(&idx))
	actual = cstring(argType)
	expected = ":" // arg 1 is a selector, encoded as a colon
	if actual != expected {
		t.Errorf("expected %v, but got: %v", expected, actual)
	}
}
