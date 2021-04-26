package embertree

import (
	"github.com/dufourgilles/emberlib/asn1"
	"github.com/dufourgilles/emberlib/errors"
)

type Label struct {
	BasePath    asn1.RelativeOID
	Description string
}

var LabelApplication = asn1.Application(18)

func NewLabel(basePath asn1.RelativeOID, description string) *Label {
	return &Label{BasePath: basePath, Description: description}
}

func DecodeLabel(reader *asn1.ASNReader) (*Label, errors.Error) {
	l := &Label{}
	err := l.Decode(reader)
	return l, err
}

func (l *Label) Encode(writer *asn1.ASNWriter) errors.Error {
	err := writer.StartSequence(LabelApplication)
	if err != nil {
		return errors.Update(err)
	}
	err = writer.StartSequence(asn1.Context(0))
	if err != nil {
		return errors.Update(err)
	}
	err = writer.WriteRelativeOID(l.BasePath)
	if err != nil {
		return errors.Update(err)
	}
	err = writer.EndSequence()
	if err != nil {
		return errors.Update(err)
	}
	err = writer.StartSequence(asn1.Context(1))
	if err != nil {
		return errors.Update(err)
	}
	err = writer.WriteString(l.Description)
	if err != nil {
		return errors.Update(err)
	}
	err = writer.EndSequence()
	if err != nil {
		return errors.Update(err)
	}
	return writer.EndSequence()
}

func (l *Label) Decode(reader *asn1.ASNReader) errors.Error {
	_, labelReader, err := reader.ReadSequenceStart(LabelApplication)
	if err != nil {
		return errors.Update(err)
	}
	_, ctxtReader, err := labelReader.ReadSequenceStart(asn1.Context(0))
	if err != nil {
		return errors.Update(err)
	}
	oid, err := ctxtReader.ReadOID(asn1.EMBER_RELATIVE_OID)
	if err != nil {
		return errors.Update(err)
	}
	l.BasePath = oid
	err = ctxtReader.ReadSequenceEnd()
	if err != nil {
		return errors.Update(err)
	}
	_, ctxtReader, err = labelReader.ReadSequenceStart(asn1.Context(1))
	if err != nil {
		return errors.Update(err)
	}
	description, err := ctxtReader.ReadString()
	if err != nil {
		return errors.Update(err)
	}
	l.Description = description
	err = ctxtReader.ReadSequenceEnd()
	if err != nil {
		return errors.Update(err)
	}
	return labelReader.ReadSequenceEnd()
}
