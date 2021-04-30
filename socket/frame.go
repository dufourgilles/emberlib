package socket

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/dufourgilles/emberlib/errors"
)

const S101_BOF = 0xFE
const S101_EOF = 0xFF
const S101_CE = 0xFD
const S101_XOR = 0x20
const S101_INV = 0xF8
const SLOT = 0x00
const MSG_EMBER = 0x0E
const CMD_EMBER = 0x00
const CMD_KEEPALIVE_REQ = 0x01
const CMD_KEEPALIVE_RESP = 0x02
const VERSION = 0x01
const FLAG_SINGLE_PACKET = 0xC0
const FLAG_FIRST_MULTI_PACKET = 0x80
const FLAG_LAST_MULTI_PACKET = 0x40
const FLAG_EMPTY_PACKET = 0x20
const FLAG_MULTI_PACKET = 0x00

const DTD_GLOW = 0x01
const DTD_VERSION_MAJOR = 0x02
const DTD_VERSION_MINOR = 0x1F

var CRC_TABLE = []uint16{ 0x0000, 0x1189, 0x2312, 0x329b, 0x4624, 0x57ad, 0x6536, 0x74bf,
    0x8c48, 0x9dc1, 0xaf5a, 0xbed3, 0xca6c, 0xdbe5, 0xe97e, 0xf8f7,
    0x1081, 0x0108, 0x3393, 0x221a, 0x56a5, 0x472c, 0x75b7, 0x643e,
    0x9cc9, 0x8d40, 0xbfdb, 0xae52, 0xdaed, 0xcb64, 0xf9ff, 0xe876,
    0x2102, 0x308b, 0x0210, 0x1399, 0x6726, 0x76af, 0x4434, 0x55bd,
    0xad4a, 0xbcc3, 0x8e58, 0x9fd1, 0xeb6e, 0xfae7, 0xc87c, 0xd9f5,
    0x3183, 0x200a, 0x1291, 0x0318, 0x77a7, 0x662e, 0x54b5, 0x453c,
    0xbdcb, 0xac42, 0x9ed9, 0x8f50, 0xfbef, 0xea66, 0xd8fd, 0xc974,
    0x4204, 0x538d, 0x6116, 0x709f, 0x0420, 0x15a9, 0x2732, 0x36bb,
    0xce4c, 0xdfc5, 0xed5e, 0xfcd7, 0x8868, 0x99e1, 0xab7a, 0xbaf3,
    0x5285, 0x430c, 0x7197, 0x601e, 0x14a1, 0x0528, 0x37b3, 0x263a,
    0xdecd, 0xcf44, 0xfddf, 0xec56, 0x98e9, 0x8960, 0xbbfb, 0xaa72,
    0x6306, 0x728f, 0x4014, 0x519d, 0x2522, 0x34ab, 0x0630, 0x17b9,
    0xef4e, 0xfec7, 0xcc5c, 0xddd5, 0xa96a, 0xb8e3, 0x8a78, 0x9bf1,
    0x7387, 0x620e, 0x5095, 0x411c, 0x35a3, 0x242a, 0x16b1, 0x0738,
    0xffcf, 0xee46, 0xdcdd, 0xcd54, 0xb9eb, 0xa862, 0x9af9, 0x8b70,
    0x8408, 0x9581, 0xa71a, 0xb693, 0xc22c, 0xd3a5, 0xe13e, 0xf0b7,
    0x0840, 0x19c9, 0x2b52, 0x3adb, 0x4e64, 0x5fed, 0x6d76, 0x7cff,
    0x9489, 0x8500, 0xb79b, 0xa612, 0xd2ad, 0xc324, 0xf1bf, 0xe036,
    0x18c1, 0x0948, 0x3bd3, 0x2a5a, 0x5ee5, 0x4f6c, 0x7df7, 0x6c7e,
    0xa50a, 0xb483, 0x8618, 0x9791, 0xe32e, 0xf2a7, 0xc03c, 0xd1b5,
    0x2942, 0x38cb, 0x0a50, 0x1bd9, 0x6f66, 0x7eef, 0x4c74, 0x5dfd,
    0xb58b, 0xa402, 0x9699, 0x8710, 0xf3af, 0xe226, 0xd0bd, 0xc134,
    0x39c3, 0x284a, 0x1ad1, 0x0b58, 0x7fe7, 0x6e6e, 0x5cf5, 0x4d7c,
    0xc60c, 0xd785, 0xe51e, 0xf497, 0x8028, 0x91a1, 0xa33a, 0xb2b3,
    0x4a44, 0x5bcd, 0x6956, 0x78df, 0x0c60, 0x1de9, 0x2f72, 0x3efb,
    0xd68d, 0xc704, 0xf59f, 0xe416, 0x90a9, 0x8120, 0xb3bb, 0xa232,
    0x5ac5, 0x4b4c, 0x79d7, 0x685e, 0x1ce1, 0x0d68, 0x3ff3, 0x2e7a,
    0xe70e, 0xf687, 0xc41c, 0xd595, 0xa12a, 0xb0a3, 0x8238, 0x93b1,
    0x6b46, 0x7acf, 0x4854, 0x59dd, 0x2d62, 0x3ceb, 0x0e70, 0x1ff9,
    0xf78f, 0xe606, 0xd49d, 0xc514, 0xb1ab, 0xa022, 0x92b9, 0x8330,
    0x7bc7, 0x6a4e, 0x58d5, 0x495c, 0x3de3, 0x2c6a, 0x1ef1, 0x0f78}

