package asn1

import (
	"bytes"
	"math"

	"github.com/dufourgilles/emberlib/errors"
)

type ASNReader struct {
	data           *bytes.Reader
	startingOffset int
	len            int
}

func NewASNReader(p []byte) *ASNReader {
	return &ASNReader{data: bytes.NewReader(p), len: len(p), startingOffset: 0}
}

func (a *ASNReader) Len() int {
	return a.data.Len()
}

func (a *ASNReader) Offset() int {
	return a.len - a.data.Len()
}

func (a *ASNReader) TopOffset() int {
	return a.startingOffset + a.Offset()
}

func (a *ASNReader) NewReader(length int) (*ASNReader, errors.Error) {
	if length <= 0 {
		return NewASNReader([]byte{}), nil
	}
	b := make([]byte, length)
	offset := a.TopOffset()
	_, err := a.data.Read(b)
	if err != nil {
		return nil, errors.New("Failed To read byte at offset %d. %s", a.TopOffset(), err)
	}
	newReader := NewASNReader(b)
	newReader.startingOffset = offset
	return newReader, nil
}

func (a *ASNReader) ReadByte() (byte, errors.Error) {
	b, err := a.data.ReadByte()
	if err != nil {
		return b, errors.New("Failed To read byte at offset %d. %s", a.TopOffset(), err)
	}
	return b, nil
}

func (a *ASNReader) ReadBoolean() (bool, errors.Error) {
	offset := a.TopOffset()
	b, err := a.ReadByte()
	if err != nil {
		return false, errors.New("Failed To tag at offset %d. %s", offset, err)
	}
	if b != EMBER_BOOLEAN {
		return false, errors.New("Invalid Boolean type")
	}
	offset = a.TopOffset()
	l, err := a.ReadLength()
	if err != nil {
		return false, errors.New("Failed To length at offset %d. %s", offset, err)
	}
	if l != 1 {
		return false, errors.New("Invalid boolean length %d at offset %d", l, offset)
	}
	b, e := a.ReadByte()
	if err != nil {
		return false, errors.Update(e)
	}
	if b == 0 {
		return false, nil
	}
	return true, nil
}

func (a *ASNReader) ReadOID(tag uint8) (RelativeOID, errors.Error) {
	var oid RelativeOID
	offset := a.TopOffset()
	tag, err := a.ReadByte()
	if err != nil {
		return nil, errors.New("Failed to read OID tag at offset %d. %s", offset, err)
	}
	buf, err := a.readStringBuffer()
	if err != nil {
		return nil, errors.Update(err)
	}

	value := int32(0)
	for i := 0; i < len(buf); i++ {
		b := buf[i] & 0xff
		value <<= 7
		value += int32(b & 0x7f)
		if b&0x80 == 0 {
			oid = append(oid, value)
			value = 0
		}
	}

	return oid, nil
}

func (a *ASNReader) Peek() (byte, errors.Error) {
	b, err := a.data.ReadByte()
	if err != nil {
		return b, errors.New("Failed To read byte at offset %d. %s", a.TopOffset(), err)
	}
	err = a.data.UnreadByte()
	return b, errors.NewError(err)
}

func (a *ASNReader) ReadLength() (int, errors.Error) {
	offset := a.TopOffset()
	lenB, err := a.ReadByte()
	if err != nil {
		return 0, errors.New("Failed to read length at offset %d. %s", offset, err)
	}

	if lenB&0x80 == 0x80 {
		lenB &= 0x7F
		if lenB == 0 {
			return -1, nil
		} else if lenB > 4 {
			return 0, errors.New("Length higher than 4 at offset %d", offset)
		} else {
			len := 0
			for i := 0; i < int(lenB); i++ {
				var val byte
				val, err = a.ReadByte()
				if err != nil {
					return 0, errors.Update(err)
				}
				len = len<<8 + int(val)
			}
			return len, nil
		}
	}
	return int(lenB), nil
}

