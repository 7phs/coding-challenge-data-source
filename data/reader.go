package data

import (
	"bufio"
	"io"
	"sync"
)

type ValidatedRecord interface {
	Validate() error
}

type RecordFabric func() ValidatedRecord
type RecordFilter func(interface{}) bool
type ErrorListener func(error)

type Source struct {
	reader            *bufio.Reader
	recordFabric      RecordFabric
	unMarshaling      RecordUnMarshaling
	filters           []RecordFilter
	errorListen       []ErrorListener
	waitErrorListener sync.WaitGroup
}

func NewSource(reader io.Reader, fabric RecordFabric) *Source {
	o := &Source{
		reader:       bufio.NewReader(reader),
		recordFabric: fabric,
	}

	o.unMarshaling = o.matchUnMarshaling

	return o
}

func (o *Source) matchUnMarshaling(line []byte, record interface{}) (err error) {
	o.unMarshaling, err = MatchUnMarshaling(line, record)

	return err
}

func (o *Source) Filter(filter RecordFilter) *Source {
	o.filters = append(o.filters, filter)

	return o
}

func (o *Source) raiseError(err error) error {
	o.waitErrorListener.Add(1)
	go func() {
		defer o.waitErrorListener.Done()

		for _, errCallback := range o.errorListen {
			errCallback(err)
		}
	}()

	return err
}

func (o *Source) Catch(errCallback ErrorListener) *Source {
	o.errorListen = append(o.errorListen, errCallback)

	return o
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

	if err := record.Validate(); err != nil {
		return nil, err
	}

	return record, nil
}

func (o *Source) parseFirstLine() (interface{}, error) {
	line, err := o.readLine()
	if err != nil {
		return nil, err
	}

	record, err := o.unmarshalRecord(line)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (o *Source) Collect() (result RecordList) {
	defer o.waitErrorListener.Wait()

	record, err := o.parseFirstLine()
	if err != nil {
		o.raiseError(err)
		return
	}

	result.Add(record)

	for {
		line, err := o.readLine()
		if err != nil {
			if err != io.EOF {
				o.raiseError(err)
			}
			return
		}

		record, err := o.unmarshalRecord(line)
		if err != nil {
			o.raiseError(err)
			continue
		}

		if record == nil {
			continue
		}

		if !(func() bool {
			for _, filter := range o.filters {
				if !filter(record) {
					return false
				}
			}
			return true
		})() {
			continue
		}

		result.Add(record)
	}

	return
}
