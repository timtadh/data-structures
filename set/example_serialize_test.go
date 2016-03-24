package set_test

import (
	"log"
)

import (
	"github.com/timtadh/data-structures/list"
	"github.com/timtadh/data-structures/set"
	"github.com/timtadh/data-structures/types"
)

func makeSet() *set.SortedSet {
	return set.FromSlice([]types.Hashable{types.Int(1), types.Int(-1), types.Int(3)})
}

func serialize(s *set.SortedSet) ([]byte, error) {
	marshal, unmarshal := types.IntMarshals()
	m := set.NewMSortedSet(s, marshal, unmarshal)
	return m.MarshalBinary()
}

func deserialize(bytes []byte) (*set.SortedSet, error) {
	marshal, unmarshal := types.IntMarshals()
	m := &set.MSortedSet{MSorted: list.MSorted{MList: list.MList{MarshalItem: marshal, UnmarshalItem: unmarshal}}}
	err := m.UnmarshalBinary(bytes)
	if err != nil {
		return nil, err
	}
	return m.SortedSet(), nil
}

func Example_serialize() {
	a := makeSet()
	b := makeSet()
	if !a.Equals(b) {
		log.Panic("a was not equal to b")
	}
	bytes, err := serialize(a)
	if err != nil {
		log.Panic(err)
	}
	log.Println(bytes)
	c, err := deserialize(bytes)
	if err != nil {
		log.Panic(err)
	}
	if !c.Equals(b) {
		log.Panic("c was not equal to b")
	}
	log.Println("success")
}
