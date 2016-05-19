package exc

import (
	"fmt"
	"runtime"
	"reflect"
	"strings"
)

import (
	"github.com/timtadh/data-structures/errors"
)

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

type catch struct {
	exception reflect.Type
	catch func(Throwable)
}

type Block struct{
	try func()
	catches []catch
	finally func()
}

func Try(try func()) *Block {
	return &Block{try: try}
}

func (b *Block) Catch(exc Throwable, do func(Throwable)) *Block {
	b.catches = append(b.catches, catch{reflect.TypeOf(exc), do})
	return b
}

func (b *Block) Finally(finally func()) *Block {
	b.finally = finally
	return b
}

func (b *Block) Unwind() {
	err := b.run()
	if err != nil {
		panic(err)
	}
}

func (b *Block) Error() error {
	err := b.run()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (b *Block) Exception() Throwable {
	t := b.run()
	if t != nil {
		return t
	} else {
		return nil
	}
}

func (b *Block) run() (Throwable) {
	err := b.exec()
	if err == nil {
		return nil
	}
	t := reflect.TypeOf(err)
	for _, c := range b.catches {
		errors.Logf("DEBUG", "try catch %v with %v", t, c.exception)
		if t.AssignableTo(c.exception) {
			err = Try(func(){c.catch(err)}).exec()
			break
		}
	}
	if b.finally != nil {
		b.finally()
	}
	return err
}

func (b *Block) exec() (err Throwable) {
	defer func() {
		if e := recover(); e != nil {
			switch exc := e.(type) {
			case Throwable:
				err = exc
			default:
				panic(e)
			}
		}
	}()
	b.try()
	return
}