type S101CodecEvent string
const (
	KEEP_ALIVE_REQUEST S101CodecEvent = "keepAliveReq"
    EMBER_PACKET = "emberPacket"
    KEEP_ALIVE_RESPONSE = "keepAliveResp"
)

func CalculateCRC(reader *bytes.Reader) uint16 {
	var crc uint16 = 0xFFFF
	for  ; ;  {
		b, e := reader.ReadByte()
		if e != nil { 
			break; 
		}
		crc = ((crc >> 8) ^ CRC_TABLE[(crc ^ uint16(b)) & 0xFF]) & 0xFFFF;
		//fmt.Printf("%d - %d\n", b ,crc)
	}
	return crc;
}

func CalculateCRCCE(reader *bytes.Reader) uint16 {
	var crc uint16 = 0xFFFF;
	for   ; ;  {
		b, e := reader.ReadByte()
		if e != nil { 
			break; 
		}
		if (b == S101_CE) {
			nb,_ := reader.ReadByte()
			b = S101_XOR ^ nb;
		}
		crc = ((crc >> 8) ^ CRC_TABLE[(crc ^ uint16(b)) & 0xFF]) & 0xFFFF;
		//fmt.Printf("%d - %d\n", b ,crc)
	}
	return crc;
}

func GetKeepaliveRequest() *bytes.Buffer {
	packet := &bytes.Buffer{}
	packet.WriteByte(S101_BOF);
	packet.WriteByte(SLOT);
	packet.WriteByte(MSG_EMBER);
	packet.WriteByte(CMD_KEEPALIVE_REQ);
	packet.WriteByte(VERSION);
	return finalizeBuffer(packet);
}

func GetKeepAliveResponse() *bytes.Buffer {
	packet := &bytes.Buffer{}
	packet.WriteByte(S101_BOF);
	packet.WriteByte(SLOT);
	packet.WriteByte(MSG_EMBER);
	packet.WriteByte(CMD_KEEPALIVE_RESP);
	packet.WriteByte(VERSION);
	return finalizeBuffer(packet);
}

func makeBERFrame(flags uint8, data []byte) *bytes.Buffer {
	var frame bytes.Buffer
	frame.WriteByte(S101_BOF);
	frame.WriteByte(SLOT);
	frame.WriteByte(MSG_EMBER);
	frame.WriteByte(CMD_EMBER);
	frame.WriteByte(VERSION);
	frame.WriteByte(flags);
	frame.WriteByte(DTD_GLOW);
	frame.WriteByte(2); // number of app bytes
	frame.WriteByte(DTD_VERSION_MINOR);
	frame.WriteByte(DTD_VERSION_MAJOR);
	frame.Write(data);
	return finalizeBuffer(&frame);
}

func finalizeBuffer(smartbuf *bytes.Buffer) *bytes.Buffer {
	reader := bytes.NewReader(smartbuf.Bytes())
	// skip first byte
	reader.ReadByte()
	var crc = CalculateCRCCE(reader)
	fmt.Println(crc)
	crc = (^crc) &0xFFFF
	fmt.Println(crc)
	var crc_hi uint8 = uint8(crc >> 8);
	var crc_lo uint8 = uint8(crc & 0xFF);
	fmt.Println(crc_hi, crc_lo)
	if (crc_lo < S101_INV) {
		smartbuf.WriteByte(crc_lo);
	} else {
		smartbuf.WriteByte(S101_CE);
		smartbuf.WriteByte(crc_lo ^ S101_XOR);
	}

	if (crc_hi < S101_INV) {
		smartbuf.WriteByte(crc_hi);
	} else {
		smartbuf.WriteByte(S101_CE);
		smartbuf.WriteByte(crc_hi ^ S101_XOR);
	}

	smartbuf.WriteByte(S101_EOF);
	return smartbuf
}

