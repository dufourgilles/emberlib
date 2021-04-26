package embertree

import (
	"github.com/dufourgilles/emberlib/asn1"
	"github.com/dufourgilles/emberlib/errors"
)

func EncodeSignals(writer *asn1.ASNWriter, ctxt uint8, signals []Signal) errors.Error {
	err := writer.StartSequence(ctxt)
	if err != nil {
		return errors.Update(err)
	}
	err = writer.StartSequence(asn1.EMBER_SEQUENCE)
	if err != nil {
		return errors.Update(err)
	}

	for _, signal := range signals {
		err = writer.StartSequence(asn1.Context(0))
		if err != nil {
			return errors.Update(err)
		}
		err = signal.Encode(writer)
		if err != nil {
			return errors.Update(err)
		}
		err = writer.EndSequence()
		if err != nil {
			return errors.Update(err)
		}
	}

	writer.EndSequence()
	return writer.EndSequence()
}

func (element *Element) EncodeTargets(writer *asn1.ASNWriter) errors.Error {
	if !element.isMatrix {
		return errors.New("Element not a Matrix.")
	}
	return EncodeSignals(writer, TargetContext, element.targets)
}

func (element *Element) EncodeSources(writer *asn1.ASNWriter) errors.Error {
	if !element.isMatrix {
		return errors.New("Element not a Matrix.")
	}
	return EncodeSignals(writer, SourceContext, element.sources)
}

func DeodeSignals(reader *asn1.ASNReader, ctxt uint8) ([]Signal, errors.Error) {
	var signals []Signal
	_, signalsReader, err := reader.ReadSequenceStart(ctxt)
	if err != nil {
		return nil, errors.Update(err)
	}
	_, setReader, err := signalsReader.ReadSequenceStart(asn1.EMBER_SEQUENCE)
	if err != nil {
		return nil, errors.Update(err)
	}

	for setReader.Len() > 0 {
		_, ctxtReader, err := setReader.ReadSequenceStart(asn1.Context(0))
		if err != nil {
			return nil, errors.Update(err)
		}
		signal, err := DecodeSignal(ctxtReader)
		if err != nil {
			return nil, errors.Update(err)
		}
		signals = append(signals, signal)
		err = ctxtReader.ReadSequenceEnd()
		if err != nil {
			return nil, errors.Update(err)
		}
		end, err := setReader.CheckSequenceEnd()
		if end {
			break
		}
		if err != nil {
			return nil, errors.Update(err)
		}
	}
	return signals, signalsReader.ReadSequenceEnd()
}

func (element *Element) DecodeTargets(reader *asn1.ASNReader) errors.Error {
	if !element.isMatrix {
		return errors.New("Element not a Matrix.")
	}
	signals, err := DeodeSignals(reader, TargetContext)
	if err != nil {
		return errors.Update(err)
	}
	element.targets = signals
	return nil
}

func (element *Element) DecodeSources(reader *asn1.ASNReader) errors.Error {
	if !element.isMatrix {
		return errors.New("Element not a Matrix.")
	}
	signals, err := DeodeSignals(reader, SourceContext)
	if err != nil {
		return errors.Update(err)
	}
	element.sources = signals
	return nil
}

func (element *Element) EncodeConnections(writer *asn1.ASNWriter) errors.Error {
	if !element.isMatrix {
		return errors.New("Element not a Matrix.")
	}
	err := writer.StartSequence(asn1.Context(5))
	if err != nil {
		return errors.Update(err)
	}
	err = writer.StartSequence(asn1.EMBER_SEQUENCE)
	if err != nil {
		return errors.Update(err)
	}

	for _, connection := range element.connections {
		err = writer.StartSequence(asn1.Context(0))
		if err != nil {
			return errors.Update(err)
		}
		err = connection.Encode(writer)
		if err != nil {
			return errors.Update(err)
		}
		err = writer.EndSequence()
		if err != nil {
			return errors.Update(err)
		}
	}

	writer.EndSequence()
	return writer.EndSequence()
}

func (element *Element) DecodeConnections(reader *asn1.ASNReader) errors.Error {
	if !element.isMatrix {
		return errors.New("Element not a Matrix.")
	}
	var connections []*Connection
	_, connectionReader, err := reader.ReadSequenceStart(asn1.Context(5))
	if err != nil {
		return errors.Update(err)
	}
	_, setReader, err := connectionReader.ReadSequenceStart(asn1.EMBER_SEQUENCE)
	if err != nil {
		return errors.Update(err)
	}
	for setReader.Len() > 0 {
		_, ctxtReader, err := setReader.ReadSequenceStart(asn1.Context(0))
		if err != nil {
			return errors.Update(err)
		}
		connection := Connection{}
		err = connection.Decode(ctxtReader)
		if err != nil {
			return errors.Update(err)
		}
		connections = append(connections, &connection)
		err = ctxtReader.ReadSequenceEnd()
		if err != nil {
			return errors.Update(err)
		}
		end, err := setReader.CheckSequenceEnd()
		if end {
			break
		}
		if err != nil {
			return errors.Update(err)
		}
	}
	element.connections = connections
	return connectionReader.ReadSequenceEnd()
}

func (element *Element) SetTargets(targets []Signal) errors.Error {
	if !element.isMatrix {
		return errors.New("Element not a matrix. Can't SetTargets.")
	}
	element.targets = targets
	return nil
}

func (element *Element) GetTargets() ([]Signal, errors.Error) {
	if !element.isMatrix {
		return nil, errors.New("Element not a matrix. Can't GetTargets.")
	}
	return element.targets, nil
}

func (element *Element) SetSources(sources []Signal) errors.Error {
	if !element.isMatrix {
		return errors.New("Element not a matrix. Can't SetSources.")
	}
	element.sources = sources
	return nil
}

func (element *Element) GetSources() ([]Signal, errors.Error) {
	if !element.isMatrix {
		return nil, errors.New("Element not a matrix. Can't GetSources.")
	}
	return element.sources, nil
}

func (element *Element) SetConnections(connections []*Connection) errors.Error {
	if !element.isMatrix {
		return errors.New("Element not a matrix. Can't SetConnections.")
	}
	element.connections = connections
	return nil
}

func (element *Element) GetConnections() ([]*Connection, errors.Error) {
	if !element.isMatrix {
		return nil, errors.New("Element not a matrix. Can't GetConnections.")
	}
	return element.connections, nil
}
