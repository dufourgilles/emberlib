package socket

import (
	"container/list"
	"fmt"
	"net"
	"time"

	"github.com/dufourgilles/emberlib/asn1"
	"github.com/dufourgilles/emberlib/embertree"
	"github.com/dufourgilles/emberlib/errors"
)

const maxQueueSize = 256
const maxBufferSize = 65536

type packetQueue struct {
	queue *list.List
}

type channelMessage interface {
	getError() error
	getData() []byte
}

func newPacketQueue() *packetQueue {
	var q packetQueue
	q.queue = list.New()
	return &q
}

func (p *packetQueue) size() int {
	return p.queue.Len()
}

func (p *packetQueue) isEmpty() bool {
	return p.queue.Len() == 0
}

func (p *packetQueue) add(msg *embertree.RootElement) errors.Error{
	if p.size() > maxQueueSize {
		return errors.New("Queue size limit. Drop message.")
	}
	p.queue.PushBack(msg)
	return nil
}

func (p *packetQueue) getNext() (*embertree.RootElement, errors.Error) {
	if p.queue.Len() > 0 {
		qElement :=  p.queue.Front()
		p.queue.Remove(qElement)
		if val, ok := qElement.Value.(*embertree.RootElement); ok {			
			return val, nil
		}
		return nil, errors.New("Queue Error: Queue Datatype is incorrect")
	}
	return nil, nil
}

type _TimerCallbacks struct {
	timer *time.Timer
	callback embertree.Listener
}

type listeningNode interface {
	AddListener(listener embertree.Listener)
}

type S101Client struct {
	stats S101SocketStats
	conn  net.Conn
	raddr string
	outQ  *packetQueue
	decoder *S101Decoder
	msTimeout int
	timerCallbacks map[listeningNode]*_TimerCallbacks
	tree *embertree.RootElement
	incomingDataChan chan channelMessage
	outgoingDataChan chan channelMessage
}

func (s *S101Client)keepAliveReqHandler(kal []byte) errors.Error {
	// This should go at the head of the queue
	s.outQ.queue.PushFront(GetKeepAliveResponse())
	return nil
}

func (s *S101Client)keepAliveResponseHandler(kal []byte) errors.Error {
	// do nothing for the moment
	return nil
}

func (s *S101Client)emberPacketHandler(packet []byte) errors.Error {
	// This should be a valid EmberRoot.
	fmt.Println("Ember Frame - start decoding")
	err := s.tree.Decode(asn1.NewASNReader(packet))
	s.errorHandler(err)
	return err
}

func (s *S101Client)errorHandler(err errors.Error) {
	fmt.Println(err)
}

func NewS101Client() *S101Client {
	client := S101Client{msTimeout: 0}
	client.decoder = NewS101Decoder(client.keepAliveReqHandler,client.keepAliveResponseHandler,client.emberPacketHandler, client.errorHandler )
	client.stats.Reset()
	client.timerCallbacks = make(map[listeningNode]*_TimerCallbacks)
	client.outQ = newPacketQueue()
	client.tree = embertree.NewTree()
	client.incomingDataChan = make(chan channelMessage, 100)
	client.outgoingDataChan = make(chan channelMessage, 1)	
	return &client
}

func (s *S101Client)SetTimeout(msTimeout int) {
	if msTimeout >= 0 {
		s.msTimeout = msTimeout
	}
}

func (s *S101Client)Connect(address string, port uint16) errors.Error {
	var err error
	if s.IsConnected() {
		return errors.New("Client already connected to %s:%d", address, port)
	}
	s.raddr = fmt.Sprintf("%s:%d", address, port)
	s.conn, err = net.Dial("tcp", s.raddr)
	if err != nil {
		return errors.NewError(err)
	}
	go s.iomanager()
	return nil
}

func (s *S101Client)Disconnect() errors.Error {
	if !s.IsConnected() {
		return errors.New("Client not connected.")
	}
	err := s.conn.Close()
	s.conn = nil
	return errors.NewError(err)
}

