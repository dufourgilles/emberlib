package embertree

import (
	"fmt"

	"github.com/dufourgilles/emberlib/asn1"
	"github.com/dufourgilles/emberlib/errors"
)

type ValueType uint8

const (
	ValueTypeUnset   ValueType = iota
	ValueTypeString  ValueType = iota
	ValueTypeBool    ValueType = iota
	ValueTypeInteger ValueType = iota
	ValueTypeBuffer  ValueType = iota
	ValueTypeReal    ValueType = iota
	ValueTypeOID     ValueType = iota
)

func ValueType2String(v ValueType) string {
	switch v {
	case ValueTypeUnset:
		return "unset"
	case ValueTypeBool:
		return "bool"
	case ValueTypeBuffer:
		return "buffer"
	case ValueTypeInteger:
		return "integer"
	case ValueTypeString:
		return "string"
	case ValueTypeReal:
		return "real"
	case ValueTypeOID:
		return "oid"
	default:
		return "unknown"
	}
}

type ContentParameter struct {
	stringVal string
	intVal    int64
	boolVal   bool
	bufferVal []byte
	realVal   float64
	oid       asn1.RelativeOID
	isSet     bool
	valueType ValueType
}

func NewContentParameter() *ContentParameter {
	return &ContentParameter{isSet: false}
}

func (cp *ContentParameter)ToString() string {
	t := cp.GetType()
	switch t {
	case ValueTypeBool:
		val,_ := cp.GetBool()
		var str string
		if val {
			str = "true"
		} else {
			str = "false"
		}
		return fmt.Sprintf("%s", str)
	case ValueTypeString:
		str,_ := cp.GetString()
		return str
	case ValueTypeInteger:
		i,_ := cp.GetInt()
		return fmt.Sprintf("%d", i)
	case ValueTypeOID:
		str := ""
		oid,_ := cp.GetRelativeOID()
		for index,val := range(oid) {
			if index == 0 {
				str = fmt.Sprintf("%d",val)
			} else {
				str = fmt.Sprintf("%s.%d", str, val)
			}
		}
		return str
	case ValueTypeReal:
		r,_ := cp.GetReal()
		return fmt.Sprintf("%f", r)
	case ValueTypeBuffer:
		b,_ := cp.GetBuffer()
		str := ""
		for index,x := range(b) {
			if index == 0 {
				str = fmt.Sprintf("%x", x)
			} else {
				str = fmt.Sprintf("%s%x", str, x)
			}
		}
		return str
	}
	return ""
}

func (cp *ContentParameter) GetType() ValueType {
	return cp.valueType
}

func (cp *ContentParameter) Set(val *ContentParameter) errors.Error {
	t := cp.GetType()
	if t != ValueTypeUnset && t != val.GetType() {
		return errors.New("Type mismatch. %s - %s", ValueType2String(cp.GetType()), ValueType2String(val.GetType()))
	}
	switch val.GetType() {
	case ValueTypeUnset:
		break
	case ValueTypeBool:
		b, err := val.GetBool()
		if err != nil {
			return err
		}
		cp.SetBool(b)
		break
	case ValueTypeBuffer:
		b, err := val.GetBuffer()
		if err != nil {
			return err
		}
		cp.SetBuffer(b)
		break
	case ValueTypeString:
		b, err := val.GetString()
		if err != nil {
			return err
		}
		cp.SetString(b)
		break
	case ValueTypeInteger:
		b, err := val.GetInt()
		if err != nil {
			return err
		}
		cp.SetInt(b)
		break
	case ValueTypeReal:
		b, err := val.GetReal()
		if err != nil {
			return err
		}
		cp.SetReal(b)
		break
	}
	return nil
}

func (cp *ContentParameter) SetString(s string) {
	cp.stringVal = s
	cp.isSet = true
	cp.valueType = ValueTypeString
}

func (cp *ContentParameter) SetBool(b bool) {
	cp.boolVal = b
	cp.isSet = true
	cp.valueType = ValueTypeBool
}

func (cp *ContentParameter) SetInt(i int64) {
	cp.intVal = i
	cp.isSet = true
	cp.valueType = ValueTypeInteger
}

func (cp *ContentParameter) SetBuffer(b []byte) {
	cp.bufferVal = b
	cp.isSet = true
	cp.valueType = ValueTypeBuffer
}

func (cp *ContentParameter) SetReal(val float64) {
	cp.realVal = val
	cp.isSet = true
	cp.valueType = ValueTypeReal
}

func (cp *ContentParameter) SetRelativeOID(val asn1.RelativeOID) {
	cp.oid = val
	cp.isSet = true
	cp.valueType = ValueTypeOID
}

func (cp *ContentParameter) IsSet() bool {
	return cp.isSet
}

func (cp *ContentParameter) GetString() (string, errors.Error) {
	var err errors.Error = nil
	if !cp.isSet {
		err = errors.New("Parameter not set")
	} else if cp.valueType != ValueTypeString {
		err = errors.New("Type mismatch. Requested string but is %s", ValueType2String(cp.valueType))
	}
	return cp.stringVal, err
}

