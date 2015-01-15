package hashtable

import "testing"

import (
	"math/rand"
	"os"
)

import (
	buf "file-structures/block/buffers"
	bs "file-structures/block/byteslice"
	file "file-structures/block/file2"
	"file-structures/linhash"
	"file-structures/linhash/bucket"
	"github.com/timtadh/data-structures/tree/avl"
	. "github.com/timtadh/data-structures/types"
)

func init() {
	if urandom, err := os.Open("/dev/urandom"); err != nil {
		return
	} else {
		seed := make([]byte, 8)
		if _, err := urandom.Read(seed); err == nil {
			rand.Seed(int64(bs.ByteSlice(seed).Int64()))
		}
		urandom.Close()
	}
}

func randstr(length int) String {
	if urandom, err := os.Open("/dev/urandom"); err != nil {
		panic(err)
	} else {
		slice := make([]byte, length)
		if _, err := urandom.Read(slice); err != nil {
			panic(err)
		}
		urandom.Close()
		return String(slice)
	}
	panic("unreachable")
}

func TestMake(t *testing.T) {
	NewHashTable(12)
}

func TestHashable(t *testing.T) {
	a := String("asdf")
	b := String("asdf")
	c := String("csfd")
	if !a.Equals(b) {
		t.Error("a != b")
	}
	if a.Hash() != b.Hash() {
		t.Error("hash(a) != hash(b)")
	}
	if a.Equals(c) {
		t.Error("a == c")
	}
	if a.Hash() != c.Hash() {
		t.Error("hash(a) != hash(c)")
	}
}

func TestPutHasGetRemove(t *testing.T) {

	type record struct {
		key   String
		value String
	}

	ranrec := func() *record {
		return &record{
			String(bs.ByteSlice(randstr(20)).String()),
			String(bs.ByteSlice(randstr(20)).String()),
		}
	}

	test := func(table MapOperable) {
		records := make([]*record, 400)
		for i := range records {
			r := ranrec()
			records[i] = r
			err := table.Put(r.key, String(""))
			if err != nil {
				t.Error(err)
			}
			err = table.Put(r.key, r.value)
			if err != nil {
				t.Error(err)
			}
			if table.Size() != (i + 1) {
				t.Error("size was wrong", table.Size(), i+1)
			}
		}

		for _, r := range records {
			if has := table.Has(r.key); !has {
				t.Error(table, "Missing key")
			}
			if has := table.Has(randstr(12)); has {
				t.Error("Table has extra key")
			}
			if val, err := table.Get(r.key); err != nil {
				t.Error(err)
			} else if !(val.(String)).Equals(r.value) {
				t.Error("wrong value")
			}
		}

		for i, x := range records {
			if val, err := table.Remove(x.key); err != nil {
				t.Error(err)
			} else if !(val.(String)).Equals(x.value) {
				t.Error("wrong value")
			}
			for _, r := range records[i+1:] {
				if has := table.Has(r.key); !has {
					t.Error("Missing key")
				}
				if has := table.Has(randstr(12)); has {
					t.Error("Table has extra key")
				}
				if val, err := table.Get(r.key); err != nil {
					t.Error(err)
				} else if !(val.(String)).Equals(r.value) {
					t.Error("wrong value")
				}
			}
			if table.Size() != (len(records) - (i + 1)) {
				t.Error("size was wrong", table.Size(), (len(records) - (i + 1)))
			}
		}
	}

	test(NewHashTable(64))
	test(NewLinearHash())
}

