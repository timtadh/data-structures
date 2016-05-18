package exc

import (
	"fmt"
	"reflect"
)

import (
	"github.com/timtadh/data-structures/errors"
)

type Exception interface {
	Error() string
	Name() string
}

type BaseException struct {
	Err *errors.Error
	Id string
}

func (e *BaseException) Error() string {
	return fmt.Sprintf("Exception: %v\n%v", e.Id, e.Err.Error())
}

func (e *BaseException) String() string {
	return e.Error()
}

func (e *BaseException) Name() string {
	return e.Id
}

func Throwf(name, format string, args ...interface{}) {
	throw(&BaseException{
		Err: errors.Errorf(format, args...).(*errors.Error),
		Id: name,
	})
}

func Throw(e Exception) {
	throw(e)
}

func throw(e Exception) {
	panic(e)
}

type catch struct {
	exception reflect.Type
	catch func(Exception)
}

type Block struct{
	try func()
	catches []catch
	finally func()
}

func Try(try func()) *Block {
	return &Block{try: try}
}

func (b *Block) Catch(exc Exception, do func(Exception)) *Block {
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

func (b *Block) Error() Exception {
	return b.run()
}

func (b *Block) run() (Exception) {
	err := b.exec()
	if err == nil {
		return nil
	}
	t := reflect.TypeOf(err)
	for _, c := range b.catches {
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

func (b *Block) exec() (err Exception) {
	defer func() {
		if e := recover(); e != nil {
			switch exc := e.(type) {
			case Exception:
				err = exc
			default:
				panic(e)
			}
		}
	}()
	b.try()
	return
}

