# emberlib 

This is a conversion of javascript/typescript [Ember+ library](https://github.com/dufourgilles/node-emberplus) in GO.
An implementation of [Lawo's Ember+](https://github.com/Lawo/ember-plus) control protocol.

This version support following ember objects : Node, Parameter, Matrix, Function, QualifiedNode,
QualifiedParameter, QualifiedMatrix, QualifiedFunction.

This is the initial version.
Encoding/Decoding results have been verified using the NodeJS library mentioned above.

## Example usage

Decode an EmberTree from Root

```go
import (
   "fmt"
   "github.com/dufourgilles/emberlib/asn1"
   "github.com/dufourgilles/emberlib/embertree"
)

encodedRoot := []byte{0x60, 0x1d, 0x6b, 0x1b, 0xa0, 0x19, 0x63, 0x17, 0xa0, 03, 02, 01, 0x0a, 0xa1,
      0x10, 0x31, 0x0e, 0xa0, 07, 0x0c, 05, 0x67, 0x64, 0x6e, 0x65, 0x74, 0xa3, 03, 01, 01, 0xFF}
reader := asn1.NewASNReader(encodedRoot)
root := embertree.NewTree()
err := root.Decode(reader)
if err != nil {
  fmt.Println(err.Message)
  fmt.Println(err.Stack)
  return
}
fmt.Println(root)
```

Create a Node and encode it

```go
import (
        "fmt"
        "github.com/dufourgilles/emberlib/asn1"
        "github.com/dufourgilles/emberlib/embertree"
)

nodeID := int(10)
node := embertree.NewNode(nodeID)
nodeContents := node.CreateContent().(*embertree.NodeContents)
nodeContents.SetIdentifier("gdnet")
writer := asn1.ASNWriter{}
err := node.Encode(&writer)
if err != nil {
  fmt.Println(err.Message)
  fmt.Println(err.Stack)
  return
}
b := make([]byte, writer.Len())
writer.Read(b)
```

Create a Parameter and encode it

```go
import (
        "fmt"
        "github.com/dufourgilles/emberlib/asn1"
        "github.com/dufourgilles/emberlib/embertree"
)

paramID := int(10)
parameter := embertree.NewParameter(paramID)
parameterContent := parameter.CreateContent().(*embertree.ParameterContents)
parameterContent.SetIdentifier("gdnet")
val := parameterContent.GetValueObject()
val.SetInt(1234)
writer := asn1.ASNWriter{}
err := parameter.Encode(&writer)
if err != nil {
  fmt.Println(err.Message)
}
b := make([]byte, writer.Len())
writer.Read(b)
```

Decode Matrix Targets

```go
import (
        "fmt"
        "github.com/dufourgilles/emberlib/asn1"
        "github.com/dufourgilles/emberlib/embertree"
)

buffer := []byte{163, 29, 48, 27, 160, 7, 110, 5, 160, 3, 2, 1, 1, 160, 7, 110, 5, 160, 3, 2, 1, 3, 160, 7, 110, 5, 160, 3, 2, 1, 5}
reader := asn1.NewASNReader(buffer)
matrix, err := NewMatrix(1, OneToN, Linear)
err = matrix.DecodeTargets(reader)
targets, err := matrix.GetTargets()
for i, signal := range targets {
  target := signal.(*Target)
        fmt.Prinf("Target %d\n", target.Number)
}
```

Get a Tree

```go
import (
   "fmt"
   "github.com/dufourgilles/emberlib/socket"
   "github.com/dufourgilles/emberlib/errors"
   "github.com/dufourgilles/emberlib/embertree"
)

type AppClient struct {
  quit chan errors.Error
}

// AppClient must have a Receive function which will be called 
// when the tree has been received
func (c *AppClient)Receive(node interface{}, err errors.Error) {
  if err != nil {
     c.quit<- err
     return
   }
   if node == nil {
      c.quit<- errors.New("nil response")
      return
    }
    root := node.(*embertree.RootElement)
    fmt.Println(root.RootElementCollection[0].ToString())
    fmt.Println(root.RootElementCollection[0])
    c.quit<- nil
}

func TestClient() {
   test := &AppClient{quit: make(chan errors.Error, 1)}
   client := socket.NewS101Client()
   client.SetTimeout(2500)
   fmt.Println("Connecting")
   err := client.Connect("192.168.1.2", 9000)
   if err != nil {
      fmt.Println(err.Message)
      return
   }
   fmt.Println("Connected. Get Tree")
   client.GetTree(test)
   fmt.Println("Waiting for server response")
   select {
      case err = <-test.quit:
   }
   if err != nil {
      fmt.Println(err.Message)
   }
   fmt.Println("done.")
}
```
