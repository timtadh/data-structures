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
type KIterator func()(key Equatable, next KIterator)
type KVIterator func()(key Equatable, value interface{}, next KVIterator)
type Coroutine func(send interface{})(recv interface{}, next Coroutine)

type Iterable interface {
    Iterate() Iterator
}

type KIterable interface {
    Keys() KIterator
}

type VIterable interface {
    Values() Iterator
}

type KVIterable interface {
    Iterate() KVIterator
}

type MapIterable interface {
    KIterable
    VIterable
    KVIterable
}

type Sized interface {
    Size() int
}

type MapOperable interface {
    Sized
    Has(key Hashable) bool
    Put(key Hashable, value interface{}) (err error)
    Get(key Hashable) (value interface{}, err error)
    Remove(key Hashable) (value interface{}, err error)
}

type Map interface {
    MapIterable
    MapOperable
}

type Tree interface {
    Root() TreeNode
}

type TreeMap interface {
    Tree
    Map
}

type TreeNode interface {
    Key() Equatable
    Value() interface{}
    Children() TreeNodeIterator
    GetChild(int) TreeNode // if your tree can't support this simply panic
                           // many of the utility functions do not require this
                           // however, it is recommended that you implement it
                           // if possible (for instance, post-order traversal
                           // requires it).
    ChildCount() int // a negative value indicates this tree can't provide
                     // an accurate count.
}
type TreeNodeIterator func() (node TreeNode, next TreeNodeIterator)

type BinaryTreeNode interface {
    TreeNode
    Left() BinaryTreeNode
    Right() BinaryTreeNode
}
