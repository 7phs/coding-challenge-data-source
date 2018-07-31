package helper

import "bytes"

// A helper to store a list of error, ex. validation function.
type ErrList []error

// Adding an error to a list
func (o *ErrList) Add(err error) {
	*o = append(*o, err)
}

// Building a string using all of collected errors divided a comma.
func (o ErrList) Error() string {
	if len(o) == 0 {
		return "<nil>"
	}

	buf := bytes.NewBufferString("")
	for i, err := range o {
		if i > 0 {
			buf.WriteString("; ")
		}

		buf.WriteString(err.Error())
	}

	return buf.String()
}

func (o ErrList) Finish() error {
	if len(o) == 0 {
		return nil
	}

	return o
}
