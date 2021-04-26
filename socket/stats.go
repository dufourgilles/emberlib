package socket

type S101SocketStats struct {
	TxPackets uint64
	RxPackets uint64
	TxBytes	uint64
	RxBytes uint64
	TxErrors uint64	
}

func (s *S101SocketStats)Reset() {
	s.TxBytes = 0
	s.TxErrors = 0
	s.TxPackets = 0
	s.RxBytes = 0
	s.RxPackets = 0
}