package objc

import (
	"runtime"
	"testing"
	"time"
)

type GCReleasedTest struct {
	Object `objc:"GCReleasedTest : NSObject"`
}

func TestGCReleased(t *testing.T) {
	c := NewClassFromStruct(GCReleasedTest{})
	deallocCount := 0
	c.AddMethod("dealloc", func(o Object) {
		deallocCount++
		o.SendSuper("dealloc")
	})
	RegisterClass(c)

	o := GCReleased(GetClass("GCReleasedTest").Alloc().Init())
	runtime.GC()
	if deallocCount != 0 {
		t.Errorf("did not expect dealloc to be called while object is still live, but got: %d", deallocCount)
	}
	runtime.KeepAlive(o)
	// Call the GC multiple time if necessary to be sure it has a chance to free
	// this memory.
	for i := 0; i < 10; i++ {
		runtime.GC()
		if deallocCount > 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if deallocCount != 1 {
		t.Errorf("expected dealloc to be called once, but got: %d", deallocCount)
	}
}
