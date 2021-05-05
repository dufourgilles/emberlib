package embertree_test

import (
	"fmt"
	"testing"

	"github.com/dufourgilles/emberlib/asn1"
	"github.com/dufourgilles/emberlib/embertree"
	"github.com/dufourgilles/emberlib/errors"
)

func TestEncodeRoot(t *testing.T) {
	nodeID := int(10)
	node := embertree.NewNode(nodeID)
	nodeContents := node.CreateContent().(*embertree.NodeContents)
	nodeContents.SetIdentifier("gdnet")
	root := embertree.NewRoot()
	root.AddElement(node)
	writer := asn1.ASNWriter{}
	err := root.Encode(&writer)
	if err != nil {
		fmt.Println(err.Message)
		fmt.Println(err.Stack)
		t.Error(err)
		return
	}
	b := make([]byte, writer.Len())
	writer.Read(b)
	for _, bb := range b {
		fmt.Printf("%x ", bb)
	}
	expectedResult := []byte{0x60, 0x80, 0x6b, 0x80, 0xa0, 0x80, 0x63, 0x80, 0xa0, 0x80, 02, 01, 0x0a, 0, 0, 0xa1,
		0x80, 0x31, 0x80, 0xa0, 0x80, 0x0c, 05, 0x67, 0x64, 0x6e, 0x65, 0x74, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i, d := range b {
		if d != expectedResult[i] {
			t.Errorf("Invalid byte at %d.  Val %d expected %d.", i, d, expectedResult[i])
		}
	}
}

func TestDecodeRoot(t *testing.T) {
	expectedResult := []byte{0x60, 0x1d, 0x6b, 0x1b, 0xa0, 0x19, 0x63, 0x17, 0xa0, 03, 02, 01, 0x0a, 0xa1,
		0x10, 0x31, 0x0e, 0xa0, 07, 0x0c, 05, 0x67, 0x64, 0x6e, 0x65, 0x74, 0xa3, 03, 01, 01, 0xFF}
	reader := asn1.NewASNReader(expectedResult)
	root := embertree.NewTree()
	err := root.Decode(reader)
	if err != nil {
		fmt.Println(err.Message)
		fmt.Println(err.Stack)
		t.Error(err)
		return
	}
	fmt.Println(root)
	if len(root.RootElementCollection) != 1 {
		t.Errorf("Invalid Element size %d", len(root.RootElementCollection))
		return
	}
	node := root.RootElementCollection[10]
	if node == nil {
		t.Errorf("Invalid Element at 10")
		return
	}
	nodeContents := node.GetContent().(*embertree.NodeContents)
	if nodeContents == nil {
		t.Errorf("Invalid Node Content")
		return
	}
	identifier, err := nodeContents.GetIdentifier()
	if identifier != "gdnet" {
		t.Errorf("Invalid Node Content identifier")
		return
	}
}

func TestDecodeRootWithChildren(t *testing.T) {
	expectedResult := []byte{96, 61, 107, 59, 160, 57, 99, 55, 160, 3, 2, 1, 10, 161, 16, 49, 14,
		160, 7, 12, 5, 103, 100, 110, 101, 116, 163, 3, 1, 1, 255,
		162, 30, 100, 28, 160, 26, 97, 24, 160, 3, 2, 1, 10,
		161, 17, 49, 15,
		160, 7, 12, 5, 103, 100, 110, 101, 116, 162, 4, 2, 2, 4}
	reader := asn1.NewASNReader(expectedResult)
	root := embertree.NewTree()
	err := root.Decode(reader)
	if err != nil {
		fmt.Println(err.Message)
		fmt.Println(err.Stack)
		t.Error(err)
		return
	}
	if len(root.RootElementCollection) != 1 {
		t.Errorf("Invalid Element size %d", len(root.RootElementCollection))
		return
	}
	node := root.RootElementCollection[10]
	if node == nil {
		t.Errorf("Invalid Element at 10")
		return
	}
	nodeContents := node.GetContent().(*embertree.NodeContents)
	if nodeContents == nil {
		t.Errorf("Invalid Node Content")
		return
	}
	identifier, err := nodeContents.GetIdentifier()
	if identifier != "gdnet" {
		t.Errorf("Invalid Node Content identifier")
		return
	}
}

//96,16,107,14,160,12,98,10,160,3,2,1,32,161,3,2,1
func TestDecodeRootGetDirectory(t *testing.T) {
	expectedResult := []byte{96, 16, 107, 14, 160, 12, 98, 10, 160, 3, 2, 1, 32, 161, 3, 2, 1, 0xff}
	reader := asn1.NewASNReader(expectedResult)
	root := embertree.NewTree()
	err := root.Decode(reader)
	if err != nil {
		fmt.Println(err.Message)
		fmt.Println(err.Stack)
		t.Error(err)
		return
	}
	if len(root.RootElementCollection) != 1 {
		t.Errorf("Invalid Element size %d", len(root.RootElementCollection))
		return
	}
	cmd := root.RootElementCollection[32]
	if cmd == nil {
		t.Errorf("Invalid Element at 32")
		return
	}
	cmdContents := cmd.GetContent().(*embertree.CommandContents)
	if cmdContents == nil {
		t.Errorf("Invalid Command Content")
		return
	}
	if cmd.Number != 32 {
		t.Errorf("Invalid Command Number")
		return
	}
}

func TestDecodeRootMatrix(t *testing.T) {
	expectedResult := []byte{96, 82, 107, 80, 160, 78, 109, 76, 160, 3, 2, 1, 1, 163, 29, 48, 27, 160, 7, 110, 5, 160, 3, 2, 1, 1, 160, 7, 110, 5, 160, 3, 2, 1, 2, 160, 7, 110, 5, 160, 3, 2, 1, 3, 164, 20, 48, 18, 160, 7, 111, 5, 160, 3, 2, 1, 1, 160, 7, 111, 5, 160, 3, 2, 1, 2, 165, 16, 48, 14, 160, 12, 112, 10, 160, 3, 2, 1, 1, 161, 3, 13, 1, 2}
	reader := asn1.NewASNReader(expectedResult)
	root := embertree.NewTree()
	err := root.Decode(reader)
	if err != nil {
		fmt.Println(err.Message)
		fmt.Println(err.Stack)
		t.Error(err)
		return
	}
	if len(root.RootElementCollection) != 1 {
		t.Errorf("Invalid Element size %d", len(root.RootElementCollection))
		return
	}
	matrix := root.RootElementCollection[1]
	if matrix == nil {
		t.Errorf("Invalid Element at 1")
		return
	}
	targets, err := matrix.GetTargets()
	if err != nil {
		fmt.Println(err.Message)
		fmt.Println(err.Stack)
		t.Error(err)
		return
	}
	encodedTargets := []int32{1, 2, 3}
	for i, signal := range targets {
		target := signal.(*embertree.Target)
		if target.Number != encodedTargets[i] {
			t.Errorf("Invalid target at %d. Got %d instead of %d.", i, target.Number, encodedTargets[i])
		}
	}
}

type RootListener struct {

}

func (rl *RootListener) Receive(interface{}, errors.Error) {

}

func TestListeners(t *testing.T) {
	listener := RootListener{}
	root := embertree.NewTree()
	root.AddListener(&listener)
	if !root.HasListner(&listener) {
		t.Errorf("Add Listener failed")
	}
	root.RemoveListener(&listener)
	if root.HasListner(&listener) {
		t.Errorf("Remove Listener failed")
	}
}