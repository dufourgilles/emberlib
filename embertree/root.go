package embertree

import (
	"fmt"

	"github.com/dufourgilles/emberlib/errors"
	. "github.com/dufourgilles/emberlib/logger"
	"github.com/dufourgilles/emberlib/asn1"
)

type RootElement struct {
	RootElementCollection map[int]*Element
	logger Logger
	listeners             map[Listener]Listener
}

func NewTree() *RootElement {
	return &RootElement{
		listeners:             make(map[Listener]Listener),
		RootElementCollection: make(map[int]*Element),
		logger: NewNullLogger(),
	}
}

func NewRoot() *RootElement {
	return &RootElement{listeners: nil, RootElementCollection: make(map[int]*Element)}
}

func (root *RootElement) GetElementByNumber(number int) *Element {
	return root.RootElementCollection[number]
}

func (root *RootElement)SetLogger(logger Logger) {
	if logger != nil {
		root.logger = logger
	}
}

func (root *RootElement) GetElementByPath(path asn1.RelativeOID) (*Element,*Element) {
	if len(path) <= 0 {
		 return nil,nil
	}
	pos := 0
	var parent *Element
	parent = nil
	element := root.RootElementCollection[int(path[pos])]
	for pos = 1; element != nil && pos < len(path); pos++ {
		parent = element
		element = element.Children[int(path[pos])]
	}
	if pos + 1 < len(path) {
		return nil,nil
	}
	return parent,element
}

func (root *RootElement) updateQualifiedElement(element *Element) (*Element, errors.Error) {
	var err errors.Error
	parent, currentElement := root.GetElementByPath(element.path)
	if currentElement == nil {
		if parent != nil {
			parent.AddChild(element)
			return parent,err			
		} else {
			err = errors.New("Element path %s not connected to our tree\n.", Path2String(element.path))
		}
	} else {
		err = currentElement.Update(element)
	}
	return nil,err
}

func (root *RootElement) updateElement(element *Element) errors.Error {
	var err errors.Error
	currentElement := root.GetElementByNumber(element.Number)
	if currentElement == nil {
		root.AddElement(element)		
	} else {
		err = currentElement.Update(element)
	}
	return err
}

func (root *RootElement) Decode(reader *asn1.ASNReader) errors.Error {
	modifiedElement := make(map[string]*Element)
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
			root.logger.Debug("Updating E/QE %s.\n", Path2String(element.GetPath()))
			if element.isQualified && len(element.path) > 1 {
				parent,err := root.updateQualifiedElement(element)
				if err == nil && parent != nil {
					modifiedElement[Path2String(parent.GetPath())] = parent
				}
			} else {
				root.updateElement(element)
			}
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
	for _, listener := range(root.listeners) {
		root.logger.Debug("Updating root listener.\n")
		listener.Receive(root, nil)
	}
	for path,mElement := range(modifiedElement) {
		for _,listener := range(mElement.listeners) {
			root.logger.Debug("Updating Element %s listener.\n", path)
			listener.Receive(mElement, nil)
		}
	}
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
	r.AddListener(listener)
	return root, nil
}

func (r *RootElement) AddListener(listener Listener) {
	r.logger.Debug("Adding Root Listener.\n")
	r.listeners[listener] = listener
}

func (r *RootElement) RemoveListener(listener Listener) {
	r.logger.Debug("Removing Root Listener.\n")
	delete(r.listeners, listener)
}

func (r *RootElement) HasListner(listener Listener) bool {
	return r.listeners[listener] != nil
}

func (root *RootElement) ToString() string {
	str := ""
	for _,element := range(root.RootElementCollection) {
		str = fmt.Sprintf("%s%s\n", str, element.ToString())
	}
	return str
}