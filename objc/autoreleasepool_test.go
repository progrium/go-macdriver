package objc

import (
	"testing"
)

type AutoreleaseDealloc struct {
	Object       `objc:"AutoreleaseDealloc : NSObject"`
	deallocCount int
}

func TestAutorelease(t *testing.T) {
	c := NewClassFromStruct(AutoreleaseDealloc{})
	deallocCount := 0
	c.AddMethod("dealloc", func(o Object) {
		deallocCount++
		o.SendSuper("dealloc")
	})
	RegisterClass(c)

	AutoreleasePool(func() {
		dc := GetClass("AutoreleaseDealloc").Alloc().Init()
		dc.Autorelease()
	})
	if deallocCount != 1 {
		t.Errorf("expected dealloc to be called once, but got: %d", deallocCount)
	}
}