func ValidateFrame(buf *bytes.Reader) bool {
	return CalculateCRC(buf) == 0xF0B8;
}

type S101FrameList struct {
	numFrames int
	increment int
	frames []*bytes.Buffer
}


func NewS101FrameList(increment int) *S101FrameList {
	var list S101FrameList
	if increment == 0 {
		increment = 1
	}
	list.frames = make([]*bytes.Buffer, increment)
	list.numFrames = 0
	list.increment = increment
	return &list
}

func (s *S101FrameList)Size() int {
	return s.numFrames
}

func (s *S101FrameList)GetAt(i int) (*bytes.Buffer, errors.Error) {
	if (i >= s.numFrames || i < 0) {
		return &bytes.Buffer{}, errors.New("out of bound")
	}
	return s.frames[i], nil
}

func (s *S101FrameList)GetBytesAt(i int) ([]byte, errors.Error) {
	if (i >= s.numFrames || i < 0) {
		return nil, errors.New("out of bound")
	}
	res := make([]byte, s.frames[i].Len())
	s.frames[i].Read(res)
	return res, nil
}

func (s *S101FrameList)addFrame(frame *bytes.Buffer) {
	if s.numFrames >= len(s.frames) {
		newFrames := make([]*bytes.Buffer, s.numFrames + s.increment)
		for j:= 0; j < int(s.numFrames); j++ {
			newFrames[j] = s.frames[j]
		}
		s.frames = newFrames
	}
	s.frames[s.numFrames] = frame
	s.numFrames++
}

func EncodeMessage(data [] byte) *S101FrameList {
	var flag uint8
	var frame []byte
	frameList := NewS101FrameList(1)
    var encbuf bytes.Buffer;
    for i := 0; i < len(data); i++ {
        b := data[i];
        if b < S101_INV {
            encbuf.WriteByte(b);
        } else {
            encbuf.WriteByte(S101_CE);
            encbuf.WriteByte(b ^ S101_XOR);
        }

        if(encbuf.Len() >= 1024 && i < len(data)-1) {			
			flag = FLAG_FIRST_MULTI_PACKET
			if (frameList.numFrames >= 1) {
				flag = FLAG_MULTI_PACKET
			}
			frame = make([]byte, encbuf.Len())
			encbuf.Read(frame);
			frameList.addFrame(makeBERFrame(flag, frame));
            encbuf.Reset()
        }
    }
	flag = FLAG_SINGLE_PACKET
	frame = make([]byte, encbuf.Len())
	encbuf.Read(frame);
    if frameList.numFrames == 0 {
        flag = FLAG_SINGLE_PACKET;
    } else {
        flag = FLAG_LAST_MULTI_PACKET
    }
	frameList.addFrame(makeBERFrame(flag, frame));
    return frameList;
};

type EmberFrameHeader struct {
	Version uint8
	Flags uint8
	Dtd uint8
	AppByteLen uint8
	MinorVersion uint8
	MajorVersion uint8
}

type EmberFrame struct {
	Header EmberFrameHeader
	Payload []byte
}

func NewEmberFrameHeader() *EmberFrameHeader {
	return &EmberFrameHeader{Version: VERSION, Dtd: DTD_GLOW, AppByteLen: 2 }
}
func NewEmberFrame() *EmberFrame {
	return &EmberFrame{}
}

func (h *EmberFrameHeader)Read(reader *bytes.Reader) error {
	return binary.Read(reader, binary.LittleEndian, h)
}

type PacketHandler func([]byte) errors.Error
type ErrorHandler func(errors.Error)

type S101Decoder struct {
	escaped bool
	inbuf bytes.Buffer
	emberbuf bytes.Buffer
	keepAliveReqHandler PacketHandler
	keepAliveResHandler PacketHandler
	emberPacketHandler PacketHandler
	errorHandler ErrorHandler
}

