package asn1_test

import (
	"testing"

	. "github.com/dufourgilles/emberlib/asn1"
)

func TestASN1WriteByte(t *testing.T) {
	asn := ASNWriter{}
	b := byte(34)
	asn.WriteByte(b)
	if asn.Len() != 1 {
		t.Errorf("Invalid object length after WriteByte. Len %d not 1", asn.Len())
		return
	}
	res := make([]byte, 1)
	asn.Read(res)
	if res[0] != b {
		t.Errorf("Invalid object value after WriteByte. Received %d not %d", int(res[0]), int(b))
		return
	}
}

func TestASN1WriteInt(t *testing.T) {
	asn := ASNWriter{}
	b := 77
	asn.WriteInt(b)
	if asn.Len() != 3 {
		t.Errorf("Invalid object length after WriteInt. Len %d not 3", asn.Len())
		return
	}
	expectedResult := []byte{0x2, 0x1, 0x4d}
	res := make([]byte, 3)
	asn.Read(res)
	for i := 0; i < 3; i++ {
		if res[i] != expectedResult[i] {
			t.Errorf("Invalid object value after WriteByte. Received %d not %d", int(res[i]), int(expectedResult[i]))
			return
		}
	}

	b = 0x1FF << 23
	asn.WriteInt(b)
	expectedResult = []byte{0x2, 0x3, 0x80, 0x0, 0x0}
	res = make([]byte, len(expectedResult))
	asn.Read(res)
	//fmt.Println(asn)
	for i := 0; i < len(expectedResult); i++ {
		if res[i] != expectedResult[i] {
			t.Errorf("Invalid object value after WriteByte. Received %d not %d at %d", int(res[i]), int(expectedResult[i]), i)
			return
		}
	}
}

func TestASN1WriteString(t *testing.T) {
	asn := ASNWriter{}
	s := "gdnet"

	asn.WriteString(s)
	expectedResult := []byte{0xC, 0x05, 0x67, 0x64, 0x6e, 0x65, 0x74}
	res := make([]byte, len(expectedResult))
	asn.Read(res)
	//fmt.Println(asn)
	for i := 0; i < len(expectedResult); i++ {
		if res[i] != expectedResult[i] {
			t.Errorf("Invalid object value after WriteByte. Received %d not %d at %d", int(res[i]), int(expectedResult[i]), i)
			return
		}
	}
}

func TestASN1WriteRelativeOID(t *testing.T) {
	asn := ASNWriter{}
	oid := []int32{345, 11, 12}

	asn.WriteRelativeOID(oid)
	expectedResult := []byte{0x0d, 0x04, 0x82, 0x59, 0x0b, 0x0c}
	res := make([]byte, len(expectedResult))
	asn.Read(res)
	//fmt.Println(asn)
	for i := 0; i < len(expectedResult); i++ {
		if res[i] != expectedResult[i] {
			t.Errorf("Invalid object value after WriteByte. Received %d not %d at %d", int(res[i]), int(expectedResult[i]), i)
			return
		}
	}
}

func TestASN1WriteReal(t *testing.T) {
	asn := ASNWriter{}
	real := float64(123.456)

	err := asn.WriteReal(real)
	if err != nil {
		t.Error(err)
		return
	}
	expectedResult := []byte{0x09, 0x09, 0x80, 0x06, 0x1e, 0xdd, 0x2f, 0x1a, 0x9f, 0xbe, 0x77}
	res := make([]byte, asn.Len())
	asn.Read(res)
	// for _,b := range(res) {
	// 	fmt.Printf("%x ", b)
	// }
	for i := 0; i < len(expectedResult); i++ {
		if res[i] != expectedResult[i] {
			t.Errorf("Invalid object value after WriteByte. Received %d not %d at %d", int(res[i]), int(expectedResult[i]), i)
			return
		}
	}
}

func TestASN1ReadReal(t *testing.T) {
	asn := NewASNReader([]byte{0x09, 0x09, 0x80, 0x06, 0x1e, 0xdd, 0x2f, 0x1a, 0x9f, 0xbe, 0x77})
	expectedReal := float64(123.456)
	real, err := asn.ReadReal()
	if err != nil {
		t.Error(err)
		return
	}
	if real != expectedReal {
		t.Errorf("Real Decode error %f - %f", real, expectedReal)
		return
	}
}
