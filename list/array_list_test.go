package list

import "testing"

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	mrand "math/rand"
	"runtime/debug"
)

import (
	trand "github.com/timtadh/data-structures/rand"
	"github.com/timtadh/data-structures/types"
)

var rand *mrand.Rand

func init() {
	seed := make([]byte, 8)
	if _, err := crand.Read(seed); err == nil {
		rand = trand.ThreadSafeRand(int64(binary.BigEndian.Uint64(seed)))
	} else {
		panic(err)
	}
}

type T testing.T

func (t *T) assert(msg string, oks ...bool) {
	for _, ok := range oks {
		if !ok {
			t.Log("\n" + string(debug.Stack()))
			t.Error(msg)
			t.Fatal("assert failed")
		}
	}
}

func (t *T) assert_nil(errors ...error) {
	for _, err := range errors {
		if err != nil {
			t.Log("\n" + string(debug.Stack()))
			t.Fatal(err)
		}
	}
}

func (t *T) randslice(length int) types.ByteSlice {
	slice := make([]byte, length)
	if _, err := crand.Read(slice); err != nil {
		t.Fatal(err)
	}
	return types.ByteSlice(slice)
}

func TestAppendGet(x *testing.T) {
	t := (*T)(x)
	SIZE := 100
	list := New(10)
	items := make([]types.ByteSlice, 0, SIZE)
	for i := 0; i < SIZE; i++ {
		item := t.randslice(rand.Intn(10) + 1)
		items = append(items, item)
		t.assert_nil(list.Append(item))
	}
	for i, item := range items {
		lg, err := list.Get(i)
		t.assert_nil(err)
		t.assert(fmt.Sprintf("i %v, items[i] == list.Get(i)", i), lg.Equals(item))
	}
}

func TestAppendGetCopy(x *testing.T) {
	t := (*T)(x)
	SIZE := 100
	list := New(10)
	items := make([]types.ByteSlice, 0, SIZE)
	for i := 0; i < SIZE; i++ {
		item := t.randslice(rand.Intn(10) + 1)
		items = append(items, item)
		t.assert_nil(list.Append(item))
	}
	for i, item := range items {
		lg, err := list.Get(i)
		t.assert_nil(err)
		t.assert(fmt.Sprintf("i %v, items[i] == list.Get(i)", i), lg.Equals(item))
	}
	list2 := list.Copy()
	for i, item := range items {
		lg, err := list2.Get(i)
		t.assert_nil(err)
		t.assert(fmt.Sprintf("i %v, items[i] == list2.Get(i)", i), lg.Equals(item))
	}
}

func TestAppendMarshalUnmarshalGet(x *testing.T) {
	t := (*T)(x)
	SIZE := 100
	list := New(10)
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
}

func TestInsertGetSet(x *testing.T) {
	t := (*T)(x)
	SIZE := 100
	list := New(10)
	items := make([]types.Int, 0, SIZE)
	order := make([]int, 0, SIZE)
	{
		item := types.Int(0)
		items = append(items, item)
		order = append(order, 3)
		t.assert_nil(list.Insert(0, item))
		// [0]
	}
	{
		item := types.Int(1)
		items = append(items, item)
		order = append(order, 0)
		t.assert_nil(list.Insert(0, item))
		// [1, 0]
	}
	{
		item := types.Int(2)
		items = append(items, item)
		order = append(order, 1)
		t.assert_nil(list.Insert(1, item))
		// [1, 2, 0]
	}
	{
		item := types.Int(3)
		items = append(items, item)
		order = append(order, 4)
		t.assert_nil(list.Insert(3, item))
		// [1, 2, 0, 3]
	}
	{
		item := types.Int(4)
		items = append(items, item)
		order = append(order, 2)
		t.assert_nil(list.Insert(2, item))
		// [1, 2, 4, 0, 3]
	}
	for j, i := range order {
		item := items[j]
		lg, err := list.Get(i)
		t.assert_nil(err)
		t.assert(fmt.Sprintf("j %v, i %v, items[i], %v != %v, list.Get(i)", j, i, item, lg), lg.Equals(item))
	}
	_, err := list.Get(-1)
	t.assert(fmt.Sprintf("err != nil, err == %v", err), err != nil)
	_, err = list.Get(5)
	t.assert(fmt.Sprintf("err != nil, err == %v", err), err != nil)
	err = list.Set(5, nil)
	t.assert(fmt.Sprintf("err != nil, err == %v", err), err != nil)
	err = list.Insert(500, nil)
	t.assert(fmt.Sprintf("err != nil, err == %v", err), err != nil)
	err = list.Remove(500)
	t.assert(fmt.Sprintf("err != nil, err == %v", err), err != nil)
	for i := 0; i < list.Size(); i++ {
		list.Set(i, types.Int(i))
	}
	for i := 0; i < list.Size(); i++ {
		lg, err := list.Get(i)
		t.assert_nil(err)
		t.assert(fmt.Sprintf("i %v, items[i] == list.Get(i)", i), lg.Equals(types.Int(i)))
	}
	list.Clear()
	t.assert(fmt.Sprintf("size != 0, %v", list.Size()), list.Size() == 0)
}