func (cp *ContentParameter) GetBool() (bool, errors.Error) {
	var err errors.Error = nil
	if !cp.isSet {
		err = errors.New("Parameter not set")
	} else if cp.valueType != ValueTypeBool {
		err = errors.New("Type mismatch. Requested bool but is %s", ValueType2String(cp.valueType))
	}
	return cp.boolVal, err
}

func (cp *ContentParameter) GetInt() (int64, errors.Error) {
	var err errors.Error = nil
	if !cp.isSet {
		err = errors.New("Parameter not set")
	} else if cp.valueType != ValueTypeInteger {
		err = errors.New("Type mismatch. Requested int but is %s", ValueType2String(cp.valueType))
	}
	return cp.intVal, err
}

func (cp *ContentParameter) GetBuffer() ([]byte, errors.Error) {
	var err errors.Error = nil
	if !cp.isSet {
		err = errors.New("Parameter not set")
	} else if cp.valueType != ValueTypeBuffer {
		err = errors.New("Type mismatch. Requested buffer but is %s", ValueType2String(cp.valueType))
	}
	return cp.bufferVal, err
}

func (cp *ContentParameter) GetReal() (float64, errors.Error) {
	var err errors.Error = nil
	if !cp.isSet {
		err = errors.New("Parameter not set")
	} else if cp.valueType != ValueTypeReal {
		err = errors.New("Type mismatch. Requested buffer but is %s", ValueType2String(cp.valueType))
	}
	return cp.realVal, err
}

func (cp *ContentParameter) GetRelativeOID() (asn1.RelativeOID, errors.Error) {
	var err errors.Error = nil
	if !cp.isSet {
		err = errors.New("Parameter not set")
	} else if cp.valueType != ValueTypeOID {
		err = errors.New("Type mismatch. Requested buffer but is %s", ValueType2String(cp.valueType))
	}
	return cp.oid, err
}

func (cp *ContentParameter) Encode(context uint8, writer *asn1.ASNWriter) errors.Error {
	if !cp.isSet {
		return nil
	}
	err := writer.StartSequence(asn1.Context(context))
	if err != nil {
		return errors.Update(err)
	}

	switch cp.valueType {
	case ValueTypeBool:
		b, err := cp.GetBool()
		if err != nil {
			return errors.Update(err)
		}
		err = writer.WriteBoolean(b)
		if err != nil {
			return errors.Update(err)
		}
		break
	case ValueTypeBuffer:
		b, err := cp.GetBuffer()
		if err != nil {
			return errors.Update(err)
		}
		err = writer.WriteBuffer(b, asn1.EMBER_BITSTRING)
		if err != nil {
			return errors.Update(err)
		}
		break
	case ValueTypeInteger:
		b, err := cp.GetInt()
		if err != nil {
			return errors.Update(err)
		}
		err = writer.WriteInt64(b)
		if err != nil {
			return errors.Update(err)
		}
		break
	case ValueTypeString:
		b, err := cp.GetString()
		if err != nil {
			return errors.Update(err)
		}
		err = writer.WriteString(b)
		if err != nil {
			return errors.Update(err)
		}
		break
	case ValueTypeReal:
		b, err := cp.GetReal()
		if err != nil {
			return errors.Update(err)
		}
		err = writer.WriteReal(b)
		if err != nil {
			return errors.Update(err)
		}
		break
	case ValueTypeOID:
		b, err := cp.GetRelativeOID()
		if err != nil {
			return errors.Update(err)
		}
		err = writer.WriteRelativeOID(b)
		if err != nil {
			return errors.Update(err)
		}
		break
	default:
		return errors.New("Unknown value type")
	}
	return writer.EndSequence()
}

func DecodeValue(reader *asn1.ASNReader, ctxt uint8) (*ContentParameter, errors.Error) {
	var contentParameter ContentParameter
	pcLength, pcReader, err := reader.ReadSequenceStart(asn1.Context(ctxt))
	if err != nil {
		return nil, errors.Update(err)
	}
	if pcLength == 0 {
		return &contentParameter, nil
	}
	pcType, err := pcReader.Peek()
	if err != nil {
		return nil, errors.Update(err)
	}

	switch pcType {
	case asn1.EMBER_BOOLEAN:
		b, err := pcReader.ReadBoolean()
		if err != nil {
			return nil, errors.Update(err)
		}
		contentParameter.SetBool(b)
		break
	case asn1.EMBER_INTEGER:
		i, err := pcReader.ReadInt64()
		if err != nil {
			return nil, errors.Update(err)
		}
		contentParameter.SetInt(i)
		break
	case asn1.EMBER_REAL:
		r, err := pcReader.ReadReal()
		if err != nil {
			return nil, errors.Update(err)
		}
		contentParameter.SetReal(r)
		break
	case asn1.EMBER_BITSTRING:
		s, err := pcReader.ReadBitString()
		if err != nil {
			return nil, errors.Update(err)
		}
		contentParameter.SetBuffer(s)
		break
	case asn1.EMBER_STRING:
		//case asn1.EMBER_OCTETSTRING:
		s, err := pcReader.ReadString()
		if err != nil {
			return nil, errors.Update(err)
		}
		contentParameter.SetString(s)
		break
	default:
		return nil, errors.New("Unknown value type %d.", pcType)
	}
	err = pcReader.ReadSequenceEnd()
	return &contentParameter, errors.Update(err)
}
