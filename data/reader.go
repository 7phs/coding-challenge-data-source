package data

import (
	"bufio"
	"io"
)

type ValidatedRecord interface {
	Validate() error
}

type RecordFabric func() ValidatedRecord
type RecordFilter func(interface{}) bool

type Source struct {
	reader       *bufio.Reader
	recordFabric RecordFabric
	unMarshaling RecordUnMarshaling
	filters      []RecordFilter

	recordCount int

	fatalErr error
}

func NewSource(reader io.Reader, fabric RecordFabric) *Source {
	o := &Source{
		reader:       bufio.NewReader(reader),
		recordFabric: fabric,
	}

	o.unMarshaling = o.matchUnMarshaling

	return o
}

func (o *Source) matchUnMarshaling(line []byte, record interface{}) error {
	o.unMarshaling, o.fatalErr = MatchUnMarshaling(line, record)

	return o.fatalErr
}

func (o *Source) Filter(filter RecordFilter) *Source {
	o.filters = append(o.filters, filter)

	return o
}

func (o *Source) raiseError(err error) error {
	o.fatalErr = err

	return err
}

func (o *Source) HasError() error {
	return o.fatalErr
}

func (o *Source) readLine() ([]byte, error) {
	var (
		line     []byte
		piece    []byte
		isPrefix = true
		err      error
	)

	for isPrefix {
		piece, isPrefix, err = o.reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return nil, err
			}

			return nil, o.raiseError(err)
		}

		line = append(line, piece...)
	}

	return line, nil
}

func (o *Source) unmarshalRecord(line []byte) (interface{}, error) {
	record := o.recordFabric()
	if err := o.unMarshaling(line, record); err != nil {
		return nil, o.raiseError(err)
	}

	if record.Validate() != nil {
		return nil, nil
	}

	for _, filter := range o.filters {
		if !filter(record) {
			return nil, nil
		}
	}

	return record, nil
}

func (o *Source) Collect() (result RecordList) {
	for {
		line, err := o.readLine()
		if err != nil {
			return
		}

		record, err := o.unmarshalRecord(line)
		if err != nil {
			return
		}

		if record == nil {
			continue
		}

		result.Add(record)
	}

	return
}
