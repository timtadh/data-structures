package exc

import (
	"fmt"
	"runtime"
	"strings"
)

import ()

type Error struct {
	Err   error
	Stack []byte
}

func Errorf(format string, args ...interface{}) *Error {
	buf := make([]byte, 50000)
	n := runtime.Stack(buf, true)
	trace := make([]byte, n)
	copy(trace, buf)
	return &Error{
		Err:  fmt.Errorf(format, args...),
		Stack: trace,
	}
}

func FromError(err error) *Error {
	buf := make([]byte, 50000)
	n := runtime.Stack(buf, true)
	trace := make([]byte, n)
	copy(trace, buf)
	return &Error{
		Err: err,
		Stack: trace,
	}
}

func (e *Error) Error() string {
	if e == nil {
		return "Error <nil>"
	} else {
		return fmt.Sprintf("%v\n\n%s", e.Err, string(e.Stack))
	}
}

func (e *Error) String() string {
	return e.Error()
}

func (e *Error) Exception() *Exception {
	return &Exception{[]*Error{e}}
}

type Throwable interface {
	Exc()              *Exception
	Error()            string
	Chain(e *Error)    Throwable
}

type Exception struct {
	Errors []*Error
}

func (e *Exception) Exc() *Exception {
	return e
}

func (e *Exception) Error() string {
	errs := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		errs = append(errs, err.Error())
	}
	return fmt.Sprintf("Exception\n\n%v\n\nEnd Exception", strings.Join(errs, "\n---\n"))
}

func (e *Exception) String() string {
	return e.Error()
}

func (e *Exception) Chain(err *Error) Throwable {
	e.Errors = append(e.Errors, err)
	return e
}

func (e *Exception) Join(exc *Exception) *Exception {
	errs := make([]*Error, 0, len(e.Errors) + len(exc.Errors))
	errs = append(errs, e.Errors...)
	errs = append(errs, exc.Errors...)
	return &Exception{
		Errors: errs,
	}
}

type NamedException struct {
	*Exception
	Nomen string
}

func (e *NamedException) Name() string {
	return e.Nomen
}

func (e *NamedException) Error() string {
	errs := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		errs = append(errs, err.Error())
	}
	return fmt.Sprintf("Exception %v\n\n%v\n\nEnd Exception", e.Nomen, strings.Join(errs, "\n---\n"))
}

func (e *NamedException) Join(exc *Exception) *NamedException {
	return &NamedException{
		Exception: e.Exception.Join(exc),
		Nomen: e.Nomen,
	}
}

func Throwf(format string, args ...interface{}) {
	ThrowErr(Errorf(format, args...))
}

func ThrowErr(e *Error) {
	throw(e.Exception())
}

func Throw(e Throwable) {
	throw(e)
}

func Rethrow(e Throwable, err *Error) {
	throw(e.Chain(err))
}

func throw(e Throwable) {
	panic(e)
}
