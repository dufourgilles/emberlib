package embertree

import (
	"github.com/dufourgilles/emberlib/errors"

	"github.com/dufourgilles/emberlib/asn1"
)

type RootElement struct {
	RootElementCollection map[int]*Element
	listeners             map[*Listener]Listener
}

func NewTree() *RootElement {
	return &RootElement{
		listeners:             make(map[*Listener]Listener),
		RootElementCollection: make(map[int]*Element),
	}
}

func NewRoot() *RootElement {
	return &RootElement{listeners: nil, RootElementCollection: make(map[int]*Element)}
}

func (root *RootElement) GetElementByNumber(number int) *Element {
	return root.RootElementCollection[number]
}

func (root *RootElement) updateElement(element *Element) errors.Error {
	var err errors.Error
	currentElement := root.GetElementByNumber(element.Number)
	if currentElement == nil {
		root.AddElement(element)
	} else {
		err = currentElement.Update(element)
	}
	for _, listener := range element.listeners {
		listener(element, err)
	}
	return err
}

func (root *RootElement) Decode(reader *asn1.ASNReader) errors.Error {
	_, reader, err := reader.ReadSequenceStart(asn1.Application(0))
	if err != nil {
		return err
	}
	peek, err := reader.Peek()
	if err != nil {
		return err
	}
	if peek == asn1.Application(11) {
		_, collectionReader, err := reader.ReadSequenceStart(peek)
		if err != nil {
			return errors.Update(err)
		}
		for collectionReader.Len() > 0 {
			_, elementReader, err := collectionReader.ReadSequenceStart(asn1.Context(0))
			if err != nil {
				return errors.Update(err)
			}
			element, err := DecodeElement(elementReader)
			if err != nil {
				return errors.Update(err)
			}
			root.updateElement(element)
			err = elementReader.ReadSequenceEnd()
			if err != nil {
				return errors.Update(err)
			}
			end, err := collectionReader.CheckSequenceEnd()
			if end {
				break
			}
			if err != nil {
				return errors.Update(err)
			}
		}
	}
	err = reader.ReadSequenceEnd()
	return errors.Update(err)
}

func (root *RootElement) Encode(writer *asn1.ASNWriter) errors.Error {
	writer.StartSequence(asn1.Application(0))
	if len(root.RootElementCollection) > 0 {
		writer.StartSequence(asn1.Application(11))
		for _, element := range root.RootElementCollection {
			writer.StartSequence(asn1.Context(0))
			element.Encode(writer)
			writer.EndSequence()
		}
		writer.EndSequence()
	}
	// if (this.isRoot() && this._result != null) {
	// 	this._result.encode(ber);
	// }
	// if (this.streams != null) {
	// 	this.streams.encode(ber);
	// }
	return writer.EndSequence()
}

func (root *RootElement) AddElement(element *Element) {
	root.RootElementCollection[element.Number] = element
}

func (r *RootElement) GetDirectoryMsg(listener Listener) (*RootElement, errors.Error) {
	root := NewRoot()
	cmd := NewCommand(COMMAND_GETDIRECTORY)
	root.AddElement(cmd)
	return root, nil
}
