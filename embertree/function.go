package embertree

import (
	"fmt"

	"github.com/dufourgilles/emberlib/errors"

	"github.com/dufourgilles/emberlib/asn1"
)

type ParameterType uint8

const (
	ParameterTypeNull    ParameterType = 0
	ParameterTypeInteger ParameterType = iota
	ParameterTypeReal    ParameterType = iota
	ParameterTypeString  ParameterType = iota
	ParameterTypeBoolean ParameterType = iota
	ParameterTypeTrigger ParameterType = iota
	ParameterTypeEnum    ParameterType = iota
	ParameterTypeOcets   ParameterType = iota
)

type ParamaterAccess uint8

const (
	ParamaterAccessNone      ParamaterAccess = 0
	ParamaterAccessRead      ParamaterAccess = iota
	ParamaterAccessWrite     ParamaterAccess = iota
	ParamaterAccessReadWrite ParamaterAccess = iota
)

type TupleDescription struct {
	Type ParameterType
	Name      string
}

type FunctionContents struct {
	identifier        ContentParameter
	description       ContentParameter
	arguments         []*TupleDescription
	result            []*TupleDescription
	templateReference asn1.RelativeOID
}

var FunctionApplication = asn1.Application(19)
var TupleDescriptionApplication = asn1.Application(21)

func NewArgument(t ParameterType, name string) *TupleDescription {
	return &TupleDescription{Type: t, Name: name}
}

func NewFunctionContents() EmberContents {
	return &FunctionContents{}
}

func NewFunction(number int) *Element {
	return NewElement(FunctionApplication, number, NewFunctionContents)
}

func (contents *FunctionContents) SetIdentifier(identifer string) {
	contents.identifier.SetString(identifer)
}

func (contents *FunctionContents) GetIdentifier() (string, errors.Error) {
	return contents.identifier.GetString()
}

func (contents *FunctionContents) SetDescription(description string) {
	contents.description.SetString(description)
}

func (contents *FunctionContents) GetDescription() (string, errors.Error) {
	return contents.description.GetString()
}

func (contents *FunctionContents) SetArguments(arguments []*TupleDescription) {
	contents.arguments = arguments
}

func (contents *FunctionContents) GetArguments() []*TupleDescription {
	return contents.arguments
}

func (contents *FunctionContents) SetResult(result []*TupleDescription) {
	contents.result = result
}

func (contents *FunctionContents) GetResult() []*TupleDescription {
	return contents.result
}

func decodeTupleDescriptions(reader *asn1.ASNReader) ([]*TupleDescription, errors.Error) {
	arguments := []*TupleDescription{}
	_, seqReader, err := reader.ReadSequenceStart(asn1.EMBER_SEQUENCE)
	if err != nil {
		return nil, errors.Update(err)
	}
	for seqReader.Len() > 0 {
		_, argReader, err := seqReader.ReadSequenceStart(asn1.Context(0))
		if err != nil {
			return nil, errors.Update(err)
		}
		tuple := &TupleDescription{}
		err = tuple.Decode(argReader)
		if err != nil {
			return nil, errors.Update(err)
		}
		arguments = append(arguments, tuple)
		err = argReader.ReadSequenceEnd()
		if err != nil {
			return nil, errors.Update(err)
		}
		end, err := seqReader.CheckSequenceEnd()
		if end {
			break
		}
		if err != nil {
			return nil, errors.Update(err)
		}
	}
	return arguments, nil
}

func encodeTupleDescriptions(ctxt uint8, tuples []*TupleDescription, writer *asn1.ASNWriter) errors.Error {
	err := writer.StartSequence(ctxt)
	if err != nil {
		return errors.Update(err)
	}
	err = writer.StartSequence(asn1.EMBER_SEQUENCE)
	if err != nil {
		return errors.Update(err)
	}
	for _, tuple := range tuples {
		err = writer.StartSequence(asn1.Context(0))
		if err != nil {
			return errors.Update(err)
		}
		err = tuple.Encode(writer)
		if err != nil {
			return errors.Update(err)
		}
		err = writer.EndSequence()
	if err != nil {
		return errors.Update(err)
	}
	}
	err = writer.EndSequence()
	if err != nil {
		return errors.Update(err)
	}
	return writer.EndSequence()
}

