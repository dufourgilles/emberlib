package socket

import (
	"container/list"
	"fmt"
	"net"
	"time"
	. "github.com/dufourgilles/emberlib/logger"
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
	listener embertree.Listener
	client *S101Client
}

type S101Client struct {
	stats S101SocketStats
	conn  net.Conn
	raddr string
	outQ  *packetQueue
	decoder *S101Decoder
	msTimeout int
	logger Logger
	timerCallbacks map[embertree.ListeningNode]*_TimerCallbacks
	getTreeCallback embertree.Listener
	tree *embertree.RootElement
	incomingDataChan chan channelMessage
	outgoingDataChan chan channelMessage
	freeTimerCallback embertree.Listener
}

func (s *S101Client)keepAliveReqHandler(kal []byte) errors.Error {
	// This should go at the head of the queue
	s.logger.Debug("KAL Request Received.\n")
	s.outQ.queue.PushFront(GetKeepAliveResponse())
	return nil
}

func (s *S101Client)keepAliveResponseHandler(kal []byte) errors.Error {
	// do nothing for the moment
	s.logger.Debug("KAL Response Received.\n")
	return nil
}

func (s *S101Client)emberPacketHandler(packet []byte) errors.Error {
	// This should be a valid EmberRoot.
	s.logger.Debug("Ember Frame - start decoding.\n")
	s.logger.Debugln(packet)
	err := s.tree.Decode(asn1.NewASNReader(packet))
	s.errorHandler(err)
	return err
}

func (s *S101Client)errorHandler(err errors.Error) {
	s.logger.Error(err)
}