func TestIterate(t *testing.T) {

	test := func(table Map) {
		t.Logf("%T", table)
		for k, v, next := table.Iterate()(); next != nil; k, v, next = next() {
			t.Errorf("Should never reach here %v %v %v", k, v, next)
		}
		records := make(map[String]String)
		for i := 0; i < 100; i++ {
			k := randstr(8)
			v := randstr(8)
			records[k] = v
			err := table.Put(k, String(""))
			if err != nil {
				t.Error(err)
			}
			err = table.Put(k, v)
			if err != nil {
				t.Error(err)
			}
			if table.Size() != (i + 1) {
				t.Error("size was wrong", table.Size(), i+1)
			}
		}
		newrecs := make(map[String]String)
		for k, v, next := table.Iterate()(); next != nil; k, v, next = next() {
			if v2, has := records[k.(String)]; !has {
				t.Error("bad key in table")
			} else if !v2.Equals(v.(Equatable)) {
				t.Error("values don't agree")
			}
			newrecs[k.(String)] = v.(String)
		}
		if len(records) != len(newrecs) {
			t.Error("iterate missed records")
		}
		for k, v := range records {
			if v2, has := newrecs[k]; !has {
				t.Error("key went missing")
			} else if !v2.Equals(v) {
				t.Error("values don't agree")
			}
		}
	}
	test(NewHashTable(64))
	test(NewLinearHash())
	test(avl.NewAvlTree())
	test(avl.NewImmutableAvlTree())
}

func BenchmarkGoMap(b *testing.B) {
	b.StopTimer()

	type record struct {
		key   string
		value string
	}

	records := make([]*record, 100)

	ranrec := func() *record {
		return &record{string(randstr(20)), string(randstr(20))}
	}

	for i := range records {
		records[i] = ranrec()
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m := make(map[string]string)
		for _, r := range records {
			m[r.key] = r.value
		}
		for _, r := range records {
			delete(m, r.key)
		}
	}
}

func BenchmarkHash(b *testing.B) {
	b.StopTimer()

	type record struct {
		key   String
		value String
	}

	records := make([]*record, 100)

	ranrec := func() *record {
		return &record{randstr(20), randstr(20)}
	}

	for i := range records {
		records[i] = ranrec()
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		t := NewHashTable(128)
		for _, r := range records {
			t.Put(r.key, r.value)
		}
		for _, r := range records {
			t.Remove(r.key)
		}
	}
}

func BenchmarkMLHash(b *testing.B) {
	b.StopTimer()

	type record struct {
		key   String
		value String
	}

	records := make([]*record, 100)

	ranrec := func() *record {
		return &record{randstr(20), randstr(20)}
	}

	for i := range records {
		records[i] = ranrec()
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		t := NewLinearHash()
		for _, r := range records {
			t.Put(r.key, r.value)
		}
		for _, r := range records {
			t.Remove(r.key)
		}
	}
}

func mkfile(path string, size uint64, blksize uint32) (*file.BlockFile, *file.LRUCacheFile) {
	ibf := file.NewBlockFileCustomBlockSize(path, &buf.NoBuffer{}, blksize)
	if err := ibf.Open(); err != nil {
		panic(err)
	}
	f, err := file.NewLRUCacheFile(ibf, size)
	if err != nil {
		panic(err)
	}
	return ibf, f
}

func randstr_safe(length int) string {
	if urandom, err := os.Open("/dev/urandom"); err != nil {
		panic(err)
	} else {
		slice := make(bs.ByteSlice, length)
		if _, err := urandom.Read(slice); err != nil {
			panic(err)
		}
		urandom.Close()
		return slice.String()
	}
	panic("unreachable")
}

func BenchmarkLHash(b *testing.B) {
	b.StopTimer()

	type record struct {
		key   String
		value String
	}

	records := make([]*record, 100)

	ranrec := func() *record {
		return &record{randstr(20), randstr(20)}
	}

	for i := range records {
		records[i] = ranrec()
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {

		hash_file, hf := mkfile("/tmp/linhash_"+string(randstr_safe(8))+".cache",
			1024*1024*512, // 256 MB
			4096)
		store, err := bucket.NewBytesStore(20, 20)
		if err != nil {
			panic(err)
		}
		linhash, err := linhash.NewLinearHash(hf, store)
		if err != nil {
			panic(err)
		}

		for _, r := range records {
			err := linhash.Put([]byte(r.key), []byte(r.value))
			if err != nil {
				panic(err)
			}
		}
		for _, r := range records {
			linhash.Remove([]byte(r.key))
		}

		linhash.Close()
		hash_file.Remove()
	}
}
