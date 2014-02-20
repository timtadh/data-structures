package trie

import "testing"

import (
    "os"
    "math/rand"
)

import (
    bs "file-structures/block/byteslice"
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

func randslice(length int) []byte {
    if urandom, err := os.Open("/dev/urandom"); err != nil {
        panic(err)
    } else {
        slice := make(bs.ByteSlice, length)
        if _, err := urandom.Read(slice); err != nil {
            panic(err)
        }
        urandom.Close()
        // return append([]byte("b"), slice...)
        return slice
    }
    panic("unreachable")
}

func has_zero(bytes []byte) bool {
    for _,ch := range bytes {
        if ch == 0 {
            return true
        }
    }
    return false
}

func randslice_nonzero(length int) []byte {
    slice := randslice(length)
    for ; has_zero(slice); slice = randslice(length) { }
    return slice
}

func TestPutHasGet(t *testing.T) {

    type record struct {
        key bs.ByteSlice
        value bs.ByteSlice
    }

    ranrec := func() *record {
        return &record{ randslice_nonzero(3), randslice(3) }
    }

    test := func(table *TST) {
        records := make([]*record, 1000)
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
            t.Logf("put %v, %v", r.key, []byte(r.key))
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
            } else if !(val.(bs.ByteSlice)).Eq(r.value) {
                t.Error("wrong value")
            }
        }
    }

    test(new(TST))
}

