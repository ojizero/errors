package errors

import (
	"fmt"
	"strings"
)

// DetailedError builds an error message that is detailed if the provided error is
// can be used to extract the additional details. Mainly useful for logging.
//
func DetailedError(err error) string {
	msg := strings.Builder{}
	if f, l, ok := FileAndLine(err); ok {
		msg.WriteString(fmt.Sprintf("(%s:%d) ", f, l))
	}
	msg.WriteString(err.Error())
	for werr := Unwrap(err); werr != nil; werr = Unwrap(werr) {
		msg.WriteString(fmt.Sprintf("\nCaused by: %s\n", err.Error()))
	}
	return msg.String()
}

// FileAndLine returns the file and the line where the error happened
// if the error is of our custom error type.
//
// Returns a boolean as well to indicate if the error isn't of our
// custom type or if our error doesn't contain the file and line.
//
func FileAndLine(err error) (string, int, bool) {
	if err, ok := err.(contextualError); ok {
		return err.file, err.line, err.file != ""
	}
	return "", 0, false
}

// With appends the given labels to the given error if the
// error is of the custom type Error. Otherwise it returns
// the given error as is.
//
func With(err error, tags ...string) error {
	if err, ok := err.(interface{ With(...string) error }); ok {
		return err.With(tags...)
	}
	return err
}

// Labels returns the labels of the error if it is of our own custom type,
// if any other error is passed it will return an empty list.
//
func Labels(err error) []string {
	if err, ok := err.(interface{ Labels() []string }); ok {
		return err.Labels()
	}
	return []string{}
}

// LabeledBy checks if the given error is a our Error type
// and if so checks if it contains the provided tag.
//
func LabeledBy(err error, label string) bool {
	if e, ok := err.(interface{ LabeledBy(string) bool }); ok {
		return e.LabeledBy(label)
	}
	return false
}

// LabeledByAny checks if the given error is a our Error type
// and if so checks if it contains any of the provided tags.
//
func LabeledByAny(err error, labels ...string) bool {
	if e, ok := err.(interface{ LabeledBy(string) bool }); ok {
		for _, t := range labels {
			if e.LabeledBy(t) {
				return true
			}
		}
	}
	return false
}

func equivalent(a, b contextualError) bool {
	return a.msg == b.msg && a.file == b.file && a.line == b.line
}
