package types

type Equatable interface {
    Equals(b Equatable) bool
}

type Sortable interface {
    Equatable
    Less(b Sortable) bool
}

type Hashable interface {
    Sortable
    Hash() int
}

type Iterator func()(item interface{}, next Iterator)
type KIterator func()(key Equatable, next Iterator)
type KVIterator func()(key Equatable, value interface{}, next KVIterator)
type Coroutine func(send interface{})(recv interface{}, next Coroutine)

type Iterable interface {
    Iterate() Iterable
}

type KVIterable interface {
    Iterate() KVIterator
    Keys() KIterator
}

