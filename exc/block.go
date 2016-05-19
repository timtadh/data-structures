package exc

import (
	"io"
	"reflect"
)

import ()

type catch struct {
	exception reflect.Type
	catch func(Throwable)
}

// Represents a Try/Catch/Finally block. Created by a Try or Close function.
// Executed by Unwind, Error, or Exception. You may add multiple Catch and
// Finally expressions to your *Block with the  Catch and Finally functions.
// Catch and Finally can be in arbitrary orders unlike in languages with
// syntatic support for Exceptions. I recommend you put your Catch expressions
// before your finally Expressions. Finally and Catch expressions are evaluated
// in the order they are declared.
type Block struct{
	try func()
	catches []catch
	finallies []func()
}

// Start a Try/Catch/Finally block. Any exception thrown in the Try block can
// be caught by a Catch block (if registered for the Exception type or a parent
// exception type).
func Try(try func()) *Block {
	return &Block{try: try}
}

func Close(makeCloser func() io.Closer, try func(io.Closer)) *Block {
	var c io.Closer = nil
	return Try(func(){
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

func (b *Block) Catch(exc Throwable, do func(Throwable)) *Block {
	b.catches = append(b.catches, catch{reflect.TypeOf(exc), do})
	return b
}

func (b *Block) Finally(finally func()) *Block {
	b.finallies = append(b.finallies, finally)
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
	if err != nil {
		t := reflect.TypeOf(err)
		for _, c := range b.catches {
			// errors.Logf("DEBUG", "trying to catch %v with %v %v %v", t, c.exception, t.ConvertibleTo(c.exception), t.AssignableTo(c.exception))
			if isa(t, c.exception) {
				err = Try(func(){c.catch(err)}).exec()
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

