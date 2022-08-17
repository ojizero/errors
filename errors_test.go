package errors

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorTypeWrapping(t *testing.T) {
	t.Run("Error type unwrapping nothing", testErrorTypeWrappingNothing)
	t.Run("Error type wrapping one error", testErrorTypeWrappingOneError)
	t.Run("Error type wrapping multiple errors", testErrorTypeWrappingMultipleErrors)
	t.Run("WrappedErrors type behaviours", testWrappedErrorsTypeBehaviour)
	t.Run("WrappedErrors type panics", testWrappedErrorsTypePanics)
}

func testErrorTypeWrappingNothing(t *testing.T) {
	err := New("custom error")

	assert.True(t, Unwrap(err) == nil, "must unwrap to an untyped nil if it wraps nothing")
	assert.Nil(t, Unwrap(err), "must unwrap to nil if it wraps nothing")
}

func testErrorTypeWrappingOneError(t *testing.T) {
	errRaw := String("stdlib error")
	errCustom := New("custom error").(contextualError)
	wrapping := Wrap(errCustom, errRaw)

	assert.IsType(t, contextualError{}, wrapping)
	assert.True(t, equivalent(errCustom, wrapping.(contextualError)), "Error after calling Wrap remains the same as original Error")
	assert.ErrorIs(t, wrapping, errRaw, "Error after wrapping passes `errors.Is' on it's raw wrapped error")
	assert.True(t, Is(wrapping, errCustom), "Error after calling Wrap error still passes `errors.Is' on it's previous value")
	unwrapped := Unwrap(wrapping)
	fmt.Fprintf(os.Stderr, "GOT (%v)", unwrapped)
	assert.Equal(t, errRaw, unwrapped, "unwrapping returns original error")
}

func testErrorTypeWrappingMultipleErrors(t *testing.T) {
	errRaw1 := String("stdlib error 1")
	errRaw2 := String("stdlib error 2")
	errCustom := New("custom error").(contextualError)
	wrapping := Wrap(errCustom, errRaw1, errRaw2)

	assert.IsType(t, contextualError{}, wrapping)
	assert.True(t, equivalent(errCustom, wrapping.(contextualError)), "Error after calling Wrap remains the same as original Error")
	assert.ErrorIs(t, wrapping, errRaw1, "Error after wrapping passes `errors.Is' on it's raw wrapped error")
	assert.ErrorIs(t, wrapping, errRaw2, "Error after wrapping passes `errors.Is' on it's raw wrapped error")
	assert.True(t, Is(wrapping, errCustom), "Error after calling Wrap error still passes `errors.Is' on it's previous value")
	unwrapped := Unwrap(wrapping)
	assert.IsType(t, &wrappedErrors{}, unwrapped, "when wrapping multiple errors calling Unwrap yields WrappedErrors type")
	assert.ErrorIs(t, unwrapped, errRaw1, "WrappedErrors passes `errors.Is' on it's wrapped raw errors")
	assert.ErrorIs(t, unwrapped, errRaw2, "WrappedErrors passes `errors.Is' on it's wrapped raw errors")
}

func testWrappedErrorsTypeBehaviour(t *testing.T) {
	errRaw1 := String("stdlib error 1")
	errRaw2 := String("stdlib error 2")
	errRaw3 := String("stdlib error 3")
	errRaw4 := String("stdlib error 4")
	wrapped1 := Wrap(errRaw1, errRaw2)
	wrapped2 := Wrap(errRaw1, errRaw2, errRaw3)
	wrapped3 := Wrap(wrapped1, errRaw3, errRaw4)

	// Assert Wrap is no-op when called with one error
	assert.Equal(t, errRaw1, Wrap(errRaw1), "calling Wrap on one error returns same error back (acts as no-op)")

	// Assert correctness of the returned error type from Wrap
	assert.IsType(t, wrappedErrors{}, wrapped1, "calling Wrap with multiple errors returns back a WrappedError type")
	assert.IsType(t, wrappedErrors{}, wrapped2, "calling Wrap with multiple errors returns back a WrappedError type")
	assert.IsType(t, wrappedErrors{}, wrapped3, "calling Wrap with multiple errors returns back a WrappedError type")

	// Assert simple errors wrapped and the overall behaviour or WrappedErrors
	// it should be transparent for the most part and one unwrap returns us
	// back to working with the original errors
	//
	assert.EqualError(t, wrapped1, errRaw1.Error(), "wrapped error should be identical to head of the wrapped errors")
	assert.EqualError(t, Unwrap(wrapped1), errRaw2.Error(), "unwrapped error should be identical to the next head of the wrapped errors")
	assert.Equal(t, Unwrap(wrapped1), errRaw2, "unwrapping that leaves 1 error yields that error as is without wrapping with WrappedErrors type")

	// Assert slightly more complex cases of error wrapping of 3 errors
	assert.EqualError(t, wrapped2, errRaw1.Error(), "wrapped error should be identical to head of the wrapped errors")
	assert.EqualError(t, Unwrap(wrapped2), errRaw2.Error(), "unwrapped error should be identical to the next head of the wrapped errors")
	assert.EqualError(t, Unwrap(Unwrap(wrapped2)), errRaw3.Error(), "unwrapped error should be identical to the next head of the wrapped errors")
	assert.IsType(t, wrappedErrors{}, Unwrap(wrapped2), "unwrapping that leaves more than 1 error yields a new WrappedErrors type")

	// Assert more complex cases of wrapping errors that already wrap other errors,
	// the overall behaviour should still be transparent for the most part.
	//
	assert.EqualError(t, wrapped3, errRaw1.Error(), "wrapped error should be identical to head of the wrapped errors")
	assert.EqualError(t, Unwrap(wrapped3), errRaw2.Error(), "unwrapped error should be identical to the next head of the wrapped errors")
	assert.EqualError(t, Unwrap(Unwrap(wrapped3)), errRaw3.Error(), "unwrapped error should be identical to the next head of the wrapped errors")
	assert.EqualError(t, Unwrap(Unwrap(Unwrap(wrapped3))), errRaw4.Error(), "unwrapped error should be identical to the next head of the wrapped errors")
	assert.Equal(t, Unwrap(Unwrap(Unwrap(wrapped3))), errRaw4, "final unwrap call returns the original error")
}

func testWrappedErrorsTypePanics(t *testing.T) {
	impossibleErr := wrappedErrors{}
	// We try to print the value here since Go complains if we don't
	// use the result of calling `Error` ¯\_(ツ)_/¯.
	assert.Panics(t, func() { fmt.Println(impossibleErr.Error()) }, "WrappedErrors is not an error type in and of itself and cannot return an error message with nothing wrapped in it")
	assert.Panics(t, func() { Unwrap(impossibleErr) }, "WrappedErrors is not an error type in and of itself and cannot return a wrapped error if nothing is wrapped in it")
}
