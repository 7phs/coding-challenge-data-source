package helper

import (
	"errors"
	"strings"
	"testing"
)

func TestErrList(t *testing.T) {
	var errList ErrList

	if exist := errList.Finish(); exist != nil {
		t.Error("empty list should be a nil")
	}

	errTitles := []string{
		"err1",
		"err2",
		"cool: err",
	}

	for _, s := range errTitles {
		errList.Add(errors.New(s))
	}

	expected := strings.Join(errTitles, "; ")
	if exist := errList.Finish().Error(); exist != expected {
		t.Error("failed to stringify an error list. Got '" + exist + "', but expected is '" + expected + "'")
	}
}
