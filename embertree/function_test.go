package embertree_test

import (
	//"fmt"
	"fmt"
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

func TestFunctionEncode(t *testing.T) {
	writer := asn1.NewASNWriter()
	f := embertree.NewFunction(7)
	contents := f.CreateContent();
	if contents == nil {
		t.Errorf("Failed to get function content")
		return
	}
	fc := contents.(*embertree.FunctionContents)
	fc.SetIdentifier("F1")
	argument := embertree.NewArgument(embertree.ParameterTypeInteger, "param1")
	arguments := []*embertree.TupleDescription{argument}
	fc.SetArguments(arguments)
	err := f.Encode(writer)
	if err != nil {
		t.Error(err)
		return
	}
	b := make([]byte, writer.Len())
	_,err = writer.Read(b)
	if err != nil {
		t.Error(err)
	}

	s := ""
	for _,bb := range(b) {
		s = fmt.Sprintf("%s,0x%x", s, bb)
	}
	fmt.Println(s)

	reader := asn1.NewASNReader(b)
	decodedFunction, err := embertree.DecodeElement(reader) 
	if err != nil {
		t.Error(err)
	}

	fmt.Println(decodedFunction)
	//t.Errorf("err")
}