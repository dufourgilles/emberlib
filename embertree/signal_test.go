package embertree_test

import (
	//"fmt"
	"testing"

	"github.com/dufourgilles/emberlib/asn1"
	. "github.com/dufourgilles/emberlib/embertree"
)

func TestTargetlEncodeDecode(t *testing.T) {
	number := int32(17)
	target := NewTarget(number)
	writer := asn1.NewASNWriter()
	err := target.Encode(writer)
	if err != nil {
		t.Errorf("target encode failure")
		t.Error(err)
		return
	}

	buf := make([]byte, writer.Len())
	writer.Read(buf)
	reader := asn1.NewASNReader(buf)

	decodedTarget := NewTarget(0)
	err = decodedTarget.Decode(reader)
	if err != nil {
		t.Errorf("target decode failure")
		t.Error(err)
		return
	}
	if decodedTarget.Number != number {
		t.Errorf("target decode failure. Invalid number %d instead of %d", decodedTarget.Number, number)
	}
}
