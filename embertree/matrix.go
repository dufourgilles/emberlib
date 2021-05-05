package embertree

import (
	"fmt"

	"github.com/dufourgilles/emberlib/asn1"
	"github.com/dufourgilles/emberlib/errors"
)

const (
	matrixTypeCtx               = 2
	matrixModeCtx               = iota
	targetCountCtx              = iota
	sourceCounCtx               = iota
	maximumTotalConnectsCtx     = iota
	maximumConnectsPerTargetCtx = iota
	parametersLocationCtx       = iota
	gainParameterNumberCtx      = iota
	matrixContentSize           = iota
)

type MatrixType int8

const (
	OneToN   MatrixType = 0
	OneToOne MatrixType = iota
	NToN     MatrixType = iota
)

type MatrixMode int8

const (
	Linear    MatrixMode = 0
	NonLinear MatrixMode = iota
)

var (
	MatrixApplication        = asn1.Application(13)
	labelContext             = asn1.Context(10)
	templateReferenceContext = asn1.Context(12)
)

type MatrixContent struct {
	labels            []*Label
	schemaIdentifier  ContentParameter
	templateReference asn1.RelativeOID
	table             [matrixContentSize]ContentParameter
}

func ValidateMatrixType(mtype MatrixType) errors.Error {
	if mtype < OneToN || mtype > NToN {
		return errors.New("Invalid matrix type %d", mtype)
	}
	return nil
}

func ValidateMatrixMode(mode MatrixMode) errors.Error {
	if mode < Linear || mode > NonLinear {
		return errors.New("Invalid matrix mode %d", mode)
	}
	return nil
}

func (c *MatrixContent) SetType(mtype MatrixType) errors.Error {
	err := ValidateMatrixType(mtype)
	if err != nil {
		return errors.Update(err)
	}
	c.table[matrixTypeCtx].SetInt(int64(mtype))
	return nil
}

func (contents *MatrixContent) SetIdentifier(identifer string) {
	contents.table[identifierCtx].SetString(identifer)
}

func (contents *MatrixContent) GetIdentifier() (string, errors.Error) {
	return contents.table[identifierCtx].GetString()
}

func (contents *MatrixContent) SetDescription(description string) {
	contents.table[descriptionCtx].SetString(description)
}

func (contents *MatrixContent) GetDescription() (string, errors.Error) {
	return contents.table[descriptionCtx].GetString()
}

func (c *MatrixContent) GetType() (MatrixType, errors.Error) {
	v, err := c.table[matrixTypeCtx].GetInt()
	if err != nil {
		return 0, errors.Update(err)
	}
	return MatrixType(v), nil
}

func (c *MatrixContent) SetMode(mode MatrixMode) errors.Error {
	err := ValidateMatrixMode(mode)
	if err != nil {
		return errors.Update(err)
	}
	c.table[matrixModeCtx].SetInt(int64(mode))
	return nil
}

func (c *MatrixContent) GetMode() (MatrixMode, errors.Error) {
	v, err := c.table[matrixModeCtx].GetInt()
	if err != nil {
		return 0, errors.Update(err)
	}
	return MatrixMode(v), nil
}

func (c *MatrixContent) SetTargetCount(count int) errors.Error {
	c.table[targetCountCtx].SetInt(int64(count))
	return nil
}

func (c *MatrixContent) GetTargetCount() (int, errors.Error) {
	v, err := c.table[targetCountCtx].GetInt()
	if err != nil {
		return 0, errors.Update(err)
	}
	return int(v), nil
}

func (c *MatrixContent) SetSourceCount(count int) errors.Error {
	c.table[sourceCounCtx].SetInt(int64(count))
	return nil
}

func (c *MatrixContent) GetSourceCount() (int, errors.Error) {
	v, err := c.table[sourceCounCtx].GetInt()
	if err != nil {
		return 0, errors.Update(err)
	}
	return int(v), nil
}

func (c *MatrixContent) SetMaxTotalConnects(count int) errors.Error {
	c.table[maximumTotalConnectsCtx].SetInt(int64(count))
	return nil
}