func NewS101Client() *S101Client {
	client := S101Client{msTimeout: 0}
	client.decoder = NewS101Decoder(client.keepAliveReqHandler,client.keepAliveResponseHandler,client.emberPacketHandler, client.errorHandler )
	client.stats.Reset()
	client.logger = NewNullLogger()
	client.timerCallbacks = make(map[embertree.ListeningNode]*_TimerCallbacks)
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

func (s *S101Client)SetLogger(logger Logger) {
	if logger != nil {
		s.logger = logger
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


func (t *_TimerCallbacks)Receive(node interface{}, err errors.Error) {
	treeNode := node.(embertree.ListeningNode)
	timerCB := t.client.timerCallbacks[treeNode]
	if timerCB == nil {
		return
	}
	t.client.logger.Debug("Timer Stopped.\n")
	timerCB.timer.Stop()
	delete(t.client.timerCallbacks, treeNode)
	// Remove myself from node
	treeNode.RemoveListener(t)
}

func (s *S101Client)runTimer(node *embertree.Element, callback embertree.Listener, timeoutError errors.Error) {
	var treeNode embertree.ListeningNode
	if node == nil {
		treeNode = s.tree
	} else {
		treeNode = node
	}
	s.logger.Debug("Timer running.\n")
	freeTimerCallback := &_TimerCallbacks{client: s, timer: time.NewTimer(time.Duration(int64(s.msTimeout)) * time.Millisecond)}
	treeNode.AddListener(freeTimerCallback)
	s.timerCallbacks[treeNode] = freeTimerCallback
	<-s.timerCallbacks[treeNode].timer.C
	s.logger.Debug("Timer Kicked.\n")
	callback.Receive(nil, timeoutError)
	s.logger.Debug("Cleaning Timer callback.\n")
	delete(s.timerCallbacks, treeNode)
	treeNode.RemoveListener(s.freeTimerCallback)
}

func (s *S101Client)GetDirectory(node *embertree.Element, callback embertree.Listener) errors.Error {
	var msg *embertree.RootElement
	var err errors.Error
	timeoutError := errors.New("GetDirectory timed out.")	
	if node == nil {
		s.logger.Debug("Send GetDirectory for root.\n")
		msg,err = s.tree.GetDirectoryMsg(callback)		
	} else {
		s.logger.Debug("Send GetDirectory for %s.\n", embertree.Path2String(node.GetPath()))
		msg = node.GetQualifiedDirectoryMsg(callback)
	}
	if err != nil {
		s.logger.Debug("GetDirectory.\n",err)
		return errors.Update(err)
	}
	if s.msTimeout > 0 {
		go s.runTimer(node, callback, timeoutError)
	}
	
	return s.outQ.add(msg)
}

type _ElementListeners struct {
	client *S101Client
	rootListener *_RoottListeners
}

// getElementCallback
func (l *_ElementListeners)Receive(node interface{}, err errors.Error) {
	l.client.logger.Debug("Element Callback.\n")	
	if err == nil && node != nil {
		element := node.(*embertree.Element)
		element.RemoveListener(l)
		for _,child := range(element.Children) {
			l.rootListener.IncPendingGetDir()
			elementCallback := &_ElementListeners{client: l.client, rootListener: l.rootListener}
			go l.client.GetDirectory(child, elementCallback)
		}
	}
	l.rootListener.DecPendingGetDir(node,err)
}

type _RoottListeners struct {
	client *S101Client
	pendingGetDirectory uint
	listener embertree.Listener
}

func (r *_RoottListeners)IncPendingGetDir() {
	r.pendingGetDirectory++
}

func (r *_RoottListeners)DecPendingGetDir(node interface{}, err errors.Error) {
	r.pendingGetDirectory--
	if r.pendingGetDirectory == 0 {
		r.listener.Receive(r.client.tree, err)
	}
}

func (r *_RoottListeners)Receive(node interface{}, err errors.Error) {	
	if err == nil && node != nil {
		root := node.(*embertree.RootElement)
		r.client.logger.Debug("Root CallBack.\n")
		r.client.logger.Debug(root.ToString())
		root.RemoveListener(r)
		for _,element := range(root.RootElementCollection) {
			r.IncPendingGetDir()
			elementCallback := &_ElementListeners{client: r.client, rootListener: r}
			r.client.logger.Debug("Root Callback GetDir for %s.\n", embertree.Path2String(element.GetPath()))
			go r.client.GetDirectory(element, elementCallback)
		}
	}
	if r.pendingGetDirectory == 0 {
		r.listener.Receive(r.client.tree, err)
		return
	}
}

func (s *S101Client)GetTree(listener embertree.Listener) errors.Error {
	rootCallback := &_RoottListeners{client: s, listener: listener, pendingGetDirectory: 0}
	return s.GetDirectory(nil, rootCallback)
}

func (s *S101Client)processBuffer(l int, buffer []byte) {
	//Decode the message
	s.decoder.DecodeBuffer(l, buffer)
}

func (s *S101Client)iomanager() {
	s.logger.Debug("iomanager started.\n")
	buffer := make([]byte, maxBufferSize)
	for s.IsConnected() {
		// Send messages present in our Q - but no more than max
		messagesToSend := s.outQ.queue.Len()
		if messagesToSend > 3 {
			messagesToSend = 3;
		}
		if messagesToSend > 1 {
			s.logger.Debug("iomanager to send %d messages.\n", messagesToSend)
		}
		count := 0
		for ; count < messagesToSend; count++ {
			root,_ := s.outQ.getNext()
			if root != nil {
				err := s.sendBERNode(root)
				if err != nil {
					s.logger.Error(err)
				}
			} else {
				s.logger.Warn("iomanager invalid nil root to send.\n")
			}
		}

		// collect inbound messages during next 100ms
		deadline := time.Now().Add(100)
		s.conn.SetReadDeadline(deadline)
		for {
			len,_ := s.conn.Read(buffer)
			if len > 0 {
				s.logger.Debug("iomanager received a message of %d bytes.\n", len)
				s.stats.RxBytes += uint64(len)
				s.stats.RxPackets++
				//process the buffer
				s.processBuffer(len, buffer)
			} else {
				break
			}
		}
	}
	s.logger.Debug("iomanager stopped.\n")
}
