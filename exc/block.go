package exc

import (
	"io"
	"reflect"
)

import ()

type catch struct {
	exception reflect.Type
	catch     func(Throwable)
}

// Represents a Try/Catch/Finally block. Created by a Try or Close function.
// Executed by Unwind, Error, or Exception. You may add multiple Catch and
// Finally expressions to your *Block with the  Catch and Finally functions.
// Catch and Finally can be in arbitrary orders unlike in languages with
// syntatic support for Exceptions. I recommend you put your Catch expressions
// before your finally Expressions. Finally and Catch expressions are evaluated
// in the order they are declared.
type Block struct {
	try       func()
	catches   []catch
	finallies []func()
}

// Start a Try/Catch/Finally block. Any exception thrown in the Try block can
// be caught by a Catch block (if registered for the Exception type or a parent
// exception type).
func Try(try func()) *Block {
	return &Block{try: try}
}

// Start a Try/Catch/Finally block. You supply to functions, the first one
// creates a closable resource and returns it. This resources is supplied to the
// the second function which acts as a normal try. Whether the try fails or
// succeeds the Close() function is always called on the resource that was
// created (and returned) by the first function. Further, catch and finally
// functions may be chained onto the Close function. However, they will be run
// after the Close function is called on the resource and will not have access
// to it.
//
// Finally, if the function to created the closable object fails, it will not be
// cleaned up if it was partially initialized. This is because it was never
// returned. There are also ways to deal with that situation using an outer
// Try/Finally.
func Close(makeCloser func() io.Closer, try func(io.Closer)) *Block {
	var c io.Closer = nil
	return Try(func() {
		Try(func() {
			c = makeCloser()
			try(c)
		}).Finally(func() {
			if c != nil {
				ThrowOnError(c.Close())
			}
		}).Unwind()
	})
}

// Add a catch function for a specific Throwable. If your Throwable struct
// "inherits" from another struct like so:
//
//     type MyException struct {
//     	exc.Exception
//     }
//
// Then you can catch *MyException with *Exception. eg:
//
//     exc.Try(func() {
//     	Throw(&MyException{*Errorf("My Exception").Exception()})
//     }).Catch(&Exception{}, func(t Throwable) {
//     	log.Log("caught!")
//     }).Unwind()
//
// Catch blocks are only run in the case of an thrown exception. Regular panics
// are ignored and will behave as normal.
//
func (b *Block) Catch(exc Throwable, do func(Throwable)) *Block {
	b.catches = append(b.catches, catch{reflect.TypeOf(exc), do})
	return b
}

// Add a finally block. These will be run whether or not an exception was
// thrown. However, if a regular panic occurs this function will not be run and
// the panic will behave as normal.
func (b *Block) Finally(finally func()) *Block {
	b.finallies = append(b.finallies, finally)
	return b
}

// Run the Try/Catch/Finally *Block. If there is an uncaught (or rethrown)
// exception continue to propogate it up the stack (unwinding it). This would be
// the normal behavoir in language which natively support exceptions.
//
// The Block will NOT BE RUN unless this method, Error, or Exception is called.
// This could lead to an difficult to track down bug!
func (b *Block) Unwind() {
	err := b.run()
	if err != nil {
		panic(err)
	}
}

// Run the Try/Catch/Finally *Block. If there is an uncaught (or rethrown)
// exception continue return it as an error.
//
// The Block will NOT BE RUN unless this method, Unwind, or Exception is called.
// This could lead to an difficult to track down bug!
func (b *Block) Error() error {
	err := b.run()
	if err != nil {
		return err
	} else {
		return nil
	}
}

// Run the Try/Catch/Finally *Block. If there is an uncaught (or rethrown)
// exception continue return it as a Throwable.
//
// The Block will NOT BE RUN unless this method, Unwind, or Error is called.
// This could lead to an difficult to track down bug!
func (b *Block) Exception() Throwable {
	t := b.run()
	if t != nil {
		return t
	} else {
		return nil
	}
}

func (b *Block) run() Throwable {
	err := b.exec()
	if err != nil {
		t := reflect.TypeOf(err)
		for _, c := range b.catches {
			// errors.Logf("DEBUG", "trying to catch %v with %v %v %v", t, c.exception, t.ConvertibleTo(c.exception), t.AssignableTo(c.exception))
			if isa(t, c.exception) {
				err = Try(func() { c.catch(err) }).exec()
				break
			}
		}
	}
	for _, finally := range b.finallies {
		finally()
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

func isa(a, b reflect.Type) bool {
	// errors.Logf("DEBUG", "isa a %v b %v a.Kind() %v", a, b, a.Kind())
	if is(a, b) {
		return true
	}
	if !(a.Kind() == reflect.Ptr && b.Kind() == reflect.Ptr) {
		return false
	}
	a = a.Elem()
	for a.Kind() == reflect.Struct && a.NumField() > 0 {
		a = a.Field(0).Type
		if is(reflect.PtrTo(a), b) {
			return true
		}
	}
	return false
}

func is(a, b reflect.Type) bool {
	// errors.Logf("DEBUG", "is a %v b %v %v %v", a, b, a.ConvertibleTo(b), a.AssignableTo(b))
	if a.AssignableTo(b) {
		return true
	}
	return false
}