func (c *MatrixContent) GetMaxTotalConnects() (int, errors.Error) {
	v, err := c.table[maximumTotalConnectsCtx].GetInt()
	if err != nil {
		return 0, errors.Update(err)
	}
	return int(v), nil
}

func (c *MatrixContent) SetMaxConnectsPerTarget(count int) errors.Error {
	c.table[maximumConnectsPerTargetCtx].SetInt(int64(count))
	return nil
}

func (c *MatrixContent) GetMaxConnectsPerTarget() (int, errors.Error) {
	v, err := c.table[maximumConnectsPerTargetCtx].GetInt()
	if err != nil {
		return 0, errors.Update(err)
	}
	return int(v), nil
}

func (c *MatrixContent) SetParameterLocation(oid asn1.RelativeOID) errors.Error {
	c.table[parametersLocationCtx].SetRelativeOID(oid)
	return nil
}

func (c *MatrixContent) GetParameterLocation() (asn1.RelativeOID, errors.Error) {
	v, err := c.table[parametersLocationCtx].GetRelativeOID()
	if err != nil {
		return nil, errors.Update(err)
	}
	return v, nil
}

func (c *MatrixContent) SetGainParameterNumber(count int) errors.Error {
	c.table[gainParameterNumberCtx].SetInt(int64(count))
	return nil
}

func (c *MatrixContent) GetGainParameterNumber() (int, errors.Error) {
	v, err := c.table[gainParameterNumberCtx].GetInt()
	if err != nil {
		return 0, errors.Update(err)
	}
	return int(v), nil
}

func (c *MatrixContent) SetSchemaIdentifier(schema string) errors.Error {
	c.schemaIdentifier.SetString(schema)
	return nil
}

func (c *MatrixContent) GetSchemaIdentifier() (string, errors.Error) {
	v, err := c.schemaIdentifier.GetString()
	if err != nil {
		return "", errors.Update(err)
	}
	return v, nil
}

func NewMatrixContent(mtype MatrixType, mode MatrixMode) (*MatrixContent, errors.Error) {
	content := MatrixContent{}
	err := content.SetType(mtype)
	if err != nil {
		return nil, errors.Update(err)
	}
	err = content.SetMode(mode)
	if err != nil {
		return nil, errors.Update(err)
	}
	return &content, nil
}

func NewDefaultMatrixContents() EmberContents {
	mc, _ := NewMatrixContent(OneToN, Linear)
	return mc
}

func NewMatrix(number int, mtype MatrixType, mode MatrixMode) (*Element, errors.Error) {
	content, err := NewMatrixContent(mtype, mode)
	if err != nil {
		return nil, errors.Update(err)
	}
	element := NewElement(MatrixApplication, number, nil)
	element.SetContents(content)
	element.isMatrix = true
	return element, nil
}

