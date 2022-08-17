package errors

import (
	"errors"
	"fmt"
)

// These are the same error handling functions found in stdlib error package
// but proxied through our own errors package in order to provide more
// convenience and ease of use.
//
var (
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap

	// String proxies over to stdlib's `errors.New'.
	String = errors.New

	// Fmt proxies over to stdlib's `fmt.Errorf'.
	Fmt = fmt.Errorf
)
