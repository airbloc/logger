package logger

import (
	"github.com/pkg/errors"
	"io"
)

// CloseWithErrCapture is used if you want to close and fail the function or
// method on a `io.Closer.Close()` error (make sure the `error` return argument is
// named as `err`). If the error is already present, `CloseWithErrCapture`
// will append (not wrap) the error caused by `Close` if any.
func CloseWithErrCapture(c io.Closer, errCap *error, msg string) {
	if err := c.Close(); *errCap == nil {
		err = errors.Wrap(err, msg)
	}
}

// CloseWithLogOnErr closes given `io.Closer` and logs error with given message if any.
func CloseWithLogOnErr(log Logger, c io.Closer, msg string) {
	if err := c.Close(); err != nil {
		log.Error(msg, err)
	}
}
