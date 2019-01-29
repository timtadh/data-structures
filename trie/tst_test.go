package trie

import (
	"fmt"
	"os"
	"sort"
	"testing"

	crand "crypto/rand"
	"encoding/binary"
	mrand "math/rand"

	trand "github.com/timtadh/data-structures/rand"
	"github.com/timtadh/data-structures/test"
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

func randslice(length int) []byte {
	return test.RandSlice(length)
}

func randstr(length int) types.String {
	return types.String(test.RandStr(length))
}

func has_zero(bytes []byte) bool {
	for _, ch := range bytes {
		if ch == 0 {
			return true
		}
	}
	return false
}

func randslice_nonzero(length int) []byte {
	slice := randslice(length)
	for ; has_zero(slice); slice = randslice(length) {
	}
	return slice
}

func write(name, contents string) {
	file, _ := os.Create(name)
	fmt.Fprintln(file, contents)
	file.Close()
}

type ByteSlices []types.ByteSlice

func (self ByteSlices) Len() int {
	return len(self)
}

func (self ByteSlices) Less(i, j int) bool {
	return self[i].Less(self[j])
}

func (self ByteSlices) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func TestEmptyIter(t *testing.T) {
	table := New()
	for _, _, next := table.Iterate()(); next != nil; _, _, next = next() {
		t.Errorf("should have been empty")
	}

}

func TestIteratorPrefixFindDotty(t *testing.T) {
	items := ByteSlices{
		types.ByteSlice("0:java.io.File;"),
		types.ByteSlice("cat"),
		types.ByteSlice("catty"),
		types.ByteSlice("car"),
		types.ByteSlice("cow"),
		types.ByteSlice("candy"),
		types.ByteSlice("coo"),
		types.ByteSlice("coon"),
		types.ByteSlice("0:java.io.File;1:new,0:java.util.ArrayList;"),
		types.ByteSlice("andy"),
		types.ByteSlice("alex"),
		types.ByteSlice("andrie"),
		types.ByteSlice("alexander"),
		types.ByteSlice("alexi"),
		types.ByteSlice("bob"),
		types.ByteSlice("0:java.io.File;"),
		types.ByteSlice("bobcat"),
		types.ByteSlice("barnaby"),
		types.ByteSlice("baskin"),
		types.ByteSlice("balm"),
	}
	table := New()
	for _, key := range items {
		if err := table.Put(key, nil); err != nil {
			t.Error(table, err)
		}
		if has := table.Has(key); !has {
			t.Error(table, "Missing key")
		}
	}
	write("TestDotty.dot", table.Dotty())
	sort.Sort(items)
	i := 0
	for k, _, next := table.Iterate()(); next != nil; k, _, next = next() {
		if !k.Equals(types.ByteSlice(items[i])) {
			t.Error(string(k.(types.ByteSlice)), "!=", string(items[i]))
		}
		i++
		for i < len(items) && items[i].Equals(items[i-1]) {
			i++
		}
	}
	co_items := ByteSlices{
		types.ByteSlice("coo"),
		types.ByteSlice("coon"),
		types.ByteSlice("cow"),
	}
	i = 0
	for k, _, next := table.PrefixFind([]byte("co"))(); next != nil; k, _, next = next() {
		if !k.Equals(types.ByteSlice(co_items[i])) {
			t.Error(string(k.(types.ByteSlice)), "!=", string(co_items[i]))
		}
		i++
	}
}

func TestComplete4(t *testing.T) {
	items := ByteSlices{
		types.ByteSlice("abaa"),
		types.ByteSlice("abab"),
		types.ByteSlice("abac"),
		types.ByteSlice("abad"),
		types.ByteSlice("abba"),
		types.ByteSlice("abbb"),
		types.ByteSlice("abbc"),
		types.ByteSlice("abbd"),
		types.ByteSlice("abca"),
		types.ByteSlice("abcb"),
		types.ByteSlice("abcc"),
		types.ByteSlice("abcd"),
		types.ByteSlice("abda"),
		types.ByteSlice("abdb"),
		types.ByteSlice("abdc"),
		types.ByteSlice("abdd"),
		types.ByteSlice("aaaa"),
		types.ByteSlice("aaab"),
		types.ByteSlice("aaac"),
		types.ByteSlice("aaad"),
		types.ByteSlice("aaba"),
		types.ByteSlice("aabb"),
		types.ByteSlice("aabc"),
		types.ByteSlice("aabd"),
		types.ByteSlice("aaca"),
		types.ByteSlice("aacb"),
		types.ByteSlice("aacc"),
		types.ByteSlice("aacd"),
		types.ByteSlice("aada"),
		types.ByteSlice("aadb"),
		types.ByteSlice("aadc"),
		types.ByteSlice("aadd"),
		types.ByteSlice("adaa"),
		types.ByteSlice("adab"),
		types.ByteSlice("adac"),
		types.ByteSlice("adad"),
		types.ByteSlice("adba"),
		types.ByteSlice("adbb"),
		types.ByteSlice("adbc"),
		types.ByteSlice("adbd"),
		types.ByteSlice("adca"),
		types.ByteSlice("adcb"),
		types.ByteSlice("adcc"),
		types.ByteSlice("adcd"),
		types.ByteSlice("adda"),
		types.ByteSlice("addb"),
		types.ByteSlice("addc"),
		types.ByteSlice("addd"),
		types.ByteSlice("acaa"),
		types.ByteSlice("acab"),
		types.ByteSlice("acac"),
		types.ByteSlice("acad"),
		types.ByteSlice("acba"),
		types.ByteSlice("acbb"),
		types.ByteSlice("acbc"),
		types.ByteSlice("acbd"),
		types.ByteSlice("acca"),
		types.ByteSlice("accb"),
		types.ByteSlice("accc"),
		types.ByteSlice("accd"),
		types.ByteSlice("acda"),
		types.ByteSlice("acdb"),
		types.ByteSlice("acdc"),
		types.ByteSlice("addd"),
	}
	table := new(TST)
	for _, key := range items {
		if err := table.Put(key, nil); err != nil {
			t.Error(table, err)
		}
		if has := table.Has(key); !has {
			t.Error(table, "Missing key")
		}
	}
	write("TestComplete4.dot", table.Dotty())
	sort.Sort(items)
	i := 0
	for k, _, next := table.Iterate()(); next != nil; k, _, next = next() {
		if !k.Equals(types.ByteSlice(items[i])) {
			t.Error(string(k.(types.ByteSlice)), "!=", string(items[i]))
		}
		i++
		for i+1 < len(items) && items[i].Equals(items[i-1]) {
			i++
		}
	}
}

func TestPutHasGetRemove(t *testing.T) {

	type record struct {
		key   types.ByteSlice
		value types.ByteSlice
	}

	ranrec := func() *record {
		return &record{randslice_nonzero(3), randslice(3)}
	}

	test := func(table *TST) {
		records := make([]*record, 500)
		for i := range records {
			r := ranrec()
			records[i] = r
			err := table.Put(r.key, "")
			if err != nil {
				t.Error(err)
			}
			err = table.Put(r.key, r.value)
			if err != nil {
				t.Error(err)
			}
		}

		for _, r := range records {
			if has := table.Has(r.key); !has {
				t.Error(table, "Missing key")
			}
			if has := table.Has(randslice(12)); has {
				t.Error("Table has extra key")
			}
			if val, err := table.Get(r.key); err != nil {
				t.Error(err)
			} else if !(val.(types.ByteSlice)).Equals(r.value) {
				t.Error("wrong value")
			}
		}
		rs := randslice(12)
		for i, x := range records {
			if val, err := table.Remove(x.key); err != nil {
				t.Error(err)
			} else if !(val.(types.ByteSlice)).Equals(x.value) {
				t.Error("wrong value")
			}
			for _, r := range records[i+1:] {
				if has := table.Has(r.key); !has {
					t.Error("Missing key")
				}
				if has := table.Has(rs); has {
					t.Error("Table has extra key")
				}
				if val, err := table.Get(r.key); err != nil {
					t.Error(err)
				} else if !(val.(types.ByteSlice)).Equals(r.value) {
					t.Error("wrong value")
				}
			}
		}
	}

	test(new(TST))
}

func BenchmarkTST(b *testing.B) {
	b.StopTimer()

	type record struct {
		key   types.ByteSlice
		value types.ByteSlice
	}

	records := make([]*record, 100)

	ranrec := func() *record {
		return &record{randslice(20), randslice(20)}
	}

	for i := range records {
		records[i] = ranrec()
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		t := new(TST)
		for _, r := range records {
			t.Put(r.key, r.value)
		}
		for _, r := range records {
			t.Remove(r.key)
		}
	}
}
