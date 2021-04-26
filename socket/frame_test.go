package socket_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dufourgilles/emberlib/socket"
)

func TestCalculateCRC(t *testing.T) {
	buffer := make([]byte, 4)
	for i := 0; i < 4; i++ {
		buffer[i] = uint8(i) + 1
		//fmt.Printf("%d", buffer[i])
	}
	//fmt.Printf("\n")
	crc := socket.CalculateCRC(bytes.NewReader(buffer))
	if crc != 50798 {
		t.Errorf("CRC Failure 0x%x/0xFFFF", crc)
		return
	}

	buffer = []byte{0, 14, 0, 1, 192, 1, 2, 31, 2, 1, 2, 3, 4}
	crc = socket.CalculateCRC(bytes.NewReader(buffer))
	if crc != 10588 {
		t.Errorf("CRC Failure 0x%x/0xFFFF", crc)
		return
	}
}

func TestEncodeMessage(t *testing.T) {
	buffer := make([]byte, 4)
	for i := 0; i < 4; i++ {
		buffer[i] = uint8(i) + 1
		//fmt.Printf("%d", buffer[i])
	}
	frameList := socket.EncodeMessage(buffer)
	if frameList.Size() != 1 {
		t.Errorf("Received invalid number of frames. %d received.  1 expected", frameList.Size())
	}
	expectedResponse := []uint8{254, 0, 14, 0, 1, 192, 1, 2, 31, 2, 1, 2, 3, 4, 163, 214, 255}
	frame, err := frameList.GetAt(0)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	frameBytes := make([]byte, frame.Len())
	textFrame := ""
	frame.Read(frameBytes)
	for i := 0; i < len(frameBytes); i++ {
		textFrame = fmt.Sprintf("%s%d ", textFrame, frameBytes[i])
		if frameBytes[i] != expectedResponse[i] {
			t.Errorf("Encoded Buffer mismatch %s", textFrame)
			return
		}
	}
}

type TestHandler struct {
	data []byte
	err  error
}

func (h *TestHandler) packetHandler(packet []byte) error {
	fmt.Println("Received packet")
	h.data = packet
	return nil
}
func (h *TestHandler) errorHandler(err error) {
	h.err = err
}

func TestDecoder(t *testing.T) {
	var h TestHandler
	h.err = nil
	decoder := socket.NewS101Decoder(
		h.packetHandler,
		h.packetHandler,
		h.packetHandler,
		h.errorHandler)
	frame := []uint8{254, 0, 14, 0, 1, 192, 1, 2, 31, 2, 1, 2, 3, 4, 163, 214, 255}
	decoder.DecodeBuffer(frame)
	if h.err != nil {
		t.Error(h.err)
	}
	if len(h.data) == 0 {
		t.Errorf("Failed to decode frame")
		return
	}
}
