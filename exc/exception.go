package exc

import (
	"fmt"
	"runtime"
	"strings"
)

import ()

// The *Error struct wraps up a regular error (could be any error) and a stack
// trace. The idea is that it track where the error was created.
type Error struct {
	Err   error
	Stack []byte
}

// A drop in replacement for either fmt.Errorf or errors.Errorf. Creates and
// *Error with a stack trace from where Errorf() was called from.
func Errorf(format string, args ...interface{}) *Error {
	buf := make([]byte, 50000)
	n := runtime.Stack(buf, true)
	trace := make([]byte, n)
	copy(trace, buf)
	return &Error{
		Err:   fmt.Errorf(format, args...),
		Stack: trace,
	}
}

// Create an *Error from an existing error value. It will create a stack trace
// to attach to the error from where FromError() was called from.
func FromError(err error) *Error {
	buf := make([]byte, 50000)
	n := runtime.Stack(buf, true)
	trace := make([]byte, n)
	copy(trace, buf)
	return &Error{
		Err:   err,
		Stack: trace,
	}
}

// Format the error and stack trace.
func (e *Error) Error() string {
	if e == nil {
		return "Error <nil>"
	} else {
		return fmt.Sprintf("%v\n\n%s", e.Err, string(e.Stack))
	}
}

// Format the error and stack trace.
func (e *Error) String() string {
	return e.Error()
}

// Create an Exception object from the Error.
func (e *Error) Exception() *Exception {
	return &Exception{[]*Error{e}}
}

// The interface that represents what can be thrown (and caught) by this
// library. All Throwables must be convertible to and *Exception, implement the
// error interface and allow chaining on of extra errors so that they can be
// Rethrown.
type Throwable interface {
	Exc() *Exception
	Error() string
	Chain(e *Error) Throwable
}

// An implementation of Throwable you can base your custom Exception types off
// of. It is also the type of Throwable thrown by Throwf. To "inherit" from
// Exception use this formula:
//
//     type MyException struct {
//     	exc.Exception
//     }
//
// This ensures that your new exception will be catchable when *Exception is
// supplied. See *Block.Catch for details.
//
type Exception struct {
	Errors []*Error
}

// Return itself.
func (e *Exception) Exc() *Exception {
	return e
}

// Format an error string from the list of *Error.
func (e *Exception) Error() string {
	errs := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		errs = append(errs, err.Error())
	}
	return fmt.Sprintf("Exception\n\n%v\n\nEnd Exception", strings.Join(errs, "\n---\n"))
}

// Format an error string from the list of *Error.
func (e *Exception) String() string {
	return e.Error()
}

// Add another *Error to the list of *Error
func (e *Exception) Chain(err *Error) Throwable {
	e.Errors = append(e.Errors, err)
	return e
}

// Join this exception with another exception.
func (e *Exception) Join(exc *Exception) *Exception {
	errs := make([]*Error, 0, len(e.Errors)+len(exc.Errors))
	errs = append(errs, e.Errors...)
	errs = append(errs, exc.Errors...)
	return &Exception{
		Errors: errs,
	}
}

// Throw a new *Exception created from an *Error made with Errorf. Basically a
// drop in replacement for everywhere you used fmt.Errorf but would now like to
// throw an exception without creating a custom exception type.
func Throwf(format string, args ...interface{}) {
	ThrowErr(Errorf(format, args...))
}

// Throw an *Exception created from an *Error.
func ThrowErr(e *Error) {
	throw(e.Exception())
}

// Throw an *Exception created from a regular error interface object. The stack
// trace for the exception will point to the throw rather than the creation of
// the error object.
func ThrowOnError(err error) {
	if err != nil {
		ThrowErr(FromError(err))
	}
}

// Throw a Throwable object. This is how you throw a custom exception:
//
//     type MyException struct {
//     	exc.Exception
//     }
//
//     exc.Try(func() {
//     	exc.Throw(&MyException{*Errorf("My Exception").Exception()})
//     }).Catch(&Exception{}, func(t Throwable) {
//     	log.Log("caught!")
//     }).Unwind()
//
func Throw(e Throwable) {
	throw(e)
}

// Rethrow an object. It chains on the *Error as the reason for the rethrow and
// where it occured. If you are inside of a catch block you should use this
// method instead of Throw*
//
//    exc.Try(func() {
//    	exc.Throwf("wat!@!")
//    }).Catch(&Exception{}, func(e Throwable) {
//    	t.Log("Caught", e)
//    	exc.Rethrow(e, Errorf("rethrow"))
//    }).Unwind()
//
func Rethrow(e Throwable, err *Error) {
	throw(e.Chain(err))
}

func throw(e Throwable) {
	panic(e)
}