func NewS101Decoder(
	keepAliveReqHandler PacketHandler, 
	keepAliveResHandler PacketHandler, 
	emberPacketHandler PacketHandler,
	errorHandler ErrorHandler) *S101Decoder {
	var s101Decoder S101Decoder
	s101Decoder.escaped = false
	s101Decoder.keepAliveReqHandler = keepAliveReqHandler
	s101Decoder.keepAliveResHandler = keepAliveResHandler
	s101Decoder.emberPacketHandler = emberPacketHandler
	s101Decoder.errorHandler = errorHandler
	return &s101Decoder
}

func (decoder *S101Decoder)DecodeBuffer( bufLen int, buf []byte) {
	for i := 0; i < bufLen; i++ {
		b := buf[i];
		if (decoder.escaped) {
			decoder.inbuf.WriteByte(b ^ S101_XOR);
			decoder.escaped = false;
		} else if b == S101_CE {
			decoder.escaped = true;
		} else if b == S101_BOF {
			decoder.inbuf.Reset();
			decoder.escaped = false;
		} else if b == S101_EOF {
			fmt.Printf("End of Frame - frame size %d", decoder.inbuf.Len())
			fmt.Println(decoder.inbuf)
			buffer := make([]byte, decoder.inbuf.Len())
			decoder.inbuf.Read(buffer)
			decoder.inbuf.Reset();
			err := decoder.HandleFrame(buffer);
			if err != nil {
				decoder.errorHandler(err)
			}
		} else {
			decoder.inbuf.WriteByte(b);
		}
	}
}

func (decoder *S101Decoder)HandleFrame(buffer []byte) errors.Error {
	fmt.Printf("Frame parsing. total length %d\n", len(buffer))
	if !ValidateFrame(bytes.NewReader(buffer)) {
		return errors.New("dropping frame with invalid CRC")
	}
	var (
		slot byte
		message byte
		command byte
		err error
	)
	// remove CRC - 2 bytes
	frame := bytes.NewReader(buffer[:len(buffer) - 2])
	slot,err = frame.ReadByte()
	if err != nil {
		return errors.NewError(err)
	}
	message,err = frame.ReadByte()
	if err != nil {
		return errors.NewError(err)
	}
	if (slot != SLOT || message != MSG_EMBER) {
		return errors.New(fmt.Sprintf("dropping frame (not an ember frame; slot=%d, msg=%d)", slot, message))
	}
	command,err = frame.ReadByte()
	if err != nil {
		return errors.NewError(err)
	}
	if command == CMD_KEEPALIVE_REQ {
		return decoder.keepAliveReqHandler(buffer)
	} else if command == CMD_KEEPALIVE_RESP {
		return decoder.keepAliveResHandler(buffer)
	} else if command == CMD_EMBER {
		return decoder.HandleEmberFrame(frame);
	} else {		
		return errors.New(fmt.Sprintf("Unknown command type %d", command))
	}
}

func (decoder *S101Decoder)HandleEmberFrame(frame *bytes.Reader) errors.Error {
	fmt.Printf("Ember Frame parsing. Total size %d\n", frame.Len())
	var emberFrame EmberFrame
	err := emberFrame.Header.Read(frame)
	if err != nil { return errors.NewError(err) }
	
	if emberFrame.Header.Version != VERSION {
		// ok to accept different version
	}

	if emberFrame.Header.Dtd != DTD_GLOW {
		return errors.New(fmt.Sprintf("Dropping frame with non-Glow DTD %d", emberFrame.Header.Dtd))
	}

	if emberFrame.Header.AppByteLen != 2 {
		return errors.New("Frame with unknown DTD length")
	}

	if emberFrame.Header.Flags & FLAG_FIRST_MULTI_PACKET != 0 {
		fmt.Println("Ember Frame first multi packet")
		decoder.emberbuf.Reset()
	}

	if emberFrame.Header.Flags & FLAG_EMPTY_PACKET == 0 {
		fmt.Printf("Ember Frame NOT empty packet. Size %d to %d\n", frame.Len(), decoder.emberbuf.Len())
		frame.WriteTo(&decoder.emberbuf)
	}
	if emberFrame.Header.Flags & FLAG_LAST_MULTI_PACKET != 0 {
		fmt.Printf("Ember Frame last packet. Total message size %d\n", decoder.emberbuf.Len())
		emberPacket := make([]byte, decoder.emberbuf.Len())
		decoder.emberbuf.Read(emberPacket)
		decoder.emberbuf.Reset()
		decoder.emberPacketHandler(emberPacket)
	}
	fmt.Println("Ember Frame parsing complete")
	return nil
}