package embertree_test

import (
	"testing"

	"github.com/dufourgilles/emberlib/asn1"
	"github.com/dufourgilles/emberlib/embertree"
)

func TestIsQualifiedTag(t *testing.T) {
	ok := embertree.IsQualifiedTag(embertree.QualifiedNodeApplication)
	if ok == false {
		t.Errorf("IsQualified failed for QualitiedNode")
		return
	}
	ok = embertree.IsQualifiedTag(embertree.NodeApplication)
	if ok == true {
		t.Errorf("IsQualified failed for Node")
		return
	}
}

func TestDecodeQualifiedNode(t *testing.T) {
	buffer := []byte{106, 7, 160, 5, 13, 3, 1, 2, 3}
	reader := asn1.NewASNReader(buffer)
	element, err := embertree.DecodeElement(reader)
	if err != nil {
		t.Error(err)
		return
	}
	if element == nil {
		t.Errorf("QualifiedNode decode failure")
		return
	}
	expectedPath := asn1.RelativeOID{1, 2, 3}
	path := element.GetPath()
	for i, p := range path {
		if p != expectedPath[i] {
			t.Errorf("QualifiedNode decode failure at pos %d. Got %d instead of %d", i, p, expectedPath[i])
			return
		}
	}
}
