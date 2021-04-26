package embertree

import (
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

func NewQualifiedElement(tag uint8, path asn1.RelativeOID, contentCreator ContentCreator) *Element {
	q := &Element{path: path, tag: tag, isMatrix: false, isQualified: true, contentsCreator: contentCreator, contents: contentCreator()}
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
	root := NewRoot()
	root.AddElement(dupElement)
	return root
}
