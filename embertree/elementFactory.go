package embertree

import (
	"fmt"

	"github.com/dufourgilles/emberlib/asn1"
	"github.com/dufourgilles/emberlib/errors"
)

func getContentCreator(tag uint8) (ContentCreator, errors.Error) {
	switch tag {
	case QualifiedParameterApplication:
		fallthrough
	case ParameterApplication:
		return NewParameterContents, nil
	case CommandApplication:
		return NewCommandContents, nil
	case QualifiedNodeApplication:
		fallthrough
	case NodeApplication:
		return NewNodeContents, nil
	case QualifiedMatrixApplication:
		fallthrough
	case MatrixApplication:
		return NewDefaultMatrixContents, nil
	case FunctionApplication:
		return NewFunctionContents, nil
	default:
		return nil, errors.New("Unknown Application %d.", tag)
	}
}

func decodeContents(element *Element, ctxt uint8, reader *asn1.ASNReader) (interface{}, errors.Error) {
	var (
		contents interface{}
	)
	_, contentReader, err := reader.ReadSequenceStart(ctxt)
	if err != nil {
		return nil, errors.Update(err)
	}
	switch element.tag {
	case QualifiedParameterApplication:
		fallthrough
	case ParameterApplication:
		pcontents := NewParameterContents().(*ParameterContents)
		err = pcontents.Decode(contentReader)
		contents = pcontents
		break
	case CommandApplication:
		ccontents := NewCommandContents().(*CommandContents)
		err = ccontents.Decode(contentReader)
		contents = ccontents
		break
	case FunctionApplication:
		fallthrough
	case QualifiedFunctionApplication:
		fcontents := NewFunctionContents().(*FunctionContents)
		err = fcontents.Decode(contentReader)
		contents = fcontents
		break
	case QualifiedNodeApplication:
		fallthrough
	case NodeApplication:
		ncontents := NewNodeContents().(*NodeContents)
		err = ncontents.Decode(contentReader)
		contents = ncontents
		break
	case QualifiedMatrixApplication:
		fallthrough
	case MatrixApplication:
		mcontents, _ := NewMatrixContent(OneToN, Linear)
		err = mcontents.Decode(contentReader)
		contents = mcontents
		break
	default:
		return nil, errors.New("Unknown Application 0x%x at offset %d.", element.tag, reader.TopOffset())
	}
	if err != nil {
		err = errors.Update(err)
		e := errors.New("Failed to decode contents tag 0x%x at offset %d. %s", element.tag, reader.TopOffset(), err.Message)
		e.Stack = err.Stack
		return nil, e
	}
	err = contentReader.ReadSequenceEnd()
	return contents, err
}

func decodeChildren(ctxt uint8, element *Element, reader *asn1.ASNReader) errors.Error {
	_, ctxtReader, err := reader.ReadSequenceStart(ctxt)
	_, childrenReader, err := ctxtReader.ReadSequenceStart(asn1.Application(4))
	if err != nil {
		return errors.Update(err)
	}
	for childrenReader.Len() > 0 {
		_, childReader, err := childrenReader.ReadSequenceStart(asn1.Context(0))
		if err != nil {
			return errors.Update(err)
		}
		child, err := DecodeElement(childReader)
		if err != nil {
			return errors.Update(err)
		}
		err = childReader.ReadSequenceEnd()
		if err != nil {
			return errors.Update(err)
		}
		element.AddChild(child)
		end, err := childrenReader.CheckSequenceEnd()
		if end {
			break
		}
		if err != nil {
			return errors.Update(err)
		}
	}
	return ctxtReader.ReadSequenceEnd()
}

func DecodeElement(reader *asn1.ASNReader) (*Element, errors.Error) {
	var (
		element  *Element
		path     asn1.RelativeOID
		number   int
		contents interface{}
	)
	tag, err := reader.Peek()
	if err != nil {
		return nil, errors.Update(err)
	}
	_, elementReader, err := reader.ReadSequenceStart(tag)
	if err != nil {
		return nil, errors.Update(err)
	}
	_, ctxtReader, err := elementReader.ReadSequenceStart(asn1.Context(0))
	if err != nil {
		return nil, errors.Update(err)
	}
	if IsQualifiedTag(tag) {
		path, err = ctxtReader.ReadOID(asn1.EMBER_RELATIVE_OID)
	} else {
		number, err = ctxtReader.ReadInt()
	}
	if err != nil {
		return nil, errors.Update(err)
	}
	err = ctxtReader.ReadSequenceEnd()
	if err != nil {
		return nil, errors.Update(err)
	}
	contentCreator, err := getContentCreator(tag)
	if err != nil {
		return nil, errors.Update(err)
	}
	if tag == MatrixApplication {
		element, _ = NewMatrix(number, OneToN, Linear)
	} else if tag == QualifiedMatrixApplication {
		element = NewQualifiedMatrix(path, OneToN, Linear)
	} else if IsQualifiedTag(tag) {
		element = NewQualifiedElement(tag, path, contentCreator)
	} else {
		element = NewElement(tag, number, contentCreator)
	}
	for elementReader.Len() > 0 {
		b, err := elementReader.Peek()
		if err != nil {
			return nil, errors.Update(err)
		}
		if b == asn1.Context(1) {
			contents, err = decodeContents(element, b, elementReader)
			if err != nil {
				fmt.Printf("Failed to decode tag %d number %d. %s", tag, number, err.Message)
				return nil, errors.Update(err)
			}
			element.SetContents(contents)
		} else if b == asn1.Context(2) {
			err = decodeChildren(b, element, elementReader)
			if err != nil {
				return nil, errors.Update(err)
			}
		} else if b == asn1.Context(3) {
			err = element.DecodeTargets(elementReader)
			if err != nil {
				return nil, errors.Update(err)
			}
		} else if b == asn1.Context(4) {
			err = element.DecodeSources(elementReader)
			if err != nil {
				return nil, errors.Update(err)
			}
		} else if b == asn1.Context(5) {
			err = element.DecodeConnections(elementReader)
			if err != nil {
				return nil, errors.Update(err)
			}
		}
		end, err := elementReader.CheckSequenceEnd()
		if end {
			break
		}
		if err != nil {
			return element, errors.Update(err)
		}
	}
	return element, nil
}
