package embertree

import (
	"github.com/dufourgilles/emberlib/asn1"
	"github.com/dufourgilles/emberlib/errors"
)

type Signal interface {
	Decode(reader *asn1.ASNReader) errors.Error
	Encode(writer *asn1.ASNWriter) errors.Error
}

type Target struct {
	Number int32
}
type Source struct {
	Number int32
}

var (
	TargetApplication = asn1.Application(14)
	SourceApplication = asn1.Application(15)
	TargetContext     = asn1.Context(3)
	SourceContext     = asn1.Context(4)
)

func NewTarget(number int32) *Target {
	t := &Target{Number: number}
	return t
}
func NewSource(number int32) *Source {
	return &Source{Number: number}
}

func encodeSignal(writer *asn1.ASNWriter, number int, ctxt uint8) errors.Error {
	err := writer.StartSequence(ctxt)
	if err != nil {
		return errors.Update(err)
	}
	err = writer.StartSequence(asn1.Context(0))
	if err != nil {
		return errors.Update(err)
	}
	err = writer.WriteInt(number)
	if err != nil {
		return errors.Update(err)
	}

	writer.EndSequence()
	return writer.EndSequence()
}

func (t *Target) Encode(writer *asn1.ASNWriter) errors.Error {
	return encodeSignal(writer, int(t.Number), TargetApplication)
}

func (s *Source) Encode(writer *asn1.ASNWriter) errors.Error {
	return encodeSignal(writer, int(s.Number), SourceApplication)
}

func DecodeSignal(reader *asn1.ASNReader) (Signal, errors.Error) {
	var signal Signal
	ctxt, err := reader.Peek()
	offset := reader.TopOffset()
	if err != nil {
		return nil, errors.Update(err)
	}
	switch ctxt {
	case TargetApplication:
		signal = NewTarget(0)
		break
	case SourceApplication:
		signal = NewSource(0)
		break
	default:
		return nil, errors.New("Unknown signal tag 0x%x at offset %d.", ctxt, offset)
	}
	err = signal.Decode(reader)
	return signal, err
}

func decodeSignal(s interface{}, reader *asn1.ASNReader, ctxt uint8) errors.Error {
	_, signalReader, err := reader.ReadSequenceStart(ctxt)
	if err != nil {
		return errors.Update(err)
	}
	_, ctxtReader, err := signalReader.ReadSequenceStart(asn1.Context(0))
	if err != nil {
		return errors.Update(err)
	}
	number, err := ctxtReader.ReadInt()
	if err != nil {
		return errors.Update(err)
	}
	if ctxt == TargetApplication {
		target := s.(*Target)
		target.Number = int32(number)
	} else {
		source := s.(*Source)
		source.Number = int32(number)
	}
	err = ctxtReader.ReadSequenceEnd()
	if err != nil {
		return errors.Update(err)
	}
	return signalReader.ReadSequenceEnd()
}

func (t *Target) Decode(reader *asn1.ASNReader) errors.Error {
	return decodeSignal(t, reader, TargetApplication)
}

func (s *Source) Decode(reader *asn1.ASNReader) errors.Error {
	return decodeSignal(s, reader, SourceApplication)
}
