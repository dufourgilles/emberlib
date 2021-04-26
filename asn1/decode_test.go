package asn1_test

import (
	"fmt"
	"testing"

	. "github.com/dufourgilles/emberlib/asn1"
)

func TestASN1ReadOID(t *testing.T) {
	asn := NewASNReader([]byte{0x0d, 0x04, 0x82, 0x1, 0x8, 0x09})
	oid, err := asn.ReadOID(13)
	if err != nil {
		t.Error(err)
		return
	}
	expectedOID := RelativeOID{257, 8, 9}
	//fmt.Println(oid)
	for i, o := range expectedOID {
		if o != oid[i] {
			t.Errorf("Invalid OID value at index %d. Received %d not %d", i, oid[i], o)
			return
		}
	}
}

func TestASN1DecodeString(t *testing.T) {
	reader := NewASNReader([]byte{0x0c, 05, 0x67, 0x64, 0x6e, 0x65, 0x74})
	str, err := reader.ReadString()
	if err != nil {
		t.Error(err)
		return
	}
	if str != "gdnet" {
		t.Errorf("Invalid string")
		return
	}
	if reader.Len() > 0 {
		t.Errorf("Buffer not empty")
		return
	}
}

func TestASN1DecodeOffset(t *testing.T) {
	buf := []byte{96, 16, 107, 14, 160, 12, 98, 10, 160, 3, 2, 1, 32, 161, 3, 2, 1, 0xFF}
	topreader := NewASNReader(buf)
	fmt.Printf("Length: %d\n", topreader.Len())
	if topreader.Len() != 18 {
		t.Errorf("Reader length mismatch %d/%d", topreader.Len(), len(buf))
		return
	}
	topreader.ReadByte()
	offset := topreader.TopOffset()
	if offset != 1 {
		t.Errorf("Offset mismatch. Currently %d out of %d", offset, topreader.Len())
		return
	}
	topreader.ReadLength()
	offset = topreader.TopOffset()
	if offset != 2 {
		t.Errorf("Offset mismatch. Currently %d out of %d", offset, topreader.Len())
		return
	}
	topreader = NewASNReader(buf)
	_, reader, err := topreader.ReadSequenceStart(Application(0))
	if err != nil {
		t.Error(err)
		return
	}
	if reader.Len() != 16 {
		t.Errorf("reader length mismatch %d", reader.Len())
		return
	}
	offset = topreader.TopOffset()
	if offset != 18 {
		t.Errorf("Offset mismatch. Currently %d out of %d", offset, topreader.Len())
		return
	}
	offset = reader.Offset()
	if offset != 0 {
		t.Errorf("Offset mismatch")
		return
	}
	offset = reader.TopOffset()
	if offset != 2 {
		t.Errorf("Offset mismatch")
		return
	}

	_, subreader, err := reader.ReadSequenceStart(Application(107))
	if err != nil {
		t.Error(err)
		return
	}
	if subreader.Len() != 14 {
		t.Errorf("reader length mismatch %d", subreader.Len())
		return
	}
	offset = reader.Offset()
	if offset != 16 {
		t.Errorf("Offset mismatch. Currently %d out of %d", offset, reader.Len())
		return
	}
	offset = subreader.Offset()
	if offset != 0 {
		t.Errorf("Offset mismatch. Currently %d out of %d", offset, reader.Len())
		return
	}
	offset = subreader.TopOffset()
	if offset != 4 {
		t.Errorf("Offset mismatch. Currently %d out of %d", offset, reader.Len())
		return
	}
}
