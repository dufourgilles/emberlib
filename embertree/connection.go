package embertree

import (
	"github.com/dufourgilles/emberlib/asn1"
	"github.com/dufourgilles/emberlib/errors"
)

type ConnectionOperation uint8

const (
	Absolute   ConnectionOperation = 0
	Connect    ConnectionOperation = iota
	Disconnect ConnectionOperation = iota
)

type ConnectionDisposition uint8

const (
	Tally    ConnectionDisposition = 0
	Modified ConnectionDisposition = iota
	Pending  ConnectionOperation   = iota
	Locked   ConnectionOperation   = iota
)

type Connection struct {
	Target      int32
	Sources     []int32
	operation   ConnectionOperation
	disposition ConnectionDisposition
}

const ConnectionApplication = 16

func (c *Connection) SetDisposition(d int) errors.Error {
	if d < int(Tally) || d > int(Locked) {
		return errors.New("Invalid disposition %d.", d)
	}
	c.disposition = ConnectionDisposition(d)
	return nil
}

func (c *Connection) GetDisposition() ConnectionDisposition {
	return c.disposition
}

func (c *Connection) SetOperation(d int) errors.Error {
	if d < int(Absolute) || d > int(Disconnect) {
		return errors.New("Invalid operation %d.", d)
	}
	c.operation = ConnectionOperation(d)
	return nil
}
func (c *Connection) GetOperation() ConnectionOperation {
	return c.operation
}

func (c *Connection) Encode(writer *asn1.ASNWriter) errors.Error {
	err := writer.StartSequence(asn1.Application(ConnectionApplication))
	if err != nil {
		return errors.Update(err)
	}

	// Target
	err = writer.StartSequence(asn1.Context(0))
	if err != nil {
		return errors.Update(err)
	}
	err = writer.WriteInt(int(c.Target))
	if err != nil {
		return errors.Update(err)
	}
	err = writer.EndSequence()
	if err != nil {
		return errors.Update(err)
	}

	//Sources
	if len(c.Sources) > 0 {
		err = writer.StartSequence(asn1.Context(1))
		if err != nil {
			return errors.Update(err)
		}
		err = writer.WriteRelativeOID(c.Sources)
		if err != nil {
			return errors.Update(err)
		}
		err = writer.EndSequence()
		if err != nil {
			return errors.Update(err)
		}
	}

	//Operation
	err = writer.StartSequence(asn1.Context(2))
	if err != nil {
		return errors.Update(err)
	}
	err = writer.WriteInt(int(c.operation))
	if err != nil {
		return errors.Update(err)
	}
	err = writer.EndSequence()
	if err != nil {
		return errors.Update(err)
	}

	//disposition
	err = writer.StartSequence(asn1.Context(3))
	if err != nil {
		return errors.Update(err)
	}
	err = writer.WriteInt(int(c.disposition))
	if err != nil {
		return errors.Update(err)
	}
	err = writer.EndSequence()
	if err != nil {
		return errors.Update(err)
	}

	return writer.EndSequence()
}

func (c *Connection) Decode(reader *asn1.ASNReader) errors.Error {
	_, connectionReader, err := reader.ReadSequenceStart(asn1.Application(ConnectionApplication))
	if err != nil {
		return errors.Update(err)
	}

	// Target
	_, targetReader, err := connectionReader.ReadSequenceStart(asn1.Context(0))
	if err != nil {
		return errors.Update(err)
	}
	target, err := targetReader.ReadInt()
	if err != nil {
		return errors.Update(err)
	}
	c.Target = int32(target)
	err = targetReader.ReadSequenceEnd()
	if err != nil {
		return errors.Update(err)
	}

	for connectionReader.Len() > 0 {
		tag, err := connectionReader.Peek()
		if err != nil {
			return errors.Update(err)
		}

		_, ctxtReader, err := connectionReader.ReadSequenceStart(tag)
		switch tag {
		case asn1.Context(1):
			//Sources
			oid, err := ctxtReader.ReadOID(asn1.EMBER_RELATIVE_OID)
			if err != nil {
				return errors.Update(err)
			}
			c.Sources = oid
			break
		case asn1.Context(2):
			operation, err := ctxtReader.ReadInt()
			if err != nil {
				return errors.Update(err)
			}
			err = c.SetOperation(operation)
			if err != nil {
				return errors.Update(err)
			}
			break
		case asn1.Context(3):
			disposition, err := ctxtReader.ReadInt()
			if err != nil {
				return errors.Update(err)
			}
			err = c.SetDisposition(disposition)
			if err != nil {
				return errors.Update(err)
			}
			break
		default:
			return errors.New("Unknown tag %d", tag)
		}
		err = ctxtReader.ReadSequenceEnd()
		if err != nil {
			return errors.Update(err)
		}
		end, err := connectionReader.CheckSequenceEnd()
		if end {
			break
		}
	}
	return nil
}
