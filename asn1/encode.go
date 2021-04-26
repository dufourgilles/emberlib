package asn1

import (
	"bytes"
	"math"

	"github.com/dufourgilles/emberlib/errors"
)

const EMBER_SEQUENCE uint8 = 0x20 | 16
const EMBER_SET uint8 = 0x20 | 17
const EMBER_BOOLEAN uint8 = 1
const EMBER_INTEGER uint8 = 2
const EMBER_BITSTRING uint8 = 3
const EMBER_OCTETSTRING uint8 = 4
const EMBER_NULL uint8 = 5
const EMBER_OBJECTIDENTIFIER uint8 = 6
const EMBER_OBJECTDESCRIPTOR uint8 = 7
const EMBER_EXTERNAL uint8 = 8
const EMBER_REAL uint8 = 9
const EMBER_ENUMERATED uint8 = 10
const EMBER_EMBEDDED uint8 = 11
const EMBER_STRING uint8 = 12
const EMBER_RELATIVE_OID uint8 = 13

type RelativeOID []int32

type ASNWriter struct {
	data bytes.Buffer
}

func NewASNWriter() *ASNWriter {
	return &ASNWriter{}
}

func (asn *ASNWriter) Read(p []byte) (int, error) {
	return asn.data.Read(p)
}

func (asn *ASNWriter) Len() int {
	return asn.data.Len()
}

func (asn *ASNWriter) WriteByte(b byte) errors.Error {
	err := asn.data.WriteByte(b)
	return errors.NewError(err)
}

// tag should be asn1.
func (asn *ASNWriter) WriteIntTag(i int, tag uint8) errors.Error {
	mask := 0xff800000
	intsize := 4
	for {
		if !((i&mask == 0 || (i&mask) == mask) && intsize > 1) {
			break
		}
		intsize--
		i <<= 8
	}
	err := asn.data.WriteByte(tag)
	if err != nil {
		return errors.NewError(err)
	}
	err = asn.data.WriteByte(byte(intsize))
	if err != nil {
		return errors.NewError(err)
	}
	mask = 0xff000000
	for {
		if intsize <= 0 {
			break
		}
		intsize--
		err = asn.data.WriteByte(byte((i & mask) >> 24))
		if err != nil {
			return errors.NewError(err)
		}
		i <<= 8
	}
	return nil
}

func (asn *ASNWriter) WriteInt(i int) errors.Error {
	return asn.WriteIntTag(i, EMBER_INTEGER)
}

func (asn *ASNWriter) WriteInt64Tag(i int64, tag uint8) errors.Error {
	intsize := 1
	temp := i
	for temp > 127 {
		intsize++
		temp >>= 8
	}

	for temp < -128 {
		intsize++
		temp >>= 8
	}

	err := asn.data.WriteByte(tag)
	if err != nil {
		return errors.NewError(err)
	}
	err = asn.data.WriteByte(byte(intsize))
	if err != nil {
		return errors.NewError(err)
	}
	for ; intsize > 0; intsize-- {
		err = asn.data.WriteByte(byte(i >> uint((intsize-1)*8)))
		if err != nil {
			return errors.NewError(err)
		}
	}
	return nil
}

func (asn *ASNWriter) WriteInt64(i int64) errors.Error {
	return asn.WriteInt64Tag(i, EMBER_INTEGER)
}

func (asn *ASNWriter) WriteNull() errors.Error {
	err := asn.data.WriteByte(EMBER_NULL)
	return errors.NewError(err)
}

func (asn *ASNWriter) WriteEnum(i int) errors.Error {
	return asn.WriteIntTag(i, EMBER_ENUMERATED)
}

func (asn *ASNWriter) WriteBoolean(b bool) errors.Error {
	err := asn.data.WriteByte(EMBER_BOOLEAN)
	if err != nil {
		return errors.NewError(err)
	}
	err = asn.data.WriteByte(1)
	if err != nil {
		return errors.NewError(err)
	}
	if b {
		err = asn.data.WriteByte(0xFF)
		return errors.NewError(err)
	}
	err = asn.data.WriteByte(0x00)
	return errors.NewError(err)
}

func shorten(value int) (int, int) {
	size := 4
	for ((value&0xff800000) == 0 || (value&0xff800000) == 0xff800000) && (size > 1) {
		size--
		value <<= 8
	}
	return size, value
}

func shortenLong(value uint64) (int, uint64) {
	mask := uint64(0xff80000000000000)
	size := 8

	for (value&mask == 0 || value&mask == mask) && size > 1 {
		size--
		value <<= 8
	}
	return size, value
}

