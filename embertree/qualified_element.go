package embertree

import (
	"github.com/dufourgilles/emberlib/errors"
	"github.com/dufourgilles/emberlib/asn1"
)

var QualifiedParameterApplication = asn1.Application(9)
var QualifiedNodeApplication = asn1.Application(10)
var QualifiedMatrixApplication = asn1.Application(17)
var QualifiedFunctionApplication = asn1.Application(20)

var qualifiedTags = map[uint8]bool{
	QualifiedParameterApplication: true,
	QualifiedNodeApplication:      true,
	QualifiedMatrixApplication:    true,
	QualifiedFunctionApplication:  true}


func getNumberFromPath(path asn1.RelativeOID) (int, errors.Error) {
	l := len(path)
	if l <= 0 {
		return -1, errors.New("Invalid path.")
	}
	return int(path[l-1]),nil
}

func NewQualifiedElement(tag uint8, path asn1.RelativeOID, contentCreator ContentCreator) *Element {
	number,err := getNumberFromPath(path)
	if err != nil { return nil }
	q := NewElement(tag, number, contentCreator)
	q.isQualified = true
	q.path = path
	return q
}

func NewQualifiedNode(path asn1.RelativeOID) *Element {
	return NewQualifiedElement(QualifiedNodeApplication, path, NewNodeContents)
}

func NewQualifiedParameter(path asn1.RelativeOID) *Element {
	return NewQualifiedElement(QualifiedParameterApplication, path, NewParameterContents)
}

func NewQualifiedMatrix(path asn1.RelativeOID, mtype MatrixType, mode MatrixMode) *Element {
	q := NewQualifiedElement(QualifiedMatrixApplication, path, NewDefaultMatrixContents)
	q.isMatrix = true
	return q
}

func IsQualifiedTag(tag uint8) bool {
	return qualifiedTags[tag] == true
}

func (element *Element) GetPath() asn1.RelativeOID {
	if len(element.path) <= 0 {
		var parentPath = asn1.RelativeOID{}
		if element.parent != nil {
			parentPath = element.parent.GetPath()
		}
		element.path = append(parentPath, int32(element.Number))
	}
	return element.path
}

func (element *Element) GetQualifiedDirectoryMsg(listener Listener) *RootElement {
	path := element.GetPath()
	dupElement := NewQualifiedElement(element.tag, path, nil)
	cmd := NewCommand(COMMAND_GETDIRECTORY)
	dupElement.AddChild(cmd)
	root := NewRoot()
	root.AddElement(dupElement)
	if listener != nil {
		element.AddListener(listener)
	}
	return root
}
