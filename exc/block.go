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

type Block struct{
	try func()
	catches []catch
	finallies []func()
}

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
			if t.AssignableTo(c.exception) {
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

