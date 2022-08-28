# errors

A set of helpful error abstractions for use in Go, attempting to be a very thin layer on top the
already versatile Go standard library.

## Usage

    go get github.com/ojizero/errors

Now simply import and use `github.com/ojizero/errors` in places where you would normally use `error`

```go
import "github.com/ojizero/errors"

// This is identical to a standard error except it holds more contextual info
// such as the file and line that defined it, along with the ability to
// add tags and wrap other errors in it.
errors.New("some error")

// Whe creating an error you can simply pass (polyvariadically) other errors
// to be wrapped by it.
errors.New("some error wrapping other errors", errors.New("a wrapper error"))

// You can also use Wrap function to wrap other errors (even standard library ones)
errors.Wrap(err1, err2, ...etc)
```

This library also provides proxy functions to other standard library error functions

- `errors.Fmt` proxies to stdlib's `fmt.Errorf`
- `errors.String` proxies to stdlib's `errors.New`
- `errors.Is` proxies to stdlib's `errors.Is`
- `errors.As` proxies to stdlib's `errors.As`
- `errors.Unwrap` proxies to stdlib's `errors.Unwrap`

Additionally we provide the following set of functionalities that can be used to provide more context
regarding the errors being handled;

- `errors.With` which can be used to attach custom labels to the errors, those labels can later be used while handling the
  errors for example being passed as fingerprints to Sentry while being reported.
- `errors.Labels`, `errors.LabeledBy`, and `errors.LabeledByAny` which can be used to work with labels assigned using
  `errors.With`.
- `errors.DetailedError` can be used to print the error message along with any additional details such as the
  file and line if they are known along with any wrapped errors into one neat message.
- `errors.FileAndLine` which can be used to pull the file and line the error was defined in (if using this package's errors).