func (asn *ASNWriter) WriteReal(r float64) errors.Error {
	err := asn.WriteByte(EMBER_REAL)
	if err != nil {
		return err
	}
	if r == 0.0 {
		return asn.writeLength(0)
	} else if math.IsInf(r, 1) {
		e := asn.writeLength(1)
		if e != nil {
			return errors.Update(e)
		}
		return asn.WriteByte(0x40)
	} else if math.IsInf(r, -1) {
		e := asn.writeLength(1)
		if e != nil {
			return errors.Update(e)
		}
		return asn.WriteByte(0x41)
	} else if math.IsNaN(r) {
		e := asn.writeLength(1)
		if e != nil {
			return errors.Update(e)
		}
		return asn.WriteByte(0x42)
	}
	// value := []byte(strconv.FormatFloat(r, 'G', -1, 64))
	// var buf []byte
	// if bytes.Contains(value, []byte{'E'}) {
	// 	buf = []byte{0x03}
	// } else {
	// 	buf = []byte{0x02}
	// }
	// buf = append(buf, value...)
	// return asn.writeUntagBuffer(buf)

	bits := math.Float64bits(r)
	significand := (uint64(0x000FFFFFFFFFFFFF) & bits) | 0x0010000000000000
	exponent := int(((uint64(0x7FF0000000000000) & bits) >> 52)) - 1023

	for significand&uint64(0xFF) == 0 {
		significand = significand >> 8
	}

	for significand&uint64(0x01) == 0 {
		significand = significand >> 1
	}

	expSize, expVal := shorten(exponent)

	significandSize, significandVal := shortenLong(significand)

	asn.writeLength(1 + expSize + significandSize)
	preamble := uint8(0x80)
	if r < 0 {
		preamble |= 0x40
	}
	asn.WriteByte(preamble)

	for i := 0; i < expSize; i++ {
		asn.WriteByte(uint8((expVal & 0xFF000000) >> 24))
		expVal <<= 8
	}

	const mask = uint64(0xFF00000000000000)
	for i := 0; i < significandSize; i++ {
		b := significandVal & mask
		err = asn.WriteByte(uint8(b >> 56))
		significandVal <<= 8
	}
	return err
}

func (asn *ASNWriter) writeLength(len int) errors.Error {
	var err error
	if len <= 0x7f {
		err = asn.data.WriteByte(byte(len & 0xFF))
	} else if len <= 0xff {
		err = asn.data.WriteByte(0x81)
		if err != nil {
			return errors.NewError(err)
		}
		err = asn.data.WriteByte(byte(len & 0xFF))
	} else if len <= 0xffff {
		err = asn.data.WriteByte(0x82)
		if err != nil {
			return errors.NewError(err)
		}
		err = asn.data.WriteByte(byte((len >> 8) & 0xFF))
		if err != nil {
			return errors.NewError(err)
		}
		err = asn.data.WriteByte(byte(len & 0xFF))
	} else if len <= 0xffffff {
		err = asn.data.WriteByte(0x82)
		if err != nil {
			return errors.NewError(err)
		}
		err = asn.data.WriteByte(byte((len >> 16) & 0xFF))
		if err != nil {
			return errors.NewError(err)
		}
		err = asn.data.WriteByte(byte((len >> 8) & 0xFF))
		if err != nil {
			return errors.NewError(err)
		}
		err = asn.data.WriteByte(byte(len & 0xFF))
	}
	return errors.NewError(err)
}

func (asn *ASNWriter) WriteString(s string) errors.Error {
	err := asn.data.WriteByte(EMBER_STRING)
	if err != nil {
		return errors.NewError(err)
	}
	l := len(s)
	e := asn.writeLength(l)
	if e != nil {
		return errors.Update(e)
	}
	if l > 0 {
		var res int
		res, err = asn.data.WriteString(s)
		if err != nil {
			return errors.NewError(err)
		}
		if l != res {
			return errors.New("Failed to write full string.")
		}
	}
	return nil
}

func (asn *ASNWriter) writeUntagBuffer(b []byte) errors.Error {
	l := len(b)
	err := asn.writeLength(l)
	if err != nil {
		return err
	}
	if l > 0 {
		var res int
		res, e := asn.data.Write(b)
		if e != nil {
			return errors.NewError(e)
		}
		if l != res {
			return errors.New("Failed to write full buffer.")
		}
	}
	return err
}

func (asn *ASNWriter) WriteBuffer(b []byte, tag uint8) errors.Error {
	err := asn.data.WriteByte(tag)
	if err != nil {
		return errors.NewError(err)
	}
	return asn.writeUntagBuffer(b)
}

func (asn *ASNWriter) WriteRelativeOID(oid RelativeOID) errors.Error {
	buffer := make([]byte, 0)
	for i := 0; i < len(oid); i++ {
		val := oid[i]
		encodedBytes := make([]byte, 0)
		for {
			if val < 0x7F {
				break
			}
			encodedBytes = append(encodedBytes, byte(val&0x7F))
			val = val >> 7
		}
		encodedBytes = append(encodedBytes, byte(val))
		for j := len(encodedBytes) - 1; j > 0; j-- {
			buffer = append(buffer, encodedBytes[j]|0x80)
		}
		buffer = append(buffer, encodedBytes[0])
	}
	return asn.WriteBuffer(buffer, EMBER_RELATIVE_OID)
}

func (asn *ASNWriter) StartSequence(tag uint8) errors.Error {
	err := asn.data.WriteByte(tag)
	if err != nil {
		return errors.NewError(err)
	}
	err = asn.data.WriteByte(0x80)
	return errors.NewError(err)
}

func (asn *ASNWriter) EndSequence() errors.Error {
	l, err := asn.data.Write([]byte{0, 0})
	if l != 2 {
		return errors.New("Failed to terminate sequence")
	}
	if err != nil {
		return errors.NewError(err)
	}
	return nil
}

func Application(num uint8) uint8 {
	return 0x60 | num
}

func Context(num uint8) uint8 {
	return 0xa0 | num
}

func Universal(num uint8) uint8 {
	return num
}
