package socket_test

import (
	"fmt"
	"testing"
	
	"github.com/dufourgilles/emberlib/embertree"
	"github.com/dufourgilles/emberlib/errors"
	"github.com/dufourgilles/emberlib/socket"
	"github.com/dufourgilles/emberlib/logger"
)

type testClient struct {
	quit chan errors.Error
}

func (c *testClient)Receive(node interface{}, err errors.Error) {
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

func TestClient(t *testing.T) {
	test := &testClient{quit: make(chan errors.Error, 1)}
	client := socket.NewS101Client()
	client.SetLogger(logger.NewConsoleLogger(logger.DebugLevel))
	client.SetTimeout(2500)
	fmt.Println("Connecting")
	err := client.Connect("192.168.1.2", 9000)
	if err != nil {
		t.Error(err.Message)
	}
	fmt.Println("Connected. Get Directory")
	client.GetTree(test)
	fmt.Println("Waiting for server response")
	select {
	case err = <-test.quit:
	}
	if err != nil {
		t.Error(err.Message)
	}
	t.Errorf("done.")
}