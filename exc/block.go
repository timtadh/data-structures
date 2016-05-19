package exc

import (
	"reflect"
)

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

