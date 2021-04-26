package socket

import (
	"container/list"
	"fmt"
	"net"

	"github.com/dufourgilles/emberlib/asn1"
	"github.com/dufourgilles/emberlib/embertree"
	"github.com/dufourgilles/emberlib/errors"
)

type packetQueue struct {
	queue *list.List
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

func (p *packetQueue) add(value []byte) {
	p.queue.PushBack(value)
}

func (p *packetQueue) getNext() ([]byte, error) {
	if p.queue.Len() > 0 {
		if val, ok := p.queue.Front().Value.([]byte); ok {
			return val, nil
		}
		return nil, fmt.Errorf("Queue Error: Queue Datatype is incorrect")
	}
	return nil, fmt.Errorf("Queue Error: Queue is empty")
}

type S101Client struct {
	stats S101SocketStats
	conn  *net.UDPConn
	raddr net.UDPAddr
	outQ  *packetQueue
}

func NewS101Client() *S101Client {
	client := S101Client{}
	client.stats.Reset()
	client.outQ = newPacketQueue()
	return &client
}

func (s *S101Client) Connect(raddr *net.UDPAddr) error {
	var err error
	s.conn, err = net.DialUDP("udp", nil, raddr)
	if err != nil {
		s.raddr = *raddr
	}
	return err
}

func (s *S101Client) IsConnected() bool {
	return s.conn != nil
}

func (s *S101Client) SendBER(data []byte) errors.Error {
	if s.IsConnected() {
		return errors.New("Not connected")
	}
	frames := EncodeMessage(data)
	for i := 0; i < frames.Size(); i++ {
		message, err := frames.GetBytesAt(i)
		if err != nil {
			s.stats.TxErrors++
			return errors.NewError(err)
		}
		var res int
		res, err = s.conn.Write(message)
		if err != nil {
			s.stats.TxErrors++
			return errors.NewError(err)
		}
		s.stats.TxPackets++
		s.stats.TxBytes += uint64(res)
	}
	return nil
}

func (s *S101Client) SendBERNode(node *embertree.RootElement) errors.Error {
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
	return s.SendBER(data)
}

func (s *S101Client) GetDirectory(node *embertree.Element, callback embertree.Listener) errors.Error {
	msg, err := node.GetDirectoryMsg(callback)
	if err != nil {
		return errors.Update(err)
	}
	return s.SendBERNode(msg)
}

func (s *S101Client) receiver() {

}