func TestAppendPop(x *testing.T) {
	t := (*T)(x)
	SIZE := 100
	list := New(10)
	items := make([]types.ByteSlice, 0, SIZE)
	for i := 0; i < SIZE; i++ {
		item := t.randslice(rand.Intn(10) + 1)
		items = append(items, item)
		t.assert_nil(list.Append(item))
	}
	for i := SIZE - 1; i >= 0; i-- {
		item, err := list.Pop()
		t.assert_nil(err)
		t.assert(fmt.Sprintf("i %v, items[i] == list.Pop()", i), item.Equals(items[i]))
	}
}

func TestExtend(x *testing.T) {
	t := (*T)(x)
	SIZE := 100
	a := New(10)
	b := New(10)
	items := make([]types.ByteSlice, 0, SIZE)
	for i := 0; i < SIZE; i++ {
		item := t.randslice(rand.Intn(10) + 1)
		items = append(items, item)
		if i < SIZE/2 {
			t.assert_nil(a.Append(item))
		} else {
			t.assert_nil(b.Append(item))
		}
	}
	t.assert_nil(a.Extend(b.Items()))
	for i := SIZE - 1; i >= 0; i-- {
		item, err := a.Pop()
		t.assert_nil(err)
		t.assert(fmt.Sprintf("i %v, items[i] == list.Pop()", i), item.Equals(items[i]))
	}
}

func TestLess(x *testing.T) {
	t := (*T)(x)
	a := FromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(3)})
	b := FromSlice([]types.Hashable{types.Int(3), types.Int(2), types.Int(3)})
	c := FromSlice([]types.Hashable{types.Int(1), types.Int(0), types.Int(3)})
	d := FromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(3)})
	e := FromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(4)})
	small := FromSlice([]types.Hashable{types.Int(2), types.Int(4)})
	big := FromSlice([]types.Hashable{types.Int(0), types.Int(1), types.Int(2), types.Int(4)})
	t.assert("a < b", a.Less(b))
	t.assert("c < a", c.Less(a))
	t.assert("c < b", c.Less(b))
	t.assert("a !< d", a.Less(d) == false)
	t.assert("a !< e", e.Less(a) == false)
	t.assert("a < big", a.Less(big))
	t.assert("small < a", small.Less(a))
	t.assert("a !< small", a.Less(small) == false)
}

func TestEqualsHash(x *testing.T) {
	t := (*T)(x)
	a := FromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(3)})
	b := FromSlice([]types.Hashable{types.Int(3), types.Int(2), types.Int(3)})
	c := FromSlice([]types.Hashable{types.Int(1), types.Int(0), types.Int(3)})
	d := FromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(3)})
	small := FromSlice([]types.Hashable{types.Int(2), types.Int(4)})
	empty := FromSlice([]types.Hashable{})
	t.assert("a != b", !a.Equals(b))
	t.assert("c != a", !c.Equals(a))
	t.assert("c != b", !c.Equals(b))
	t.assert("a == d", a.Equals(d))
	t.assert("c != small", !c.Equals(small))
	t.assert("a.Hash() != b.Hash()", a.Hash() != b.Hash())
	t.assert("c.Hash() != b.Hash()", c.Hash() != b.Hash())
	t.assert("a.Hash() != d.Hash()", a.Hash() == d.Hash())
	t.assert("d.Hash() != b.Hash()", d.Hash() != b.Hash())
	t.assert("d.Hash() != empty.Hash()", d.Hash() != empty.Hash())
}
