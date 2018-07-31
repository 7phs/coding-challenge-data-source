package data

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sync"
)

// An interface for validation record, which required for Reader.
type ValidatedRecord interface {
	Validate() error
}

// Short descriptions for callback function
type RecordFabric func() ValidatedRecord
type RecordFilter func(interface{}) bool
type ErrorListener func(error)

// A data source is using to read data and parsing lines to records.
// A source is possible filtering a parsed record using a user callback.
// Supporting automatically find a good fit unmarshaller from supported.
type Source struct {
	reader            *bufio.Reader
	recordFabric      RecordFabric
	unMarshaling      RecordUnMarshaling
	filters           []RecordFilter
	errorListen       []ErrorListener
	waitErrorListener sync.WaitGroup

	recordIndex int
}

// A constructor of the source
func NewSource(reader io.Reader, fabric RecordFabric) *Source {
	o := &Source{
		reader:       bufio.NewReader(reader),
		recordFabric: fabric,
	}

	return o
}

// A
func (o *Source) matchUnMarshaling(line []byte) (err error) {
	o.unMarshaling, err = MatchUnMarshaling(line)

	return err
}

// Register a user function as a filter callback.
func (o *Source) Filter(filter RecordFilter) *Source {
	o.filters = append(o.filters, filter)

	return o
}

// Calling all user callbacks in a goroutine to transfer an error into it.
func (o *Source) raiseError(err error) error {
	if o.recordIndex > 0 {
		err = fmt.Errorf("line #%d: %s", o.recordIndex, err)
	}

	if len(o.errorListen) > 0 {
		o.waitErrorListener.Add(1)
		go func() {
			defer o.waitErrorListener.Done()

			for _, errCallback := range o.errorListen {
				errCallback(err)
			}
		}()
	}

	return err
}

// Register a function as an error listener.
func (o *Source) Catch(errCallback ErrorListener) *Source {
	o.errorListen = append(o.errorListen, errCallback)

	return o
}

// Reading a whole line from the data without dividing on a pieces.
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

			return nil, err
		}

		line = append(line, piece...)
	}

	return line, nil
}

// Unmarshal a read line into a required record using an automatically detected unmarshaller.
func (o *Source) unmarshalRecord(line []byte) (interface{}, error) {
	record := o.recordFabric()
	if err := o.unMarshaling(line, record); err != nil {
		return nil, err
	}

	if err := record.Validate(); err != nil {
		return nil, err
	}

	return record, nil
}

// Reading a line and unmarshaling it into a record.
func (o *Source) readRecord() (interface{}, error) {
	line, err := o.readLine()
	if err != nil {
		if err != io.EOF {
			return nil, err
		}
		return nil, io.EOF
	}

	if len(bytes.TrimSpace(line)) == 0 {
		return nil, nil
	}

	// the first time checking marshaling
	if o.recordIndex == 0 {
		if err := o.matchUnMarshaling(line); err != nil {
			return nil, err
		}
	}

	o.recordIndex++

	record, err := o.unmarshalRecord(line)
	if err != nil {
		// just skip an unmarshaling error, but reporting listeners
		o.raiseError(err)
		return nil, nil
	}

	return record, nil
}

// A primary workflow. Reading records from data, filtering it and collecting its into a result list.
func (o *Source) Collect() (result RecordList) {
	defer o.waitErrorListener.Wait()

	for {
		record, err := o.readRecord()
		if err != nil {
			if err != io.EOF {
				o.raiseError(err)
			}

			return
		} else if record == nil {
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
