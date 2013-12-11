package hashtable

type entry struct {
    key Hashable
    value interface{}
    next *entry
}

type hash struct {
    table []*entry
    size int
}

type String string


func (self String) Equals(other Hashable) bool {
    if o, ok := other.(String); ok {
        return self == o
    } else {
        return false
    }
}

func (self String) Less(other Hashable) bool {
    if o, ok := other.(String); ok {
        return self < o
    } else {
        return false
    }
}

func (self String) Hash() int {
    bytes := []byte(self)
    hash := 0
    for i, c := range bytes {
        hash += (i+1)*int(c)
    }
    return hash
}

func (self *entry) Put(key Hashable, value interface{}) (e *entry, appended bool) {
    if self == nil {
        return &entry{key, value, nil}, true
    }
    if self.key.Equals(key) {
        self.value = value
        return self, false
    } else {
        self.next, appended = self.next.Put(key, value)
        return self, appended
    }
}

func (self *entry) Get(key Hashable) (has bool, value interface{}) {
    if self == nil {
        return false, nil
    } else if self.key.Equals(key) {
        return true, self.value
    } else {
        return self.next.Get(key)
    }
}

func (self *entry) Remove(key Hashable) *entry {
    if self == nil {
        panic(Errors["list-not-found"])
    }
    if self.key.Equals(key) {
        return self.next
    } else {
        self.next = self.next.Remove(key)
        return self
    }
}

func NewHashTable(initial_size int) HashTable {
    return &hash{
        table: make([]*entry, initial_size),
        size: 0,
    }
}

func (self *hash) bucket(key Hashable) int {
    return key.Hash() % len(self.table)
}

func (self *hash) Size() int { return self.size }

func (self *hash) Put(key Hashable, value interface{}) (err error) {
    bucket := self.bucket(key)
    var appended bool
    self.table[bucket], appended = self.table[bucket].Put(key, value)
    if appended {
        self.size += 1
    }
    if self.size * 2 > len(self.table) {
        return self.expand()
    }
    return nil
}

func (self *hash) expand() error {
    table := self.table
    self.table = make([]*entry, len(table)*2)
    self.size = 0
    for _, E := range table {
        for e := E; e != nil; e = e.next {
            if err := self.Put(e.key, e.value); err != nil {
                return err
            }
        }
    }
    return nil
}

func (self *hash) Get(key Hashable) (value interface{}, err error) {
    bucket := self.bucket(key)
    if has, value := self.table[bucket].Get(key); has {
        return value, nil
    } else {
        return nil, Errors["not-found"]
    }
}

func (self *hash) Has(key Hashable) (has bool) {
    has, _ = self.table[self.bucket(key)].Get(key)
    return
}

func (self *hash) Remove(key Hashable) (value interface{}, err error) {
    bucket := self.bucket(key)
    has, value := self.table[bucket].Get(key)
    if !has {
        return nil, Errors["not-found"]
    }
    self.table[bucket] = self.table[bucket].Remove(key)
    self.size -= 1
    return value, nil
}


