package core_test

import (
	"testing"

	"github.com/progrium/macdriver/core"
)

func TestNSDictionaryCGo(t *testing.T) {
	k := core.NSString_FromString("key")
	v := core.NSString_FromString("value")
	d := core.NSDictionary_Init(k, v)
	count := d.Count()
	if count != 1 {
		t.Errorf("expected count to be 1, but got: %d", count)
	}
}
