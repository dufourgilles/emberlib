package embertree

import (
	"github.com/dufourgilles/emberlib/asn1"
	"github.com/dufourgilles/emberlib/errors"
)

type EmberContents interface {
	Encode(writer *asn1.ASNWriter) errors.Error
	Decode(reader *asn1.ASNReader) errors.Error
}

type ContentCreator func() EmberContents

type EmberObject interface {
	Encode(writer *asn1.ASNWriter) errors.Error
	CreateContent() interface{}
	GetContent() interface{}
	GetParent() (EmberObject, errors.Error)
	AddChild(child interface{}) errors.Error
	GetTag() uint8
}

type Listener func(interface{}, errors.Error)

type Element struct {
	Number          int
	tag             uint8
	path            asn1.RelativeOID
	children        map[int]*Element
	parent          *Element
	contents        EmberContents
	contentsCreator ContentCreator
	listeners       map[*Listener]Listener

	// for qualified element
	isQualified bool

	// for matrix only
	isMatrix    bool
	targets     []Signal
	sources     []Signal
	connections []*Connection
}

func NewElement(tag uint8, number int, contentsCreator ContentCreator) *Element {
	return &Element{
		Number:          number,
		tag:             tag,
		contents:        nil,
		parent:          nil,
		children:        make(map[int]*Element),
		isMatrix:        false,
		isQualified:     false,
		contentsCreator: contentsCreator,
		listeners:       make(map[*Listener]Listener)}
}

func (element *Element) CreateContent() interface{} {
	if element.contentsCreator == nil {
		return nil
	}
	element.contents = element.contentsCreator().(EmberContents)
	return element.contents
}

func (element *Element) GetContent() interface{} {
	return element.contents
}

func (element *Element) GetTag() uint8 {
	return element.tag
}

func (element *Element) AddChild(child *Element) errors.Error {
	element.children[child.Number] = child
	child.SetParent(element)
	return nil
}

func (element *Element) updateListeners(err errors.Error) {
	for _, listener := range element.listeners {
		listener(element, err)
	}
}
func (element *Element) SetContents(contents interface{}) errors.Error {
	element.contents = contents.(EmberContents)
	element.updateListeners(nil)
	return nil
}

func (element *Element) Update(newElement *Element) errors.Error {
	var err errors.Error
	if element.Number != newElement.Number || element.tag != newElement.tag {
		return errors.New("Attempt to update different element Number %d/%d Tag %d/%d", element.Number, newElement.Number, element.tag, newElement.tag)
	}
	content := newElement.GetContent()
	if content != nil {
		element.contents = content.(EmberContents)
	}
	for number, newChild := range newElement.children {
		child := element.children[number]
		if child == nil {
			err = element.AddChild(newChild)
		} else {
			err = child.Update(newChild)
		}
		if err != nil {
			break
		}
	}
	element.updateListeners(err)
	return err
}

func (element *Element) GetContents() (EmberContents, errors.Error) {
	return element.contents, nil
}

func (element *Element) SetParent(parent *Element) errors.Error {
	element.parent = parent
	if !element.isQualified {
		element.path = asn1.RelativeOID{}
	}
	return nil
}

func (element *Element) GetParent() (*Element, errors.Error) {
	return element.parent, nil
}

func (element *Element) encode(writer *asn1.ASNWriter, asChild bool) errors.Error {
	//fmt.Printf("Starting element seq %d\n", element.seqID)
	err := writer.StartSequence(element.tag)
	if err != nil {
		return errors.Update(err)
	}

	if element.isQualified {
		//Encode Path
		err = writer.StartSequence(asn1.Context(0))
		if err != nil {
			return errors.Update(err)
		}
		err = writer.WriteRelativeOID(element.path)
		if err != nil {
			return errors.Update(err)
		}
		err = writer.EndSequence()
		if err != nil {
			return errors.Update(err)
		}
	} else {
		//Encode Number
		err = writer.StartSequence(asn1.Context(0))
		if err != nil {
			return errors.Update(err)
		}
		err = writer.WriteInt(int(element.Number))
		if err != nil {
			return errors.Update(err)
		}
		err = writer.EndSequence()
		if err != nil {
			return errors.Update(err)
		}
	}
	// Encode Contents
	if element.contents != nil {
		err = writer.StartSequence(asn1.Context(1))
		if err != nil {
			return errors.Update(err)
		}
		err = element.contents.Encode(writer)
		if err != nil {
			return errors.Update(err)
		}
		err = writer.EndSequence()
		if err != nil {
			return errors.Update(err)
		}
	}

	if !asChild {
		//Encode Children
		if len(element.children) > 0 {
			err = writer.StartSequence(asn1.Context(2))
			if err != nil {
				return errors.Update(err)
			}
			err = element.EncodeChildren(writer)
			if err != nil {
				return errors.Update(err)
			}
			err = writer.EndSequence()
			if err != nil {
				return errors.Update(err)
			}
		}
	}

	if element.isMatrix {
		err = element.EncodeTargets(writer)
		if err != nil {
			return errors.Update(err)
		}
		err = element.EncodeSources(writer)
		if err != nil {
			return errors.Update(err)
		}
		err = element.EncodeConnections(writer)
		if err != nil {
			return errors.Update(err)
		}
	}

	return writer.EndSequence()
}

func (element *Element) Encode(writer *asn1.ASNWriter) errors.Error {
	return element.encode(writer, false)
}

func (element *Element) EncodeChildren(writer *asn1.ASNWriter) errors.Error {
	err := writer.StartSequence(asn1.Application(4))
	if err != nil {
		return errors.Update(err)
	}
	for _, child := range element.children {
		err = writer.StartSequence(asn1.Context(0))
		if err != nil {
			return errors.Update(err)
		}
		err := child.encode(writer, true)
		if err != nil {
			return err
		}
		err = writer.EndSequence()
		if err != nil {
			return errors.Update(err)
		}
	}
	return writer.EndSequence()
}

func (element *Element) getDupBranch(cmd *Element) (*Element, errors.Error) {
	e := element
	dupElement := NewElement(e.tag, e.Number, nil)
	if cmd != nil {
		dupElement.AddChild(cmd)
	}
	for {
		parent, err := e.GetParent()
		if err != nil {
			return nil, err
		}
		if parent == nil {
			break
		}
		dupParent := NewElement(parent.tag, parent.Number, nil)
		dupParent.AddChild(dupElement)
		e = parent
		dupElement = dupParent
	}
	return dupElement, nil
}

func (element *Element) GetDirectoryMsg(listener Listener) (*RootElement, errors.Error) {
	dupElement, err := element.getDupBranch(NewCommand(COMMAND_GETDIRECTORY))
	if err != nil {
		return nil, errors.Update(err)
	}
	if listener != nil {
		element.listeners[&listener] = listener
	}
	root := NewRoot()
	root.AddElement(dupElement)
	return root, nil
}
