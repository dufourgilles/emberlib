package embertree

import (
	"github.com/dufourgilles/emberlib/errors"

	"github.com/dufourgilles/emberlib/asn1"
)

type NodeContents struct {
	identifier        ContentParameter
	description       ContentParameter
	isRoot            ContentParameter
	isOnline          ContentParameter
	schemaIdentifiers ContentParameter
	templateReference asn1.RelativeOID
}

var NodeApplication = asn1.Application(3)

func NewNodeContents() EmberContents {
	return &NodeContents{}
}

func NewNode(number int) *Element {
	return NewElement(NodeApplication, number, NewNodeContents)
}

func (contents *NodeContents) SetIdentifier(identifer string) {
	contents.identifier.SetString(identifer)
}

func (contents *NodeContents) GetIdentifier() (string, errors.Error) {
	return contents.identifier.GetString()
}

func (contents *NodeContents) SetDescription(description string) {
	contents.description.SetString(description)
}

func (contents *NodeContents) GetDescription() (string, errors.Error) {
	return contents.description.GetString()
}

func (contents *NodeContents) SetIsRoot(isRoot bool) {
	contents.isRoot.SetBool(isRoot)
}

func (contents *NodeContents) GetIsRoot() (bool, errors.Error) {
	return contents.isRoot.GetBool()
}

func (contents *NodeContents) SetIsOnline(isOnline bool) {
	contents.isOnline.SetBool(isOnline)
}

func (contents *NodeContents) GetIsOnline() (bool, errors.Error) {
	return contents.isOnline.GetBool()
}

func (contents *NodeContents) SetSchemaIdentifiers(schemaIdentifiers string) {
	contents.schemaIdentifiers.SetString(schemaIdentifiers)
}

func (contents *NodeContents) GetSchemaIdentifiers() (string, errors.Error) {
	return contents.schemaIdentifiers.GetString()
}

func (contents *NodeContents) SetTemplateReference(templateReference asn1.RelativeOID) {
	contents.templateReference = templateReference
}

func (contents *NodeContents) GetTemplateReference() (asn1.RelativeOID, errors.Error) {
	return contents.templateReference, nil
}

func (nc *NodeContents) Encode(writer *asn1.ASNWriter) errors.Error {
	var err errors.Error
	err = writer.StartSequence(asn1.EMBER_SET)
	if err != nil {
		return errors.Update(err)
	}
	if nc.identifier.IsSet() {
		err = nc.identifier.Encode(0, writer)
		if err != nil {
			return errors.Update(err)
		}
	}

	if nc.description.IsSet() {
		nc.description.Encode(1, writer)
		if err != nil {
			return errors.Update(err)
		}
	}

	if nc.isRoot.IsSet() {
		err = nc.isRoot.Encode(2, writer)
		if err != nil {
			return errors.Update(err)
		}
	}

	if nc.isOnline.IsSet() {
		err = nc.isOnline.Encode(3, writer)
		if err != nil {
			return errors.Update(err)
		}
	}

	if nc.schemaIdentifiers.IsSet() {
		err = nc.schemaIdentifiers.Encode(4, writer)
		if err != nil {
			return errors.Update(err)
		}
	}

	if nc.templateReference != nil {
		err = writer.StartSequence(asn1.Context(5))
		if err != nil {
			return errors.Update(err)
		}
		err = writer.WriteRelativeOID(nc.templateReference)
		if err != nil {
			return errors.Update(err)
		}
		err = writer.EndSequence()
		if err != nil {
			return errors.Update(err)
		}
	}

	return writer.EndSequence() // EMBER_SET
}

func (nc *NodeContents) Decode(reader *asn1.ASNReader) errors.Error {
	var err errors.Error
	var value *ContentParameter
	_, reader, err = reader.ReadSequenceStart(asn1.EMBER_SET)
	if err != nil {
		return errors.Update(err)
	}

	for reader.Len() > 0 {
		peek, err := reader.Peek()
		if err != nil {
			return errors.Update(err)
		}
		if peek != asn1.Context(5) {
			value, err = DecodeValue(reader, peek)
			if err != nil {
				return errors.Update(err)
			}
		}
		switch peek {
		case asn1.Context(0):
			{
				nc.identifier.Set(value)
				break
			}
		case asn1.Context(1):
			{
				nc.description.Set(value)
				break
			}
		case asn1.Context(2):
			{
				nc.isRoot.Set(value)
				break
			}
		case asn1.Context(3):
			{
				nc.isOnline.Set(value)
				break
			}
		case asn1.Context(4):
			{
				nc.schemaIdentifiers.Set(value)
				break
			}
		case asn1.Context(5):
			{
				_, templateReader, err := reader.ReadSequenceStart(peek)
				templateReference, err := templateReader.ReadOID(asn1.EMBER_RELATIVE_OID)
				if err != nil {
					return errors.Update(err)
				}
				nc.templateReference = templateReference
				err = templateReader.ReadSequenceEnd()
				if err != nil {
					return errors.Update(err)
				}
			}
		}
		end, err := reader.CheckSequenceEnd()
		if end {
			break
		}
		if err != nil {
			return errors.Update(err)
		}
	}
	return nil
}
