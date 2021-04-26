package embertree_test

import (
	//"fmt"
	"testing"

	"github.com/dufourgilles/emberlib/asn1"
	. "github.com/dufourgilles/emberlib/embertree"
)

func TestLabeLEncode(t *testing.T) {
	oid := asn1.RelativeOID{1, 2, 3}
	l := NewLabel(oid, "label")
	writer := asn1.NewASNWriter()
	err := l.Encode(writer)
	if err != nil {
		t.Errorf("label encode failure")
		t.Error(err)
		return
	}
}

func TestLabeLDecode(t *testing.T) {
	buf := []byte{114, 16, 160, 5, 13, 3, 1, 2, 3, 161, 7, 12, 5, 108, 97, 98, 101, 108}
	l := Label{}
	reader := asn1.NewASNReader(buf)
	err := l.Decode(reader)
	if err != nil {
		t.Errorf("label encode failure")
		t.Error(err)
		return
	}
	oid := asn1.RelativeOID{1, 2, 3}
	for i := range l.BasePath {
		if oid[i] != l.BasePath[i] {
			t.Errorf("Invalid OID at %d.  Got %d instead of %d", i, oid[i], l.BasePath[i])
			return
		}
	}
	if l.Description != "label" {
		t.Errorf("Invalid description %s instead of label", l.Description)
		return
	}
}
