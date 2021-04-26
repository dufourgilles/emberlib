package embertree_test

import (
	"fmt"
	"testing"

	"github.com/dufourgilles/emberlib/asn1"
	"github.com/dufourgilles/emberlib/embertree"
	"github.com/dufourgilles/emberlib/errors"
)

func TestEncodeNode(t *testing.T) {
	nodeID := int(10)
	node := embertree.NewNode(nodeID)
	writer := asn1.ASNWriter{}
	err := node.Encode(&writer)
	if err != nil {
		t.Error(err)
	}
	b := make([]byte, writer.Len())
	writer.Read(b)
	expectedResult := []byte{99, 128, 160, 128, 2, 1, 10, 0, 0, 0, 0}
	for i, d := range b {
		if d != expectedResult[i] {
			t.Errorf("Invalid byte at %d.  Val %d expected %d.", i, d, expectedResult[i])
		}
	}
}

func TestEncodeNodeIdentifier(t *testing.T) {
	nodeID := int(10)
	node := embertree.NewNode(nodeID)
	nodeContents := node.CreateContent().(*embertree.NodeContents)
	nodeContents.SetIdentifier("gdnet")
	writer := asn1.ASNWriter{}
	err := node.Encode(&writer)
	if err != nil {
		fmt.Println(err.Message)
		fmt.Println(err.Stack)
		t.Error(err)
		return
	}
	b := make([]byte, writer.Len())
	writer.Read(b)
	//fmt.Println(b)
	expectedResult := []byte{99, 128, 160, 128, 2, 1, 10, 0, 0, 161, 128, 49, 128, 160, 128, 12, 5, 103, 100, 110, 101, 116, 0, 0, 0, 0, 0, 0, 0, 0}
	for i, d := range b {
		if d != expectedResult[i] {
			t.Errorf("Invalid byte at %d.  Val %d expected %d.", i, d, expectedResult[i])
		}
	}
}

type listenerTest struct {
	el  interface{}
	err errors.Error
}

func (l *listenerTest) listener(element interface{}, err errors.Error) {
	l.el = element
	l.err = err
}

func TestNodeListener(t *testing.T) {
	listener := listenerTest{el: nil, err: nil}
	nodeID := int(10)
	node := embertree.NewNode(nodeID)
	node.GetDirectoryMsg(listener.listener)
	nodeContents := node.CreateContent().(*embertree.NodeContents)
	nodeContents.SetIdentifier("gdnet")
	node.SetContents(nodeContents)
	if listener.el == nil && listener.err == nil {
		t.Errorf("Listener Failed to receive node")
		return
	}
}
