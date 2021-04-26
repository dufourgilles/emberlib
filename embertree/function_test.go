package embertree_test

import (
	//"fmt"
	"testing"

	"github.com/dufourgilles/emberlib/asn1"
	"github.com/dufourgilles/emberlib/embertree"
)

func TestFunctionDecode(t *testing.T) {
	encodedFunction := []byte{115, 38, 160, 3, 2, 1, 4, 161, 31, 49, 29, 160, 4, 12, 2, 70, 49, 162, 21, 48, 19, 160, 17, 117, 15, 160, 3, 2, 1, 1, 161, 8, 12, 6, 112, 97, 114, 97, 109, 49}
	reader := asn1.NewASNReader(encodedFunction)
	f, err := embertree.DecodeElement(reader)
	if err != nil {
		t.Error(err)
		return
	}
	contents := f.GetContent()
	arguments := contents.(*embertree.FunctionContents).GetArguments()
	if arguments == nil || len(arguments) != 1 {
		t.Errorf("Function decode error. Invalid arguments count. Got %d insteaf of 1", len(arguments))
	}

}
