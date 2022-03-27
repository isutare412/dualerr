package dualerr

import (
	"fmt"
	"runtime"
	"strings"
)

var Err = Error{}

// Error has two errors. One as embedded, which enables Error work just
// like go vanilla error. The other is simple error, which carries simple error
// for client report.
type Error struct {
	error
	simple error
}

// New creates DualError. err is embedded into DualError. So err works just like
// a go error. format and args are passed to fmt.Errorf and build a simple
// error. This simple error is carried though an error propagation without
// wrapping.
func New(err error, format string, args ...interface{}) error {
	pcs := make([]uintptr, 1)
	depth := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	caller, _ := frames.Next()
	funcName := funcName(caller.Function)
	return Error{
		error:  fmt.Errorf("%s: %w", funcName, err),
		simple: fmt.Errorf(format, args...),
	}
}

// Wrap wrpas err with callers function name.
func Wrap(err error) error {
	pcs := make([]uintptr, 1)
	depth := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	caller, _ := frames.Next()
	funcName := funcName(caller.Function)
	return fmt.Errorf("%s: %w", funcName, err)
}

// Is implements interface for errors.Is.
func (e Error) Is(err error) bool {
	_, ok := err.(Error)
	return ok
}

// Simple returns unwrapped simple error.
func (e Error) Simple() error {
	return e.simple
}

func funcName(f string) string {
	lastSlash := strings.LastIndex(f, "/")
	if lastSlash < 0 {
		return f
	}
	return f[lastSlash+1:]
}
