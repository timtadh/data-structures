package list

import "testing"

import (
	"fmt"
)

import (
	"github.com/timtadh/data-structures/types"
)

func TestFixed(x *testing.T) {
	t := (*T)(x)
	SIZE := 75
	list := Fixed(SIZE)
	items := make([]types.ByteSlice, 0, SIZE)
	for i := 0; i < SIZE; i++ {
		item := t.randslice(rand.Intn(10) + 1)
		items = append(items, item)
		t.assert_nil(list.Append(item))
	}
	t.assert(fmt.Sprintf("list.Full() (%v)", list.Full()), list.Full())
	t.assert(fmt.Sprintf("list.Size == %d (%d)", SIZE, list.Size()), list.Size() == SIZE)
	err := list.Append(t.randslice(10))
	t.assert(fmt.Sprintf("err != nil (%v)", err), err != nil)
	for i, item := range items {
		lg, err := list.Get(i)
		t.assert_nil(err)
		t.assert(fmt.Sprintf("i %v, items[i] == list.Get(i)", i), lg.Equals(item))
	}
}

func TestFixedSorted(x *testing.T) {
	t := (*T)(x)
	SIZE := 100
	set := NewFixedSorted(SIZE, false)
	items := make([]types.Int, 0, SIZE)
	for i := 0; i < SIZE; i++ {
		item := types.Int(rand.Intn(1000))
		for set.Has(item) {
			item = types.Int(rand.Intn(1000))
		}
		items = append(items, item)
		t.assert_nil(set.Add(item))
	}
	t.assert(fmt.Sprintf("set.Full() (%v)", set.Full()), set.Full())
	t.assert(fmt.Sprintf("set.Size == %d (%d)", SIZE, set.Size()), set.Size() == SIZE)
	for i, item := range items {
		t.assert(fmt.Sprintf("i %v, !set.Has(item)", i), set.Has(item))
	}
	for _, item := range items {
		t.assert_nil(set.Delete(item))
	}
	for i, item := range items {
		t.assert(fmt.Sprintf("i %v, !set.Has(item)", i), !set.Has(item))
	}
}

func TestFixedAppendMarshalUnmarshalGet(x *testing.T) {
	t := (*T)(x)
	SIZE := 100
	list := Fixed(SIZE)
	items := make([]types.Int, 0, SIZE)
	for i := 0; i < SIZE; i++ {
		item := types.Int(rand.Intn(10) + 1)
		items = append(items, item)
		t.assert_nil(list.Append(item))
	}
	for i, item := range items {
		lg, err := list.Get(i)
		t.assert_nil(err)
		t.assert(fmt.Sprintf("i %v, items[i] == list.Get(i)", i), lg.Equals(item))
	}
	err := list.Append(t.randslice(10))
	t.assert(fmt.Sprintf("err != nil (%v)", err), err != nil)
	marshal, unmarshal := types.IntMarshals()
	mlist1 := NewMList(list, marshal, unmarshal)
	bytes, err := mlist1.MarshalBinary()
	t.assert_nil(err)
	mlist2 := &MList{MarshalItem: marshal, UnmarshalItem: unmarshal}
	t.assert_nil(mlist2.UnmarshalBinary(bytes))
	for i, item := range items {
		lg, err := mlist2.Get(i)
		t.assert_nil(err)
		t.assert(fmt.Sprintf("i %v, items[i] == list.Get(i)", i), lg.Equals(item))
	}
	err = mlist2.Append(t.randslice(10))
	t.assert(fmt.Sprintf("err != nil (%v)", err), err != nil)
}
