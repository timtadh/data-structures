package hashtable

import "testing"

import (
    "os"
    "math/rand"
)

import (
    bs "file-structures/block/byteslice"
    file "file-structures/block/file2"
    buf "file-structures/block/buffers"
    "file-structures/linhash"
    "file-structures/linhash/bucket"
)


func init() {
    if urandom, err := os.Open("/dev/urandom"); err != nil {
        return
    } else {
        seed := make([]byte, 8)
        if _, err := urandom.Read(seed); err == nil {
            rand.Seed(int64(bs.ByteSlice(seed).Int64()))
        }
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
        key String
        value String
    }


    ranrec := func() *record {
        return &record{
          String(bs.ByteSlice(randstr(20)).String()),
          String(bs.ByteSlice(randstr(20)).String()),
        }
    }

    test := func(table HashTable) {
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
            if table.Size() != (i+1) {
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
            if table.Size() != (len(records) - (i+1)) {
                t.Error("size was wrong", table.Size(), (len(records) - (i+1)))
            }
        }
    }

    test(NewHashTable(64))
    test(NewLinearHash())
}

func TestPutHasGetRemoveBucket(t *testing.T) {

    type record struct {
        key String
        value String
    }

    records := make([]*record, 400)
    var tree *bst
    var err error
    var val interface{}
    var updated bool

    ranrec := func() *record {
        return &record{ randstr(20), randstr(20) }
    }

    for i := range records {
        r := ranrec()
        records[i] = r
        tree, updated = tree.Put(r.key, String(""))
        if updated {
            t.Error("should have not been updated")
        }
        tree, updated = tree.Put(r.key, r.value)
        if !updated {
            t.Error("should have been updated")
        }
        if tree.Size() != (i+1) {
            t.Error("size was wrong", tree.Size(), i+1)
        }
    }

    for _, r := range records {
        if has := tree.Has(r.key); !has {
            t.Error("Missing key")
        }
        if has := tree.Has(randstr(12)); has {
            t.Error("Table has extra key")
        }
        if val, err := tree.Get(r.key); err != nil {
            t.Error(err, bs.ByteSlice(val.(String)), bs.ByteSlice(r.value))
        } else if !(val.(String)).Equals(r.value) {
            t.Error("wrong value")
        }
    }

    for i, x := range records {
        if tree, val, err = tree.Remove(x.key); err != nil {
            t.Error(err)
        } else if !(val.(String)).Equals(x.value) {
            t.Error("wrong value")
        }
        for _, r := range records[i+1:] {
            if has := tree.Has(r.key); !has {
                t.Error("Missing key")
            }
            if has := tree.Has(randstr(12)); has {
                t.Error("Table has extra key")
            }
            if val, err := tree.Get(r.key); err != nil {
                t.Error(err)
            } else if !(val.(String)).Equals(r.value) {
                t.Error("wrong value")
            }
        }
        if tree.Size() != (len(records) - (i+1)) {
            t.Error("size was wrong", tree.Size(), (len(records) - (i+1)))
        }
    }
}


func BenchmarkGoMap(b *testing.B) {
    b.StopTimer()

    type record struct {
        key string
        value string
    }

    records := make([]*record, 100)

    ranrec := func() *record {
        return &record{ string(randstr(20)), string(randstr(20)) }
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
        key String
        value String
    }

    records := make([]*record, 100)

    ranrec := func() *record {
        return &record{ randstr(20), randstr(20) }
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
        key String
        value String
    }

    records := make([]*record, 100)

    ranrec := func() *record {
        return &record{ randstr(20), randstr(20) }
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
        key String
        value String
    }

    records := make([]*record, 100)

    ranrec := func() *record {
        return &record{ randstr(20), randstr(20) }
    }

    for i := range records {
        records[i] = ranrec()
    }
    b.StartTimer()
    for i := 0; i < b.N; i++ {

        hash_file, hf := mkfile("/tmp/linhash_" + string(randstr_safe(8)) + ".cache",
                     1024*1024*512, // 256 MB
                     4096)
        store, err := bucket.NewBytesStore(20, 20)
        if err != nil { panic(err) }
        linhash, err := linhash.NewLinearHash(hf, store)
        if err != nil { panic(err) }

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


