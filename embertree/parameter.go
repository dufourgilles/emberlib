package embertree

import (
	"github.com/dufourgilles/emberlib/errors"

	"github.com/dufourgilles/emberlib/asn1"
)

const (
	identifierCtx        = 0
	descriptionCtx       = iota
	valueCtx             = iota
	minimumCtx           = iota
	maximumCtx           = iota
	accessCtx            = iota
	formatCtx            = iota
	enumerationCtx       = iota
	factorCtx            = iota
	isOnlineCtx          = iota
	formulaCtx           = iota
	stepCtx              = iota
	defaultCtx           = iota
	typeCtx              = iota
	streamIdentifierCtx  = iota
	enumMapCtx           = iota
	streamDescriptorCtx  = iota
	schemaIdentifiersCtx = iota
	parameterContentSize = iota
)

var ParameterApplication = asn1.Application(1)

type ParameterContents struct {
	templateReference asn1.RelativeOID
	table             [parameterContentSize]ContentParameter
}

func NewParameterContents() EmberContents {
	cp := &ParameterContents{}
	return cp
}

func NewParameter(number int) *Element {
	return NewElement(ParameterApplication, number, NewParameterContents)
}

func (contents *ParameterContents) SetIdentifier(identifer string) {
	contents.table[identifierCtx].SetString(identifer)
}

func (contents *ParameterContents) GetIdentifier() (string, errors.Error) {
	return contents.table[identifierCtx].GetString()
}

func (contents *ParameterContents) SetDescription(description string) {
	contents.table[descriptionCtx].SetString(description)
}

func (contents *ParameterContents) GetDescription() (string, errors.Error) {
	return contents.table[descriptionCtx].GetString()
}

func (contents *ParameterContents) GetValueObject() *ContentParameter {
	return &contents.table[valueCtx]
}

func (contents *ParameterContents) GetMinimumObject() *ContentParameter {
	return &contents.table[minimumCtx]
}

func (contents *ParameterContents) GetMaximumObject() *ContentParameter {
	return &contents.table[maximumCtx]
}

func (contents *ParameterContents) GetDefaultObject() *ContentParameter {
	return &contents.table[defaultCtx]
}

func (contents *ParameterContents) SetAccess(access string) {
	contents.table[accessCtx].SetString(access)
}

func (contents *ParameterContents) GetAccess() (string, errors.Error) {
	return contents.table[accessCtx].GetString()
}

func (contents *ParameterContents) SetFormat(format string) {
	contents.table[formatCtx].SetString(format)
}

func (contents *ParameterContents) GetFormat() (string, errors.Error) {
	return contents.table[formatCtx].GetString()
}

func (contents *ParameterContents) SetEnumeration(enumeration string) {
	contents.table[enumerationCtx].SetString(enumeration)
}

func (contents *ParameterContents) GetEnumeration() (string, errors.Error) {
	return contents.table[enumerationCtx].GetString()
}

func (contents *ParameterContents) SetFactor(factor int64) {
	contents.table[factorCtx].SetInt(factor)
}

func (contents *ParameterContents) GetFactor() (int64, errors.Error) {
	return contents.table[factorCtx].GetInt()
}

func (contents *ParameterContents) SetOnline(online bool) {
	contents.table[isOnlineCtx].SetBool(online)
}

func (contents *ParameterContents) GetOnline() (bool, errors.Error) {
	return contents.table[isOnlineCtx].GetBool()
}

func (contents *ParameterContents) SetFormula(formula string) {
	contents.table[formulaCtx].SetString(formula)
}

func (contents *ParameterContents) GetFormula() (string, errors.Error) {
	return contents.table[formulaCtx].GetString()
}

func (contents *ParameterContents) SetStep(step int64) {
	contents.table[stepCtx].SetInt(step)
}

func (contents *ParameterContents) GetStep() (int64, errors.Error) {
	return contents.table[stepCtx].GetInt()
}

func (contents *ParameterContents) SetType(ptype string) {
	contents.table[typeCtx].SetString(ptype)
}

func (contents *ParameterContents) GetType() (string, errors.Error) {
	return contents.table[typeCtx].GetString()
}

func (contents *ParameterContents) SetStreanIdentifier(streamIdentifier int64) {
	contents.table[streamIdentifierCtx].SetInt(streamIdentifier)
}

func (contents *ParameterContents) GetStreamIdentifier() (int64, errors.Error) {
	return contents.table[streamIdentifierCtx].GetInt()
}

func (contents *ParameterContents) Encode(writer *asn1.ASNWriter) errors.Error {
	err := writer.StartSequence(asn1.EMBER_SET)
	if err != nil {
		return errors.Update(err)
	}

	for i, cp := range contents.table {
		err = cp.Encode(uint8(i), writer)
		if err != nil {
			return errors.Update(err)
		}
	}
	if contents.templateReference != nil {
		err = writer.StartSequence(asn1.Context(18))
		if err != nil {
			return errors.Update(err)
		}
		err = writer.WriteRelativeOID(contents.templateReference)
		if err != nil {
			return errors.Update(err)
		}
		err = writer.EndSequence()
		if err != nil {
			return errors.Update(err)
		}
	}
	return writer.EndSequence() // EMBER_SET
}

func (pc *ParameterContents) Decode(reader *asn1.ASNReader) errors.Error {
	var value *ContentParameter
	_, reader, err := reader.ReadSequenceStart(asn1.EMBER_SET)
	if err != nil {
		return errors.Update(err)
	}

	for reader.Len() > 0 {
		peek, err := reader.Peek()
		if err != nil {
			return errors.Update(err)
		}
		index := uint8(peek) - asn1.Context(0)
		if index < 18 {
			value, err = DecodeValue(reader, index)
			if err != nil {
				return errors.Update(err)
			}
			err = pc.table[index].Set(value)
			if err != nil {
				return errors.Update(err)
			}
		} else if index == 18 {
			_, templateReader, err := reader.ReadSequenceStart(peek)
			templateReference, err := templateReader.ReadOID(asn1.EMBER_RELATIVE_OID)
			if err != nil {
				return errors.Update(err)
			}
			pc.templateReference = templateReference
			err = templateReader.ReadSequenceEnd()
			if err != nil {
				return errors.Update(err)
			}
		}
		end, err := reader.CheckSequenceEnd()
		if end {
			break
		}
		if err != nil {
			return errors.Update(err)
		}
	}
	return nil
}