func (a *ASNReader) ReadSequenceStart(tag uint8) (int, *ASNReader, errors.Error) {
	offset := a.TopOffset()
	b, err := a.ReadByte()
	if err != nil {
		return -1, nil, errors.Update(err)
	}
	if b != tag {
		return -1, a, errors.New("Sequence TAG mismatch at offset %d. Got %d instead of %d", offset, b, tag)
	}
	length, err := a.ReadLength()
	if length >= 0 {
		newReader, err := a.NewReader(length)

		return length, newReader, err
	}
	return length, a, nil
}

func (a *ASNReader) ReadSequenceEnd() errors.Error {
	end, err := a.CheckSequenceEnd()
	if end {
		return nil
	}
	if err != nil {
		return err
	}
	if !end {
		return errors.New("Sequence End not found at offset %d with %d bytes left.", a.TopOffset(), a.Len())
	}
	return nil
}

func (a *ASNReader) CheckSequenceEnd() (bool, errors.Error) {
	if a.Len() == 0 {
		return true, nil
	}
	x1, err := a.data.ReadByte()
	if err != nil {
		return false, errors.New("Failed To read byte at offset %d. %s", a.TopOffset(), err)
	}
	if x1 != 0 {
		a.data.UnreadByte()
		return false, nil
	}
	x1, err = a.data.ReadByte()
	if err != nil {
		return false, errors.New("Failed To read byte at offset %d. %s", a.TopOffset(), err)
	}
	if x1 != 0 {
		a.data.UnreadByte()
		a.data.UnreadByte()
		return false, nil
	}
	return true, nil
}

func (a *ASNReader) ReadInt() (int, errors.Error) {
	offset := a.TopOffset()
	tag, e := a.data.ReadByte()
	if e != nil {
		return 0, errors.New("Failed To read byte at offset %d. %s", offset, e)
	}
	if tag != EMBER_INTEGER {
		return 0, errors.New("Incorrect integer tag %d at offset %d.", tag, offset)
	}

	var l int
	l, err := a.ReadLength()
	if err != nil {
		return 0, err
	}
	if l > 4 {
		return 0, errors.New("Integer length too big %d at offset %d", l, offset)
	}

	var b byte
	val := int(0)
	b, err = a.ReadByte()
	if b&0x80 > 0 {
		val = -1
	}
	for l > 0 {
		l--
		val = (val << 8) | int(b)
		if l > 0 {
			b, err = a.ReadByte()
			if err != nil {
				return 0, errors.Update(err)
			}
		}
	}
	return val, nil
}

func (a *ASNReader) readInt64(l int) (int64, errors.Error) {
	if l > 8 {
		return 0, errors.New("Integer length too big %d at offset %d.", l, a.TopOffset())
	}
	var b byte
	val := int64(0)
	b, err := a.ReadByte()
	if err != nil {
		return 0, errors.Update(err)
	}
	if b&0x80 > 0 {
		val = -1
	}
	for l > 0 {
		l--
		val = (val << 8) | int64(b)
		if l > 0 {
			b, err = a.ReadByte()
			if err != nil {
				return 0, errors.Update(err)
			}
		}
	}
	return val, err
}

func (a *ASNReader) ReadInt64() (int64, errors.Error) {
	offset := a.TopOffset()
	tag, e := a.data.ReadByte()
	if e != nil {
		return 0, errors.New("Failed To read tag at offset %d. %s", offset, e)
	}
	if tag != EMBER_INTEGER {
		return 0, errors.New("Incorrect integer tag at offset %d", offset)
	}

	var l int
	l, err := a.ReadLength()
	if err != nil {
		return 0, errors.Update(err)
	}
	if l > 8 {
		return 0, errors.New("Integer64 length too big %d at offset %d", l, offset)
	}

	val, err := a.readInt64(l)
	return val, nil
}

