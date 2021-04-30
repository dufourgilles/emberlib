package embertree

import (
	"github.com/dufourgilles/emberlib/errors"

	"github.com/dufourgilles/emberlib/asn1"
)

const COMMAND_SUBSCRIBE = 30
const COMMAND_UNSUBSCRIBE = 31
const COMMAND_GETDIRECTORY = 32
const COMMAND_INVOKE = 33

type FieldFlags int

const (
	SparseFieldFlags      FieldFlags = -2
	ALLFieldFlags         ValueType  = iota
	DefaultFieldFlags     ValueType  = iota
	IdentifierFieldFlags  ValueType  = iota
	DescriptionFieldFlags ValueType  = iota
	TreeFieldFlags        ValueType  = iota
	ValueFieldFlags       ValueType  = iota
	ConnectionsFieldFlags ValueType  = iota
)

type CommandContents struct {
	fieldFlags FieldFlags
}

var CommandApplication = asn1.Application(2)

func NewCommandContents() EmberContents {
	return &CommandContents{fieldFlags: FieldFlags(DefaultFieldFlags)}
}

func NewCommand(number int) *Element {
	command := NewElement(CommandApplication, number, NewCommandContents)
	commandContents := command.CreateContent().(*CommandContents)
	if number == COMMAND_GETDIRECTORY {
		commandContents.fieldFlags = FieldFlags(ALLFieldFlags)
	}
	return command
}

func (cc *CommandContents) SetFieldFlags(fieldFlags FieldFlags) {
	cc.fieldFlags = fieldFlags
}

func (cc *CommandContents) GetFieldFlags() FieldFlags {
	return cc.fieldFlags
}

func (cc *CommandContents) Encode(writer *asn1.ASNWriter) errors.Error {
	return writer.WriteInt(int(cc.fieldFlags))
}

func (cc *CommandContents) Decode(reader *asn1.ASNReader) errors.Error {
	val, err := reader.ReadInt()
	if err != nil {
		return err
	}
	cc.fieldFlags = FieldFlags(val)
	return err
}
