package list

import (
	"bytes"
	"encoding/binary"
	"hash/fnv"
)

import (
	"github.com/timtadh/data-structures/types"
	"github.com/timtadh/data-structures/errors"
)


type List struct {
	list []types.Hashable
}

type ItemMarshal func(types.Hashable) ([]byte, error)
type ItemUnmarshal func([]byte) (types.Hashable, error)

type MList struct {
	List
	MarshalItem ItemMarshal
	UnmarshalItem ItemUnmarshal
}

func NewMList(list *List, marshal ItemMarshal, unmarshal ItemUnmarshal) *MList {
	return &MList{
		List: List{
			list: list.list,
		},
		MarshalItem: marshal,
		UnmarshalItem: unmarshal,
	}
}

func (m *MList) MarshalBinary() ([]byte, error) {
	items := make([][]byte, 0, m.Size())
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(m.Size()))
	items = append(items, size)
	for item, next := m.Items()(); next != nil; item, next = next() {
		b, err := m.MarshalItem(item)
		if err != nil {
			return nil, err
		}
		size := make([]byte, 4)
		binary.LittleEndian.PutUint32(size, uint32(len(b)))
		items = append(items, size, b)
	}
	return bytes.Join(items, []byte{}), nil
}

func (m *MList) UnmarshalBinary(bytes []byte) (error) {
	size := int(binary.LittleEndian.Uint32(bytes[0:4]))
	off := 4
	m.list = make([]types.Hashable, 0, size)
	for i := 0; i < size; i++ {
		s := off
		e := off + 4
		size := int(binary.LittleEndian.Uint32(bytes[s:e]))
		s = e
		e = s + size
		item, err := m.UnmarshalItem(bytes[s:e])
		if err != nil {
			return err
		}
		m.Append(item)
		off = e
	}
	return nil
}

func New(initialSize int) *List {
	return &List{
		list: make([]types.Hashable, 0, initialSize),
	}
}

func FromSlice(list []types.Hashable) *List {
	l := &List{
		list: make([]types.Hashable, len(list)),
	}
	copy(l.list, list)
	return l
}

func (l *List) Clear() {
	l.list = l.list[:0]
}

func (l *List) Size() int {
	return len(l.list)
}

func (l *List) Has(item types.Hashable) (has bool) {
	for i := range l.list {
		if l.list[i].Equals(item) {
			return true
		}
	}
	return false
}

func (l *List) Equals(b types.Equatable) bool {
	if o, ok := b.(types.IterableContainer); ok {
		return l.equals(o)
	} else {
		return false
	}
}

func (l *List) equals(o types.IterableContainer) bool {
	if l.Size() != o.Size() {
		return false
	}
	for v, next := l.Items()(); next != nil; v, next = next() {
		item := v.(types.Hashable)
		if !o.Has(item) {
			return false
		}
	}
	return true
}

func (l *List) Less(b types.Sortable) bool {
	if o, ok := b.(types.IterableContainer); ok {
		return l.less(o)
	} else {
		return false
	}
}

func (l *List) less(o types.IterableContainer) bool {
	if l.Size() < o.Size() {
		return true
	} else if l.Size() > o.Size() {
		return false
	}
	cs, si := l.Items()()
	co, oi := o.Items()()
	for si != nil || oi != nil {
		if cs.Less(co) {
			return true
		} else if !cs.Equals(co) {
			return false
		}
		cs, si = si()
		co, oi = oi()
	}
	return false
}

func (l *List) Hash() int {
	h := fnv.New32a()
	if len(l.list) == 0 {
		return 0
	}
	bs := make([]byte, 4)
	for _, item := range l.list {
		binary.LittleEndian.PutUint32(bs, uint32(item.Hash()))
		h.Write(bs)
	}
	return int(h.Sum32())
}

func (l *List) Items() (it types.KIterator) {
	i := 0
	return func() (item types.Hashable, next types.KIterator) {
		if i < len(l.list) {
			item = l.list[i]
			i++
			return item, it
		}
		return nil, nil
	}
}

func (l *List) Get(i int) (item types.Hashable, err error) {
	if i < 0 || i >= len(l.list) {
		return nil, errors.Errorf("Access out of bounds. len(*List) = %v, idx = %v", len(l.list), i)
	}
	return l.list[i], nil
}

func (l *List) Set(i int, item types.Hashable) (err error) {
	if i < 0 || i >= len(l.list) {
		return errors.Errorf("Access out of bounds. len(*List) = %v, idx = %v", len(l.list), i)
	}
	l.list[i] = item
	return nil
}

func (l *List) Append(item types.Hashable) error {
	return l.Insert(len(l.list), item)
}

func (l *List) Insert(i int, item types.Hashable) error {
	if i < 0 || i > len(l.list) {
		return errors.Errorf("Access out of bounds. len(*List) = %v, idx = %v", len(l.list), i)
	}
	if len(l.list) == cap(l.list) {
		l.expand()
	}
	l.list = l.list[:len(l.list)+1]
	for j := len(l.list) - 1; j > 0; j-- {
		if j == i {
			l.list[i] = item
			break
		}
		l.list[j] = l.list[j-1]
	}
	if i == 0 {
		l.list[i] = item
	}
	return nil
}

func (l *List) Extend(it types.KIterator) (err error) {
	for item, next := it(); next != nil; item, next = next() {
		if err := l.Append(item); err != nil {
			return err
		}
	}
	return nil
}

func (l *List) Pop() (item types.Hashable, err error) {
	item, err = l.Get(len(l.list)-1)
	if err != nil {
		return nil, err
	}
	err = l.Remove(len(l.list)-1)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (l *List) Remove(i int) error {
	if i < 0 || i >= len(l.list) {
		return errors.Errorf("Access out of bounds. len(*List) = %v, idx = %v", len(l.list), i)
	}
	dst := l.list[i:len(l.list)-1]
	src := l.list[i+1:len(l.list)]
	copy(dst, src)
	l.list = l.list[:len(l.list)-1]
	l.shrink()
	return nil
}

func (l *List) expand() {
	list := l.list
	if cap(list) < 100 {
		l.list = make([]types.Hashable, len(list), cap(list)*2)
	} else {
		l.list = make([]types.Hashable, len(list), cap(list)+100)
	}
	copy(l.list, list)
}

func (l *List) shrink() {
	if (len(l.list)-1)*2 >= cap(l.list) || cap(l.list)/2 <= 10 {
		return
	}
	list := l.list
	l.list = make([]types.Hashable, len(list), cap(list)/2)
	copy(l.list, list)
}
