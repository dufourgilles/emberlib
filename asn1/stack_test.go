package asn1_test

import (
	"testing"

	. "github.com/dufourgilles/emberlib/asn1"
)

func TestStack(t *testing.T) {
	var expect = int(77)
	s := NewStack()
	s.Push(expect)
	v := s.Pop()
	if v.(int) != expect {
		t.Errorf("Invalid response %d. Expected %d", v.(int), expect)
	}
	v = s.Pop()
	if v != nil {
		t.Errorf("Stack not empty")
	}
	s.Push(expect)
	s.Push(0)
	v = s.Pop()
	if v != 0 {
		t.Errorf("Invalid response %d. Expected %d", v.(int), 0)
	}
	v = s.Pop()
	if v.(int) != expect {
		t.Errorf("Invalid response %d. Expected %d", v.(int), expect)
	}
}
