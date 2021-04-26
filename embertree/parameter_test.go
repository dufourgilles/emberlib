package embertree_test

import (
	"testing"

	"github.com/dufourgilles/emberlib/asn1"
	"github.com/dufourgilles/emberlib/embertree"
)

func TestEncodeParameter(t *testing.T) {
	paramID := int(10)
	parameter := embertree.NewParameter(paramID)
	writer := asn1.ASNWriter{}
	err := parameter.Encode(&writer)
	if err != nil {
		t.Error(err)
	}
	b := make([]byte, writer.Len())
	writer.Read(b)
	expectedResult := []byte{97, 128, 160, 128, 2, 1, 10, 0, 0, 0, 0}
	for i, d := range b {
		if d != expectedResult[i] {
			t.Errorf("Invalid byte at %d.  Val %d expected %d.", i, d, expectedResult[i])
		}
	}
}

func TestEncodeParameterIdentifier(t *testing.T) {
	paramID := int(10)
	parameter := embertree.NewParameter(paramID)
	parameterContent := parameter.CreateContent().(*embertree.ParameterContents)
	parameterContent.SetIdentifier("gdnet")
	writer := asn1.ASNWriter{}
	err := parameter.Encode(&writer)
	if err != nil {
		t.Error(err)
	}
	b := make([]byte, writer.Len())
	writer.Read(b)
	//fmt.Println(b)
	expectedResult := []byte{97, 128, 160, 128, 2, 1, 10, 0, 0, 161, 128, 49, 128, 160, 128, 12, 5, 103, 100, 110, 101, 116, 0, 0, 0, 0, 0, 0, 0, 0}
	for i, d := range b {
		if d != expectedResult[i] {
			t.Errorf("Invalid byte at %d.  Val %d expected %d.", i, d, expectedResult[i])
		}
	}
}

func TestEncodeParameterValue(t *testing.T) {
	paramID := int(10)
	parameter := embertree.NewParameter(paramID)
	parameterContent := parameter.CreateContent().(*embertree.ParameterContents)
	parameterContent.SetIdentifier("gdnet")
	val := parameterContent.GetValueObject()
	val.SetInt(1234)
	writer := asn1.ASNWriter{}
	err := parameter.Encode(&writer)
	if err != nil {
		t.Error(err)
	}
	b := make([]byte, writer.Len())
	writer.Read(b)
	//fmt.Println(b)
	expectedResult := []byte{97, 128, 160, 128, 2, 1, 10, 0, 0, 161, 128, 49, 128, 160, 128, 12, 5, 103, 100, 110, 101, 116, 0, 0, 162, 128, 2, 2, 4, 210, 0, 0, 0, 0, 0, 0, 0, 0}
	for i, d := range b {
		if d != expectedResult[i] {
			t.Errorf("Invalid byte at %d.  Val %d expected %d.", i, d, expectedResult[i])
		}
	}
}

func TestDecodeParameterValue(t *testing.T) {
	encodedParameter := []byte{97, 128, 160, 128, 2, 1, 10, 0, 0, 161, 128, 49, 128, 160, 128, 12, 5, 103, 100, 110, 101, 116, 0, 0, 162, 128, 2, 2, 4, 210, 0, 0, 0, 0, 0, 0, 0, 0}
	reader := asn1.NewASNReader(encodedParameter)
	parameter, err := embertree.DecodeElement(reader)
	if err != nil {
		t.Error(err)
		return
	}
	parameterContent := parameter.GetContent().(*embertree.ParameterContents)
	//fmt.Println("done")
	identifier, err := parameterContent.GetIdentifier()
	if err != nil {
		t.Error(err)
		return
	}
	if identifier != "gdnet" {
		t.Errorf("Invalid identifier. Expected gdnet. Received %s", identifier)
		return
	}
	val := parameterContent.GetValueObject()
	valInt, err := val.GetInt()
	if err != nil {
		t.Error(err)
		return
	}
	if valInt != 1234 {
		t.Errorf("Invalid identifier. Expected 1234. Received %d", valInt)
		return
	}
}