func (s *S101Client)IsConnected() bool {
	return s.conn != nil
}

func (s *S101Client)sendBER(data []byte) errors.Error {
	if !s.IsConnected() {
		return errors.New("Not connected")
	}
	frames := EncodeMessage(data)
	for i := 0; i < frames.Size(); i++ {
		message, err := frames.GetBytesAt(i)
		if err != nil {
			s.stats.TxErrors++
			return errors.Update(err)
		}
		var res int
		res, e := s.conn.Write(message)
		if e != nil {
			s.stats.TxErrors++
			return errors.NewError(e)
		}
		s.stats.TxPackets++
		s.stats.TxBytes += uint64(res)
	}
	return nil
}

func (s *S101Client)sendBERNode(node *embertree.RootElement) errors.Error {
	if node == nil {
		return errors.New("null node")
	}
	writer := asn1.ASNWriter{}
	err := node.Encode(&writer)
	if err != nil {
		return err
	}
	data := make([]byte, writer.Len())
	writer.Read(data)
	return s.sendBER(data)
}

func (s *S101Client)freeTimerCallback(node interface{}, err errors.Error) {
	treeNode := node.(listeningNode)
	timerCB := s.timerCallbacks[treeNode]
	if timerCB == nil {
		return
	}
	timerCB.timer.Stop()
	delete(s.timerCallbacks, treeNode)
}

func (s *S101Client)runTimer(node *embertree.Element, callback embertree.Listener, timeoutError errors.Error) {
	var treeNode listeningNode
	if node == nil {
		treeNode = s.tree
	} else {
		treeNode = node
	}
	treeNode.AddListener(s.freeTimerCallback)
	s.timerCallbacks[treeNode] = &_TimerCallbacks{
		callback: callback,
		timer: time.NewTimer(time.Duration(int64(s.msTimeout)) * time.Millisecond),
	}
	fmt.Println("Timer running")
	<-s.timerCallbacks[treeNode].timer.C
	callback(nil, timeoutError)
	delete(s.timerCallbacks, treeNode)
}

func (s *S101Client)GetDirectory(node *embertree.Element, callback embertree.Listener) errors.Error {
	var msg *embertree.RootElement
	var err errors.Error
	timeoutError := errors.New("GetDirectory timed out.")	
	if node == nil {
		msg,err = s.tree.GetDirectoryMsg(callback)
	} else {
		msg, err = node.GetDirectoryMsg(callback)
	}
	if err != nil {
		return errors.Update(err)
	}
	if s.msTimeout > 0 {
		go s.runTimer(node, callback, timeoutError)
	}
	fmt.Println("Add message to Q")
	return s.outQ.add(msg)
}

func (s *S101Client)processBuffer(l int, buffer []byte) {
	//Decode the message
	s.decoder.DecodeBuffer(l, buffer)
}

func (s *S101Client)iomanager() {
	fmt.Println("iomanager started")
	buffer := make([]byte, maxBufferSize)
	for s.IsConnected() {
		// Send messages present in our Q - but no more than max
		messagesToSend := s.outQ.queue.Len()
		if messagesToSend > 3 {
			messagesToSend = 3;
		}
		if messagesToSend > 1 {
			fmt.Printf("iomanager to send %d messages\n", messagesToSend)
		}
		count := 0
		for ; count < messagesToSend; count++ {
			root,_ := s.outQ.getNext()
			if root != nil {
				err := s.sendBERNode(root)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("iomanager invalid nil root to send")
			}
		}

		// collect inbound messages during next 100ms
		//fmt.Println("iomanager read messages")
		deadline := time.Now().Add(100)
		s.conn.SetReadDeadline(deadline)
		for {
			len,_ := s.conn.Read(buffer)
			if len > 0 {
				fmt.Printf("iomanager received a message of %d bytes\n", len)
				s.stats.RxBytes += uint64(len)
				s.stats.RxPackets++
				//process the buffer
				s.processBuffer(len, buffer)
			} else {
				break
			}
		}
	}
	fmt.Println("iomanager stopped")
}