func (c *MatrixContent) EncodeLabels(writer *asn1.ASNWriter) errors.Error {
	err := writer.StartSequence(labelContext)
	if err != nil {
		return errors.Update(err)
	}
	err = writer.StartSequence(asn1.EMBER_SEQUENCE)
	if err != nil {
		return errors.Update(err)
	}
	for _, label := range c.labels {
		err = writer.StartSequence(asn1.Context(0))
		if err != nil {
			return errors.Update(err)
		}
		err = label.Encode(writer)
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

func (c *MatrixContent) DecodeLabels(reader *asn1.ASNReader) errors.Error {
	var labels []*Label
	_, labelReader, err := reader.ReadSequenceStart(labelContext)
	if err != nil {
		return errors.Update(err)
	}
	_, seqReader, err := labelReader.ReadSequenceStart(asn1.EMBER_SEQUENCE)
	if err != nil {
		return errors.Update(err)
	}
	for seqReader.Len() > 0 {
		_, ctxtReader, err := seqReader.ReadSequenceStart(asn1.Context(0))
		if err != nil {
			return errors.Update(err)
		}
		label, err := DecodeLabel(ctxtReader)
		if err != nil {
			return errors.Update(err)
		}
		labels = append(labels, label)
		err = ctxtReader.ReadSequenceEnd()
		if err != nil {
			return errors.Update(err)
		}
		end, err := seqReader.CheckSequenceEnd()
		if end {
			break
		}
		if err != nil {
			return errors.Update(err)
		}
	}
	err = reader.ReadSequenceEnd()
	if err != nil {
		return errors.Update(err)
	}
	c.labels = labels
	return nil
}

func (c *MatrixContent) Encode(writer *asn1.ASNWriter) errors.Error {
	err := writer.StartSequence(asn1.EMBER_SET)
	if err != nil {
		return errors.Update(err)
	}

	for i, cp := range c.table {
		err = cp.Encode(uint8(i), writer)
		if err != nil {
			return errors.Update(err)
		}
	}

	err = c.EncodeLabels(writer)
	if err != nil {
		return errors.Update(err)
	}

	err = c.schemaIdentifier.Encode(11, writer)
	if err != nil {
		return errors.Update(err)
	}

	if c.templateReference != nil {
		err = writer.StartSequence(templateReferenceContext)
		if err != nil {
			return errors.Update(err)
		}
		err = writer.WriteRelativeOID(c.templateReference)
		if err != nil {
			return errors.Update(err)
		}
		err = writer.EndSequence()
		if err != nil {
			return errors.Update(err)
		}
	}

	return writer.EndSequence()
}

func (c *MatrixContent) Decode(reader *asn1.ASNReader) errors.Error {
	_, matrixContentReader, err := reader.ReadSequenceStart(asn1.EMBER_SET)
	if err != nil {
		return errors.Update(err)
	}

	for matrixContentReader.Len() > 0 {
		peek, err := matrixContentReader.Peek()
		if err != nil {
			return errors.Update(err)
		}
		index := uint8(peek) - asn1.Context(0)
		if index < matrixContentSize {
			value, err := DecodeValue(matrixContentReader, index)
			if err != nil {
				return errors.Update(err)
			}
			err = c.table[index].Set(value)
			if err != nil {
				return errors.Update(err)
			}
		} else if peek == labelContext {
			c.DecodeLabels(matrixContentReader)
		} else if index == 11 {
			value, err := DecodeValue(matrixContentReader, index)
			if err != nil {
				return errors.Update(err)
			}
			err = c.schemaIdentifier.Set(value)
			if err != nil {
				return errors.Update(err)
			}
		} else if peek == templateReferenceContext {
			_, templateReader, err := matrixContentReader.ReadSequenceStart(peek)
			templateReference, err := templateReader.ReadOID(asn1.EMBER_RELATIVE_OID)
			if err != nil {
				return errors.Update(err)
			}
			c.templateReference = templateReference
			err = templateReader.ReadSequenceEnd()
			if err != nil {
				return errors.Update(err)
			}
		}
		end, err := matrixContentReader.CheckSequenceEnd()
		if end {
			break
		}
		if err != nil {
			return errors.Update(err)
		}
	}

	return nil
}


func (c *MatrixContent) ToString() string {
	str:= ""
	valStr,err := c.GetIdentifier()
	if err == nil {
		str = fmt.Sprintf("%s  identifier: %s\n",str, valStr)
	}
	valStr,err = c.GetDescription()
	if err == nil {
		str = fmt.Sprintf("%s  description: %s\n",str, valStr)
	}
	t,err := c.GetType()
	if t == OneToN {
		str = fmt.Sprintf(("%s  type: OneToN"))
	} else if t ==  NToN {
		str = fmt.Sprintf(("%s  type: NToN"))
	} else {
		str = fmt.Sprintf(("%s  type: OneToOne"))
	}
	m,err := c.GetMode()
	if m == Linear {
		str = fmt.Sprintf(("%s  mode: Linear"))
	} else {
		str = fmt.Sprintf(("%s  mode: Non-Linear"))
	}
	return fmt.Sprintf("{\n%s}\n", str)
}