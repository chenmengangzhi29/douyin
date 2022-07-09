package errno

import (
	"fmt"
	"runtime"
)

type stack []uintptr

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

type withCode struct {
	err   error
	code  int
	cause error
	msg   string
	*stack
}

func WithCode(code int, format string, args ...interface{}) error {
	return &withCode{
		err:   fmt.Errorf(format, args...),
		code:  code,
		msg:   format,
		stack: callers(),
	}
}

func WrapC(err error, code int, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &withCode{
		err:   fmt.Errorf(format, args...),
		code:  code,
		cause: err,
		msg:   format,
		stack: callers(),
	}
}

func (w *withCode) Error() string { return w.msg }

func (w *withCode) Cause() error { return w.cause }

func (w *withCode) Unwrap() error { return w.cause }
