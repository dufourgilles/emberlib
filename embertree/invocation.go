package embertree

import (
	"github.com/dufourgilles/emberlib/errors"

	"github.com/dufourgilles/emberlib/asn1"
)

var InvocationApplication = asn1.Application(22)

type Invocation struct {
	invocationID int
	arguments    []*ContentParameter
}

func NewInvocation(id int, arguments []*ContentParameter) *Invocation {
	return &Invocation{invocationID: id, arguments: arguments}
}

func (i *Invocation) Encode(writer *asn1.ASNWriter) errors.Error {
	err := writer.StartSequence(InvocationApplication)
	if err != nil {
		return errors.Update(err)
	}

	err = writer.StartSequence(asn1.Context(0))
	if err != nil {
		return errors.Update(err)
	}
	err = writer.WriteInt(i.invocationID)
	if err != nil {
		return errors.Update(err)
	}
	err = writer.EndSequence()
	if err != nil {
		return errors.Update(err)
	}

	if len(i.arguments) > 0 {
		err = writer.StartSequence(asn1.Context(1))
		if err != nil {
			return errors.Update(err)
		}
		err = writer.StartSequence(asn1.EMBER_SEQUENCE)
		if err != nil {
			return errors.Update(err)
		}
		for _, tuple := range i.arguments {
			err = tuple.Encode(0, writer)
			if err != nil {
				return errors.Update(err)
			}
		}
		err = writer.EndSequence()
		if err != nil {
			return errors.Update(err)
		}
		err = writer.EndSequence()
		if err != nil {
			return errors.Update(err)
		}
	}
	return writer.EndSequence()
}

func (i *Invocation) Decode(reader *asn1.ASNReader) errors.Error {
	_, invocationReader, err := reader.ReadSequenceStart(InvocationApplication)
	if err != nil {
		return errors.Update(err)
	}

	for invocationReader.Len() > 0 {
		peek, err := invocationReader.Peek()
		if err != nil {
			return errors.Update(err)
		}
		_, ctxtReader, err := invocationReader.ReadSequenceStart(peek)
		if err != nil {
			return errors.Update(err)
		}
		switch peek {
		case asn1.Context(0):
			id, err := ctxtReader.ReadInt()
			if err != nil {
				return errors.Update(err)
			}
			i.invocationID = id
			break
		case asn1.Context(1):
			arguments := []*ContentParameter{}
			_, seqReader, err := ctxtReader.ReadSequenceStart(asn1.EMBER_SEQUENCE)
			if err != nil {
				return errors.Update(err)
			}
			for seqReader.Len() > 0 {
				value, err := DecodeValue(seqReader, 0)
				if err != nil {
					return errors.Update(err)
				}
				arguments = append(arguments, value)
				end, err := seqReader.CheckSequenceEnd()
				if end {
					break
				}
				if err != nil {
					return errors.Update(err)
				}
			}
			break
		default:
			return errors.New("Invovation decode error: Unknwon tag %d", peek)
		}
		err = ctxtReader.ReadSequenceEnd()
		if err != nil {
			return errors.Update(err)
		}
		end, err := invocationReader.CheckSequenceEnd()
		if end || err != nil {
			break
		}
	}
	return errors.Update(err)
}
