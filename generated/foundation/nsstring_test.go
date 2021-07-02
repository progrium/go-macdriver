package foundation

import "testing"

func TestNSString(t *testing.T) {
	s := NSString_class.Alloc().Init()
	len := s.Length()
	if len != 0 {
		t.Errorf("expected len to be 0, but got %d", len)
	}
}

func TestNSString_fromGo(t *testing.T) {
	text := "foobar"
	s := NSString_FromString(text)
	size := s.Length()
	if int(size) != len(text) {
		t.Errorf("expected len to be %d, but got %d", len(text), size)
	}
	first := s.CharacterAtIndex_(0)
	if first != 'f' {
		t.Errorf("got unexpected first character: %s", string(rune(first)))
	}
}