func (a *ASNReader) ReadIndifiniteLengthData() ([]byte, errors.Error) {
	var b byte
	ret := make([]byte, 0)
	for {
		end, err := a.CheckSequenceEnd()
		if err != nil {
			return ret, errors.Update(err)
		}
		if end {
			break
		}
		b, err = a.ReadByte()
		if err != nil {
			return ret, errors.Update(err)
		}
		ret = append(ret, b)
	}
	return ret, nil
}

func (a *ASNReader) readStringBuffer() ([]byte, errors.Error) {
	var l int
	l, err := a.ReadLength()
	if err != nil {
		return nil, errors.Update(err)
	}

	var b []byte
	if l < 0 {
		b, err = a.ReadIndifiniteLengthData()
		if err != nil {
			return nil, errors.Update(err)
		}
		return b, nil
	} else if l == 0 {
		return nil, nil
	}
	b = make([]byte, l)
	for i := 0; i < l; i++ {
		b[i], err = a.ReadByte()
		if err != nil {
			return nil, errors.Update(err)
		}
	}
	return b, nil
}

func (a *ASNReader) ReadString() (string, errors.Error) {
	offset := a.TopOffset()
	tag, e := a.data.ReadByte()
	if e != nil {
		return "", errors.New("Failed To read tag at offset %d. %s", offset, e)
	}
	if tag != EMBER_STRING {
		return "", errors.New("Incorrect string tag at offset %d.", offset)
	}

	b, err := a.readStringBuffer()
	if err != nil {
		return "", errors.Update(err)
	}
	return string(b), err
}

func (a *ASNReader) ReadBitString() ([]byte, errors.Error) {
	offset := a.TopOffset()
	tag, err := a.data.ReadByte()
	if err != nil {
		return nil, errors.New("Failed To read tag at offset %d. %s", offset, err)
	}
	if tag != EMBER_BITSTRING {
		return nil, errors.New("Incorrect string tag at offset %d.", offset)
	}

	return a.readStringBuffer()
}

func (a *ASNReader) ReadReal() (float64, errors.Error) {
	offset := a.TopOffset()
	tag, err := a.data.ReadByte()
	if err != nil {
		return 0.0, errors.New("Failed To read tag at offset %d. %s", offset, err)
	}
	if tag != EMBER_REAL {
		return 0.0, errors.New("Incorrect real tag")
	}

	buf, e := a.readStringBuffer()
	if e != nil || buf == nil {
		return 0.0, errors.Update(e)
	}

	preamble := buf[0]

	if len(buf) == 1 {
		if preamble == 0x40 {
			return math.Inf(1), nil
		} else if preamble == 0x41 {
			return math.Inf(-1), nil
		} else if preamble == 0x42 {
			return math.NaN(), nil
		} else {
			return math.NaN(), errors.New("Invalid preamble")
		}
	}

	var sign int = -1
	if preamble&0x40 == 0 {
		sign = 1
	}
	exponentLength := 1 + (preamble & 3)
	significandShift := (preamble >> 2) & 3

	exponent := 0
	pos := 1
	if buf[pos]&0x80 != 0 {
		exponent = -1
	}

	if len(buf)-pos < int(exponentLength) {
		return math.NaN(), errors.New("Invalid ASN.1; not enough length to contain exponent")
	}

	for i := 0; i < int(exponentLength); i++ {
		exponent = (exponent << 8) | int(buf[pos])
		pos++
	}

	significand := uint64(0)
	for pos < len(buf) {
		significand = (significand << 8) | uint64(buf[pos])
		pos++
	}

	significand = significand << significandShift

	mask := uint64(0x7FFFF00000000000)
	for significand&mask == 0 {
		significand <<= 8
	}

	mask = uint64(0x7FF0000000000000)
	for significand&mask == 0 {
		significand <<= 1
	}

	significand = significand & uint64(0x000FFFFFFFFFFFFF)
	longExponent := uint64(exponent)
	bits := ((longExponent + 1023) << 52) | significand
	if sign < 0 {
		bits |= uint64(0x8000000000000000)
	}

	return math.Float64frombits(bits), nil

}
