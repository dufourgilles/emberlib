package embertree_test

import (
	//"fmt"
	"testing"

	"github.com/dufourgilles/emberlib/asn1"
	"github.com/dufourgilles/emberlib/embertree"
)

func TestContentParameterBool(t *testing.T) {
	c := embertree.ContentParameter{}
	val := true
	c.SetBool(val)

	s, err := c.GetBool()
	if err != nil {
		t.Error(err)
		return
	}
	if s != val {
		t.Errorf("Incorrect bool received.")
		return
	}

	_, err = c.GetString()
	if err == nil {
		t.Errorf("Unxpected string value received instead of bool")
		return
	}

	_, err = c.GetBuffer()
	if err == nil {
		t.Errorf("Unxpected buffer value received instead of bool")
		return
	}

	_, err = c.GetInt()
	if err == nil {
		t.Errorf("Unxpected int value received instead of bool")
		return
	}
}

func TestContentParameterString(t *testing.T) {
	c := embertree.ContentParameter{}
	str := "gdnet"
	c.SetString(str)

	s, err := c.GetString()
	if err != nil {
		t.Error(err)
		return
	}
	if s != str {
		t.Errorf("Incorrect string received. %s - %s", s, str)
		return
	}

	_, err = c.GetBool()
	if err == nil {
		t.Errorf("Unxpected bool value received instead of string")
		return
	}

	_, err = c.GetBuffer()
	if err == nil {
		t.Errorf("Unxpected buffer value received instead of string")
		return
	}

	_, err = c.GetInt()
	if err == nil {
		t.Errorf("Unxpected int value received instead of string")
		return
	}
}

func TestContentParameterInt(t *testing.T) {
	c := embertree.ContentParameter{}
	val := int64(77)
	c.SetInt(val)

	s, err := c.GetInt()
	if err != nil {
		t.Error(err)
		return
	}
	if s != val {
		t.Errorf("Incorrect int received. %d - %d", s, val)
		return
	}

	_, err = c.GetString()
	if err == nil {
		t.Errorf("Unxpected string value received instead of int")
		return
	}

	_, err = c.GetBuffer()
	if err == nil {
		t.Errorf("Unxpected buffer value received instead of int")
		return
	}

	_, err = c.GetBool()
	if err == nil {
		t.Errorf("Unxpected bool value received instead of int")
		return
	}
}

func TestContentParameterBuffer(t *testing.T) {
	c := embertree.ContentParameter{}
	val := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	c.SetBuffer(val)

	s, err := c.GetBuffer()
	if err != nil {
		t.Error(err)
		return
	}
	if s == nil || len(s) != len(val) {
		t.Errorf("Incorrect buffer received. %d - %d", len(s), len(val))
		return
	}
	for i, _ := range val {
		if s[i] != val[i] {
			t.Errorf("Incorrect buffer received. %d - %d", s[i], val[i])
			return
		}
	}

	_, err = c.GetString()
	if err == nil {
		t.Errorf("Unxpected string value received instead of buffer")
		return
	}

	_, err = c.GetInt()
	if err == nil {
		t.Errorf("Unxpected int value received instead of buffer")
		return
	}

	_, err = c.GetBool()
	if err == nil {
		t.Errorf("Unxpected bool value received instead of buffer")
		return
	}
}

func TestContentParameterEncodeDecode(t *testing.T) {
	c := embertree.ContentParameter{}
	val := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	c.SetBuffer(val)

	writer := asn1.NewASNWriter()
	err := c.Encode(1, writer)
	if err != nil {
		t.Error(err)
		return
	}

	p := make([]byte, writer.Len())
	writer.Read(p)

	// for _,b := range(p) {
	// 	fmt.Printf("0x%x, ", b)
	// }
	reader := asn1.NewASNReader(p)

	dup, err := embertree.DecodeValue(reader, 1)
	if err != nil {
		t.Error(err)
		return
	}
	if dup == nil {
		t.Errorf("Failed to decode")
	}

	s, err := dup.GetBuffer()
	if err != nil {
		t.Error(err)
		return
	}
	if s == nil || len(s) != len(val) {
		t.Errorf("Incorrect buffer received. %d - %d", len(s), len(val))
		return
	}
	for i, _ := range val {
		if s[i] != val[i] {
			t.Errorf("Incorrect buffer received. %d - %d", s[i], val[i])
			return
		}
	}
}
