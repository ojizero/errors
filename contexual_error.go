package errors

import "runtime"

// contextualError is a custom error type we use to provide us with additional
// context around any app-level errors.
//
// It captures the line and file of the error, allows for tagging
// it with custom labels, and supports wrapping other errors.
//
type contextualError struct {
	msg     string
	wrapped error

	file string
	line int

	labels []string
}

// Validate we implement Go's stdlib errors package useful interfaces.
var (
	_ error                       = contextualError{}
	_ interface{ Is(error) bool } = contextualError{}
	_ interface{ Unwrap() error } = contextualError{}

	// To keep in with the style of stdlib's `errors' package we define a new interface
	// with `Wrap', `With', `Labels', and `LabeledBy' methods.
	//
	_ interface{ Wrap(...error) error }   = contextualError{}
	_ interface{ With(...string) error }  = contextualError{}
	_ interface{ Labels() []string }      = contextualError{}
	_ interface{ LabeledBy(string) bool } = contextualError{}
)

// New generates a new Error instance, it captures the caller
// information as well and attaches them to the generated
// error.
//
// Any errors passed here will be wrapped within the generated
// error (same as calling Wrap) on them afterwards.
//
func New(msg string, errs ...error) error {
	err := contextualError{msg: msg}
	if _, file, line, ok := runtime.Caller(1); ok {
		err.file = file
		err.line = line
	}
	return err.Wrap(errs...)
}

// Error implements Go stdlib's error interface.
func (err contextualError) Error() string {
	return err.msg
}

// With appends labels to our Error type, labels can be used to organize
// your errors, or error handlers.
//
func (err contextualError) With(labels ...string) error {
	err.labels = append(err.labels, labels...)
	return err
}

// Labels returns all assigned labels to the given error
func (err contextualError) Labels() []string {
	return err.labels
}

// LabeledBy checks whether or not a given label is assigned to the error.
func (err contextualError) LabeledBy(label string) bool {
	for _, l := range err.labels {
		if l == label {
			return true
		}
	}
	return false
}

// Wrap can be used to wrap additional errors inside out Error type.
func (err contextualError) Wrap(errs ...error) error {
	if len(errs) == 0 {
		return err
	}
	if err.wrapped == nil {
		err.wrapped = &wrappedErrors{errs}
		return err
	}
	err.wrapped = Wrap(err.wrapped, errs...)
	return err
}

// Unwrap returns the wrapped errors back, it returns a proxy that
// behaves the same as the head of the wrapped errors but if
// unwrapped it will unwrap the next error in the chain.
//
func (err contextualError) Unwrap() error {
	// In case we are wrapping only one error no need to dance
	// around it and return the wrapped custom type, just
	// return the underlying wrapped original error.
	//
	if wrapping, ok := err.wrapped.(*wrappedErrors); err.wrapped != nil && ok && len(wrapping.errs) == 1 {
		return wrapping.Unwrap()
	}
	return err.wrapped
}

// Is implements `errors.Is` for our custom Error type, it verifies
// that the target is of our custom error type and that it is
// identical to the error it is being compared to.
//
func (err contextualError) Is(target error) bool {
	if target, ok := target.(contextualError); ok {
		return equivalent(err, target) || Is(err.wrapped, target)
	}
	return false
}
