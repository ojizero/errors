package errors

// wrappedErrors acts as a proxy of the first error attached to it
// and allows us to wrap multiple errors under one, each error
// acts as if it is wrapped by the one before it and wraps
// all errors after it.
//
type wrappedErrors struct{ errs []error }

// Validate we implement Go's stdlib errors package useful interfaces.
var (
	_ error                       = wrappedErrors{}
	_ interface{ Is(error) bool } = wrappedErrors{}
	_ interface{ Unwrap() error } = wrappedErrors{}

	// To keep in with the style of stdlib's `errors' package we define a new interface
	// with `Wrap' method that aims to mimic the `Unwrap' interface.
	//
	_ interface{ Wrap(...error) error } = wrappedErrors{}
)

// Wrap uses the given error to wrap the other errors given.
//
// If the given error is of our custom type Error or WrappedErrors
// it is used to wrap the error as if you called
// `err.Wrap(errs...)`. Otherwise a proxy
// WrappedErrors error is created.
//
func Wrap(err error, errs ...error) error {
	if len(errs) == 0 {
		return err
	}
	if err, ok := err.(interface{ Wrap(...error) error }); ok {
		return err.Wrap(errs...)
	}
	return wrappedErrors{append([]error{err}, errs...)}
}

// Error is a pass through to the head of the wrapped errors chain.
func (w wrappedErrors) Error() string {
	if len(w.errs) == 0 {
		// This should be an impossible case were if reached it is an undefined state
		// since the WrappedError type exists exclusively to allows us to
		// conveniently wrap other errors and not to be an error
		// in and of itself.
		//
		panic("WrappedErrors undefined state: no wrapped errors")
	}
	return w.errs[0].Error()
}

// Wrap adds a new error to the chain of wrapped errors.
func (w wrappedErrors) Wrap(errs ...error) error {
	w.errs = append(w.errs, errs...)
	return w
}

// Unwrap returns a new wrappedError where the first wrapped error is removed.
//
// This is intended to allow for interloping with Go's stdlib errors
// package cleanly, i.e. the first wrapped error should be the
// error of most importance.
//
// Other wrapped errors act as a chain of one error wrapping all the ones after it.
//
func (w wrappedErrors) Unwrap() error {
	if len(w.errs) == 0 {
		// This should be an impossible case were if reached it is an undefined state
		// since the WrappedError type exists exclusively to allows us to
		// conveniently wrap other errors and not to be an error
		// in and of itself.
		//
		panic("WrappedErrors undefined state: no wrapped errors")
	}
	if len(w.errs) == 1 {
		return w.errs[0]
	}
	if len(w.errs) == 2 {
		return w.errs[1]
	}
	if err := Unwrap(w.errs[0]); err != nil {
		w.errs[0] = err
		return w
	}
	w.errs = w.errs[1:]
	return w
}

// Is returns true if any of the wrapper errors matches the given target.
func (w wrappedErrors) Is(target error) bool {
	for _, e := range w.errs {
		if Is(e, target) {
			return true
		}
	}
	return false
}