func (tuple *TupleDescription) Decode(reader *asn1.ASNReader) errors.Error {
	_, tupleReader, err := reader.ReadSequenceStart(TupleDescriptionApplication)
	if err != nil {
		return errors.Update(err)
	}
	for tupleReader.Len() > 0 {
		peek, err := tupleReader.Peek()
		if err != nil {
			return errors.Update(err)
		}
		_, ctxtReader, err := tupleReader.ReadSequenceStart(peek)
		if err != nil {
			return errors.Update(err)
		}
		switch peek {
		case asn1.Context(0):
			t, err := ctxtReader.ReadInt()
			if err != nil {
				return errors.Update(err)
			}
			tuple.Type = ParameterType(t)
			break
		case asn1.Context(1):
			name, err := ctxtReader.ReadString()
			if err != nil {
				return errors.Update(err)
			}
			tuple.Name = name
			break
		default:
			return errors.New("Unknown TupleDescription tag %d", peek)
		}
		err = ctxtReader.ReadSequenceEnd()
		if err != nil {
			return errors.Update(err)
		}
		end, err := tupleReader.CheckSequenceEnd()
		if end || err != nil {
			break
		}
	}
	return errors.Update(err)
}

func (tuple *TupleDescription) Encode(writer *asn1.ASNWriter) errors.Error {
	err := writer.StartSequence(TupleDescriptionApplication)
	if err != nil {
		return errors.Update(err)
	}

	err = writer.StartSequence(asn1.Context(0))
	if err != nil {
		return errors.Update(err)
	}
	err = writer.WriteInt(int(tuple.Type))
	if err != nil {
		return errors.Update(err)
	}
	err = writer.EndSequence()
	if err != nil {
		return errors.Update(err)
	}

	err = writer.StartSequence(asn1.Context(1))
	if err != nil {
		return errors.Update(err)
	}
	err = writer.WriteString(tuple.Name)
	if err != nil {
		return errors.Update(err)
	}
	err = writer.EndSequence()
	if err != nil {
		return errors.Update(err)
	}

	err = writer.EndSequence()
	if err != nil {
		return errors.Update(err)
	}
	return nil
}

func (fc *FunctionContents) Decode(reader *asn1.ASNReader) errors.Error {
	var err errors.Error
	var value *ContentParameter
	var ctxtReader *asn1.ASNReader
	_, reader, err = reader.ReadSequenceStart(asn1.EMBER_SET)
	if err != nil {
		return errors.Update(err)
	}
	for reader.Len() > 0 {
		peek, err := reader.Peek()
		if err != nil {
			return errors.Update(err)
		}
		if peek < asn1.Context(2) {
			value, err = DecodeValue(reader, peek)
		} else {
			_, ctxtReader, err = reader.ReadSequenceStart(peek)
		}
		if err != nil {
			return errors.Update(err)
		}
		switch peek {
		case asn1.Context(0):
			fc.identifier.Set(value)
			break
		case asn1.Context(1):
			fc.description.Set(value)
			break
		case asn1.Context(2):
			arguments, err := decodeTupleDescriptions(ctxtReader)
			if err != nil {
				return errors.Update(err)
			}
			fc.arguments = arguments
			err = ctxtReader.ReadSequenceEnd()
			if err != nil {
				return errors.Update(err)
			}
			break
		case asn1.Context(3):
			result, err := decodeTupleDescriptions(ctxtReader)
			if err != nil {
				return errors.Update(err)
			}
			fc.result = result
			err = ctxtReader.ReadSequenceEnd()
			if err != nil {
				return errors.Update(err)
			}
			break
		default:
			return errors.New("Unknown function content tag %d", peek)
		}
		end, err := reader.CheckSequenceEnd()
		if end || err != nil {
			break
		}
	}
	return errors.Update(err)
}

func (fc *FunctionContents) Encode(writer *asn1.ASNWriter) errors.Error {
	err := writer.StartSequence(asn1.EMBER_SET)
	if err != nil {
		return errors.Update(err)
	}
	if fc.identifier.IsSet() {
		err = fc.identifier.Encode(0, writer)
		if err != nil {
			return errors.Update(err)
		}
	}

	if fc.description.IsSet() {
		fc.description.Encode(1, writer)
		if err != nil {
			return errors.Update(err)
		}
	}

	if len(fc.arguments) > 0 {
		err = encodeTupleDescriptions(asn1.Context(2), fc.arguments, writer)
		if err != nil {
			return errors.Update(err)
		}
	}
	
	if len(fc.result) > 0 {
		err = encodeTupleDescriptions(asn1.Context(3), fc.result, writer)
		if err != nil {
			return errors.Update(err)
		}
	}
	writer.EndSequence()
	return nil
}

func (fc *FunctionContents) ToString() string {
	str:= ""
	valStr,err := fc.GetIdentifier()
	if err == nil {
		str = fmt.Sprintf("%s  identifier: %s\n",str, valStr)
	}
	valStr,err = fc.GetDescription()
	if err == nil {
		str = fmt.Sprintf("%s  description: %s\n",str, valStr)
	}
	return fmt.Sprintf("{\n%s}\n", str)
}