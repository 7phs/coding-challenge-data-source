package helper

import (
	"errors"
	"strings"
	"testing"
)

func TestErrList(t *testing.T) {
	var errList ErrList

	expected := "<nil>"
	if exist := errList.Error(); exist != expected {
		t.Error("failed to stringify an empty error list. Got '" + exist + "', but expected is '" + expected + "'")
	}

	errTitles := []string{
		"err1",
		"err2",
		"cool: err",
	}

	for _, s := range errTitles {
		errList.Add(errors.New(s))
	}

	expected = strings.Join(errTitles, "; ")
	if exist := errList.Error(); exist != expected {
		t.Error("failed to stringify an error list. Got '" + exist + "', but expected is '" + expected + "'")
	}
}
