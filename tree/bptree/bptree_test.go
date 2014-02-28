package tree

import "testing"

import (
  "github.com/timtadh/data-structures/types"
)

func Test_balance_leaf_nodes(t *testing.T) {
    a := NewLeaf(7)
    b := NewLeaf(7)
    if err := a.put_kv(types.Int(1), 1); err != nil {
        t.Error(err)
    }
    if err := a.put_kv(types.Int(2), 2); err != nil {
        t.Error(err)
    }
    if err := a.put_kv(types.Int(3), 3); err != nil {
        t.Error(err)
    }
    if err := a.put_kv(types.Int(4), 4); err != nil {
        t.Error(err)
    }
    if err := a.put_kv(types.Int(5), 5); err != nil {
        t.Error(err)
    }
    if err := a.put_kv(types.Int(6), 6); err != nil {
        t.Error(err)
    }
    if err := a.put_kv(types.Int(7), 7); err != nil {
        t.Error(err)
    }
    balance_nodes(a, b)
    for i, k := range a.keys {
        if int(k.(types.Int)) != i+1 {
            t.Errorf("k != %d", i+1)
        }
    }
    for i, k := range b.keys {
        if int(k.(types.Int)) != 3+i+1 {
            t.Errorf("k != %d", 3+i+1)
        }
    }
    for i, v := range a.values {
        if v.(int) != i+1 {
            t.Errorf("k != %d", i+1)
        }
    }
    for i, v := range b.values {
        if v.(int) != 3+i+1 {
            t.Errorf("v != %d", 3+i+1)
        }
    }
    t.Log(a)
    t.Log(b)
}

func Test_balance_internal_nodes(t *testing.T) {
    a := NewInternal(6)
    b := NewInternal(6)
    if err := a.put_kp(types.Int(1), nil); err != nil {
        t.Error(err)
    }
    if err := a.put_kp(types.Int(2), nil); err != nil {
        t.Error(err)
    }
    if err := a.put_kp(types.Int(3), nil); err != nil {
        t.Error(err)
    }
    if err := a.put_kp(types.Int(4), nil); err != nil {
        t.Error(err)
    }
    if err := a.put_kp(types.Int(5), nil); err != nil {
        t.Error(err)
    }
    if err := a.put_kp(types.Int(6), nil); err != nil {
        t.Error(err)
    }
    balance_nodes(a, b)
    for i, k := range a.keys {
        if int(k.(types.Int)) != i+1 {
            t.Errorf("k != %d", i+1)
        }
    }
    for i, k := range b.keys {
        if int(k.(types.Int)) != 3+i+1 {
            t.Errorf("k != %d", 3+i+1)
        }
    }
    t.Log(a)
    t.Log(b)
}

func Test_put_no_root_split(t *testing.T) {
    a := NewLeaf(2)
    if err := a.put_kv(types.Int(1), 1); err != nil { t.Error(err) }
    p, err := a.put(types.Int(1), 2)
    if err != nil {
        t.Error(err)
    } else {
        if p != a {
            t.Errorf("p != a")
        }
        if !p.Has(types.Int(1)) {
            t.Error("p didn't have the right keys", p)
        }
    }
    p, err = a.put(types.Int(1), 3)
    if err != nil {
        t.Error(err)
    } else {
        if p != a {
            t.Errorf("p != a")
        }
        if !p.Has(types.Int(1)) {
            t.Error("p didn't have the right keys", p)
        }
        if p.next == nil {
            t.Error("p.next should not be nil")
        }
        t.Log(p)
        t.Log(p.next)
    }
}

func Test_put_root_split(t *testing.T) {
    a := NewLeaf(2)
    p, err := a.put(types.Int(1), 1)
    if err != nil {
        t.Error(err)
    } else {
        if p != a {
            t.Errorf("p != a")
        }
        if !p.Has(types.Int(1)) {
            t.Error("p didn't have the right keys", p)
        }
    }
    p, err = a.put(types.Int(3), 3)
    if err != nil {
        t.Error(err)
    } else {
        if p != a {
            t.Errorf("p != a")
        }
        if !p.Has(types.Int(1)) || !p.Has(types.Int(3)) {
            t.Error("p didn't have the right keys", p)
        }
    }
    p, err = a.put(types.Int(2), 2)
    if err != nil {
        t.Error(err)
    } else {
        if p == a {
            t.Errorf("p == a")
        }
        if !p.Has(types.Int(1)) || !p.Has(types.Int(3)) {
            t.Error("p didn't have the right keys", p)
        }
        if len(p.pointers) != 2 {
            t.Error("p didn't have right number of pointers", p)
        }
        if !p.pointers[0].Has(types.Int(1)) || !p.pointers[0].Has(types.Int(2)) {
            t.Error("p.pointers[0] didn't have the right keys", p.pointers[0])
        }
        if !p.pointers[1].Has(types.Int(3)) {
            t.Error("p.pointers[1] didn't have the right keys", p.pointers[1])
        }
        t.Log(p)
        t.Log(p.pointers[0])
        t.Log(p.pointers[1])
    }
}

func Test_internal_insert_no_split(t *testing.T) {
    a := NewInternal(3)
    leaf := NewLeaf(1)
    if err := leaf.put_kv(types.Int(1), 1); err != nil { t.Error(err) }
    if err := a.put_kp(types.Int(1), leaf); err != nil { t.Error(err) }
    if err := a.put_kp(types.Int(5), nil); err != nil { t.Error(err) }
    p, q, err := a.internal_insert(types.Int(2), nil)
    if err != nil {
        t.Error(err)
    } else {
        if p != a {
            t.Errorf("p != a")
        }
        if q != nil {
            t.Errorf("q != nil")
        }
        if !p.Has(types.Int(1)) || !p.Has(types.Int(2)) || !p.Has(types.Int(5)) {
            t.Error("p didn't have the right keys", p)
        }
    }
}

func Test_internal_insert_split_less(t *testing.T) {
    a := NewInternal(3)
    leaf := NewLeaf(1)
    if err := leaf.put_kv(types.Int(1), 1); err != nil { t.Error(err) }
    if err := a.put_kp(types.Int(1), leaf); err != nil { t.Error(err) }
    if err := a.put_kp(types.Int(3), nil); err != nil { t.Error(err) }
    if err := a.put_kp(types.Int(5), nil); err != nil { t.Error(err) }
    p, q, err := a.internal_insert(types.Int(2), nil)
    if err != nil {
        t.Error(err)
    } else {
        if p != a {
            t.Errorf("p != a")
        }
        if q == nil {
            t.Errorf("q == nil")
        }
        if !p.Has(types.Int(1)) || !p.Has(types.Int(2)) {
            t.Error("p didn't have the right keys", p)
        }
        if !q.Has(types.Int(3)) || !q.Has(types.Int(5)) {
            t.Error("q didn't have the right keys", q)
        }
    }
}

func Test_internal_split_less(t *testing.T) {
    a := NewInternal(3)
    if err := a.put_kp(types.Int(1), nil); err != nil { t.Error(err) }
    if err := a.put_kp(types.Int(3), nil); err != nil { t.Error(err) }
    if err := a.put_kp(types.Int(5), nil); err != nil { t.Error(err) }
    p, q, err := a.internal_split(types.Int(2), nil)
    if err != nil {
        t.Error(err)
    } else {
        if p != a {
            t.Errorf("p != a")
        }
        if q == nil {
            t.Errorf("q == nil")
        }
        if !p.Has(types.Int(1)) || !p.Has(types.Int(2)) {
            t.Error("p didn't have the right keys", p)
        }
        if !q.Has(types.Int(3)) || !q.Has(types.Int(5)) {
            t.Error("q didn't have the right keys", q)
        }
    }
}

func Test_internal_split_equal(t *testing.T) {
    a := NewInternal(3)
    if err := a.put_kp(types.Int(1), nil); err != nil { t.Error(err) }
    if err := a.put_kp(types.Int(3), nil); err != nil { t.Error(err) }
    if err := a.put_kp(types.Int(5), nil); err != nil { t.Error(err) }
    p, q, err := a.internal_split(types.Int(3), nil)
    if err == nil {
        t.Error("split succeeded should have failed", p, q)
    }
}

func Test_internal_split_greater(t *testing.T) {
    a := NewInternal(3)
    if err := a.put_kp(types.Int(1), nil); err != nil { t.Error(err) }
    if err := a.put_kp(types.Int(3), nil); err != nil { t.Error(err) }
    if err := a.put_kp(types.Int(5), nil); err != nil { t.Error(err) }
    p, q, err := a.internal_split(types.Int(4), nil)
    if err != nil {
        t.Error(err)
    } else {
        if p != a {
            t.Errorf("p != a")
        }
        if q == nil {
            t.Errorf("q == nil")
        }
        if !p.Has(types.Int(1)) {
            t.Error("p didn't have the right keys", p)
        }
        if !q.Has(types.Int(3)) || !q.Has(types.Int(4)) || !q.Has(types.Int(5)) {
            t.Error("q didn't have the right keys", q)
        }
    }
}

func Test_leaf_insert_no_split(t *testing.T) {
    a := NewLeaf(3)
    insert_linked_list_node(a, nil, nil)
    if err := a.put_kv(types.Int(1), 1); err != nil { t.Error(err) }
    if err := a.put_kv(types.Int(3), 3); err != nil { t.Error(err) }
    p, q, err := a.leaf_insert(types.Int(2), 2)
    if err != nil {
        t.Error(err)
    } else {
        if p != a {
            t.Errorf("p != a")
        }
        if q != nil {
            t.Errorf("q != nil")
        }
        if !p.Has(types.Int(1)) || !p.Has(types.Int(2)) || !p.Has(types.Int(3)) {
            t.Error("p didn't have the right keys", p)
        }
    }
}

// tests the defer to split logic
func Test_leaf_insert_split_less(t *testing.T) {
    a := NewLeaf(3)
    insert_linked_list_node(a, nil, nil)
    if err := a.put_kv(types.Int(1), 1); err != nil { t.Error(err) }
    if err := a.put_kv(types.Int(3), 3); err != nil { t.Error(err) }
    if err := a.put_kv(types.Int(5), 5); err != nil { t.Error(err) }
    p, q, err := a.leaf_insert(types.Int(2), 2)
    if err != nil {
        t.Error(err)
    } else {
        if p != a {
            t.Errorf("p != a")
        }
        if q == nil {
            t.Errorf("q == nil")
        }
        if !p.Has(types.Int(1)) || !p.Has(types.Int(2)) {
            t.Error("p didn't have the right keys", p)
        }
        if !q.Has(types.Int(3)) || !q.Has(types.Int(5)) {
            t.Error("q didn't have the right keys", q)
        }
    }
}

func Test_leaf_split_less(t *testing.T) {
    a := NewLeaf(3)
    insert_linked_list_node(a, nil, nil)
    if err := a.put_kv(types.Int(1), 1); err != nil { t.Error(err) }
    if err := a.put_kv(types.Int(3), 3); err != nil { t.Error(err) }
    if err := a.put_kv(types.Int(5), 5); err != nil { t.Error(err) }
    p, q, err := a.leaf_split(types.Int(2), 2)
    if err != nil {
        t.Error(err)
    } else {
        if p != a {
            t.Errorf("p != a")
        }
        if q == nil {
            t.Errorf("q == nil")
        }
        if !p.Has(types.Int(1)) || !p.Has(types.Int(2)) {
            t.Error("p didn't have the right keys", p)
        }
        if !q.Has(types.Int(3)) || !q.Has(types.Int(5)) {
            t.Error("q didn't have the right keys", q)
        }
    }
}

func Test_leaf_split_equal(t *testing.T) {
    a := NewLeaf(3)
    insert_linked_list_node(a, nil, nil)
    if err := a.put_kv(types.Int(1), 1); err != nil { t.Error(err) }
    if err := a.put_kv(types.Int(3), 3); err != nil { t.Error(err) }
    if err := a.put_kv(types.Int(5), 5); err != nil { t.Error(err) }
    p, q, err := a.leaf_split(types.Int(3), 2)
    if err != nil {
        t.Error(err)
    } else {
        if p != a {
            t.Errorf("p != a")
        }
        if q == nil {
            t.Errorf("q == nil")
        }
        if !p.Has(types.Int(1)) {
            t.Error("p didn't have the right keys", p)
        }
        if !q.Has(types.Int(3)) || q.Count(types.Int(3)) != 2 || !q.Has(types.Int(5)) {
            t.Error("q didn't have the right keys", q, q.Count(types.Int(3)))
        }
    }
}

func Test_leaf_split_greater(t *testing.T) {
    a := NewLeaf(3)
    insert_linked_list_node(a, nil, nil)
    if err := a.put_kv(types.Int(1), 1); err != nil { t.Error(err) }
    if err := a.put_kv(types.Int(3), 3); err != nil { t.Error(err) }
    if err := a.put_kv(types.Int(5), 5); err != nil { t.Error(err) }
    p, q, err := a.leaf_split(types.Int(4), 2)
    if err != nil {
        t.Error(err)
    } else {
        if p != a {
            t.Errorf("p != a")
        }
        if q == nil {
            t.Errorf("q == nil")
        }
        if !p.Has(types.Int(1)) {
            t.Error("p didn't have the right keys", p)
        }
        if !q.Has(types.Int(3)) || !q.Has(types.Int(4)) || !q.Has(types.Int(5)) {
            t.Error("q didn't have the right keys", q)
        }
    }
}

// tests the defer logic
func Test_pure_leaf_insert_split_less(t *testing.T) {
    a := NewLeaf(2)
    insert_linked_list_node(a, nil, nil)
    b := NewLeaf(2)
    insert_linked_list_node(b, a, nil)
    c := NewLeaf(2)
    insert_linked_list_node(c, b, nil)
    d := NewLeaf(2)
    insert_linked_list_node(d, c, nil)
    if err := a.put_kv(types.Int(3), 1); err != nil { t.Error(err) }
    if err := a.put_kv(types.Int(3), 2); err != nil { t.Error(err) }
    if err := b.put_kv(types.Int(3), 3); err != nil { t.Error(err) }
    if err := b.put_kv(types.Int(3), 4); err != nil { t.Error(err) }
    if err := c.put_kv(types.Int(3), 5); err != nil { t.Error(err) }
    if err := c.put_kv(types.Int(3), 6); err != nil { t.Error(err) }
    if err := d.put_kv(types.Int(4), 6); err != nil { t.Error(err) }
    p, q, err := a.leaf_insert(types.Int(2), 1)
    if err != nil {
        t.Error(err)
    } else {
        if q != a {
            t.Errorf("q != a")
        }
        if p == nil || len(p.keys) != 1 || !p.keys[0].Equals(types.Int(2)) {
            t.Errorf("p did not contain the right key")
        }
        if p.prev != nil { t.Errorf("expected p.prev == nil") }
        if p.next != a { t.Errorf("expected p.next == a") }
        if a.prev != p { t.Errorf("expected a.prev == p") }
        if a.next != b { t.Errorf("expected a.next == b") }
        if b.prev != a { t.Errorf("expected b.prev == a") }
        if b.next != c { t.Errorf("expected b.next == c") }
        if c.prev != b { t.Errorf("expected c.prev == b") }
        if c.next != d { t.Errorf("expected c.next == d") }
        if d.prev != c { t.Errorf("expected d.prev == c") }
        if d.next != nil { t.Errorf("expected d.next == nil") }
    }
}

func Test_pure_leaf_split_less(t *testing.T) {
    a := NewLeaf(2)
    insert_linked_list_node(a, nil, nil)
    b := NewLeaf(2)
    insert_linked_list_node(b, a, nil)
    c := NewLeaf(2)
    insert_linked_list_node(c, b, nil)
    d := NewLeaf(2)
    insert_linked_list_node(d, c, nil)
    if err := a.put_kv(types.Int(3), 1); err != nil { t.Error(err) }
    if err := a.put_kv(types.Int(3), 2); err != nil { t.Error(err) }
    if err := b.put_kv(types.Int(3), 3); err != nil { t.Error(err) }
    if err := b.put_kv(types.Int(3), 4); err != nil { t.Error(err) }
    if err := c.put_kv(types.Int(3), 5); err != nil { t.Error(err) }
    if err := c.put_kv(types.Int(3), 6); err != nil { t.Error(err) }
    if err := d.put_kv(types.Int(4), 6); err != nil { t.Error(err) }
    p, q, err := a.pure_leaf_split(types.Int(2), 1)
    if err != nil {
        t.Error(err)
    } else {
        if q != a {
            t.Errorf("q != a")
        }
        if p == nil || len(p.keys) != 1 || !p.keys[0].Equals(types.Int(2)) {
            t.Errorf("p did not contain the right key")
        }
        if p.prev != nil { t.Errorf("expected p.prev == nil") }
        if p.next != a { t.Errorf("expected p.next == a") }
        if a.prev != p { t.Errorf("expected a.prev == p") }
        if a.next != b { t.Errorf("expected a.next == b") }
        if b.prev != a { t.Errorf("expected b.prev == a") }
        if b.next != c { t.Errorf("expected b.next == c") }
        if c.prev != b { t.Errorf("expected c.prev == b") }
        if c.next != d { t.Errorf("expected c.next == d") }
        if d.prev != c { t.Errorf("expected d.prev == c") }
        if d.next != nil { t.Errorf("expected d.next == nil") }
    }
}

func Test_pure_leaf_split_equal(t *testing.T) {
    a := NewLeaf(2)
    insert_linked_list_node(a, nil, nil)
    b := NewLeaf(2)
    insert_linked_list_node(b, a, nil)
    c := NewLeaf(2)
    insert_linked_list_node(c, b, nil)
    d := NewLeaf(2)
    insert_linked_list_node(d, c, nil)
    if err := a.put_kv(types.Int(3), 1); err != nil { t.Error(err) }
    if err := a.put_kv(types.Int(3), 2); err != nil { t.Error(err) }
    if err := b.put_kv(types.Int(3), 3); err != nil { t.Error(err) }
    if err := b.put_kv(types.Int(3), 4); err != nil { t.Error(err) }
    if err := c.put_kv(types.Int(3), 5); err != nil { t.Error(err) }
    if err := d.put_kv(types.Int(4), 6); err != nil { t.Error(err) }
    p, q, err := a.pure_leaf_split(types.Int(3), 1)
    if err != nil {
        t.Error(err)
    } else {
        if p != a {
            t.Errorf("p != a")
        }
        if q != nil {
            t.Errorf("q != nil")
        }
        if a.prev != nil { t.Errorf("expected a.prev == nil") }
        if a.next != b { t.Errorf("expected a.next == b") }
        if b.prev != a { t.Errorf("expected b.prev == a") }
        if b.next != c { t.Errorf("expected b.next == c") }
        if c.prev != b { t.Errorf("expected c.prev == b") }
        if c.next != d { t.Errorf("expected c.next == d") }
        if d.prev != c { t.Errorf("expected d.prev == c") }
        if d.next != nil { t.Errorf("expected d.next == nil") }
    }
}

func Test_pure_leaf_split_greater(t *testing.T) {
    a := NewLeaf(2)
    insert_linked_list_node(a, nil, nil)
    b := NewLeaf(2)
    insert_linked_list_node(b, a, nil)
    c := NewLeaf(2)
    insert_linked_list_node(c, b, nil)
    d := NewLeaf(2)
    insert_linked_list_node(d, c, nil)
    if err := a.put_kv(types.Int(3), 1); err != nil { t.Error(err) }
    if err := a.put_kv(types.Int(3), 2); err != nil { t.Error(err) }
    if err := b.put_kv(types.Int(3), 3); err != nil { t.Error(err) }
    if err := b.put_kv(types.Int(3), 4); err != nil { t.Error(err) }
    if err := c.put_kv(types.Int(3), 5); err != nil { t.Error(err) }
    if err := d.put_kv(types.Int(5), 6); err != nil { t.Error(err) }
    p, q, err := a.pure_leaf_split(types.Int(4), 1)
    if err != nil {
        t.Error(err)
    } else {
        if p != a {
            t.Errorf("p != a")
        }
        if q == nil || len(q.keys) != 1 || !q.keys[0].Equals(types.Int(4)) {
            t.Errorf("q != nil")
        }
        if a.prev != nil { t.Errorf("expected a.prev == nil") }
        if a.next != b { t.Errorf("expected a.next == b") }
        if b.prev != a { t.Errorf("expected b.prev == a") }
        if b.next != c { t.Errorf("expected b.next == c") }
        if c.prev != b { t.Errorf("expected c.prev == b") }
        if c.next != q { t.Errorf("expected c.next == q") }
        if q.prev != c { t.Errorf("expected q.prev == c") }
        if q.next != d { t.Errorf("expected q.next == d") }
        if d.prev != q { t.Errorf("expected d.prev == q") }
        if d.next != nil { t.Errorf("expected d.next == nil") }
    }
}

func Test_find_end_of_pure_run(t *testing.T) {
    a := NewLeaf(2)
    insert_linked_list_node(a, nil, nil)
    b := NewLeaf(2)
    insert_linked_list_node(b, a, nil)
    c := NewLeaf(2)
    insert_linked_list_node(c, b, nil)
    d := NewLeaf(2)
    insert_linked_list_node(d, c, nil)
    if err := a.put_kv(types.Int(3), 1); err != nil { t.Error(err) }
    if err := a.put_kv(types.Int(3), 2); err != nil { t.Error(err) }
    if err := b.put_kv(types.Int(3), 3); err != nil { t.Error(err) }
    if err := b.put_kv(types.Int(3), 4); err != nil { t.Error(err) }
    if err := c.put_kv(types.Int(3), 5); err != nil { t.Error(err) }
    if err := c.put_kv(types.Int(3), 6); err != nil { t.Error(err) }
    if err := d.put_kv(types.Int(4), 6); err != nil { t.Error(err) }
    e := a.find_end_of_pure_run()
    if e != c {
        t.Errorf("end of run should have been block c %v %v", e, c)
    }
}

func Test_insert_linked_list_node(t *testing.T) {
    a := NewLeaf(1)
    insert_linked_list_node(a, nil, nil)
    b := NewLeaf(2)
    insert_linked_list_node(b, a, nil)
    c := NewLeaf(3)
    insert_linked_list_node(c, b, nil)
    d := NewLeaf(4)
    insert_linked_list_node(d, a, b)
    if a.prev != nil { t.Errorf("expected a.prev == nil") }
    if a.next != d { t.Errorf("expected a.next == d") }
    if d.prev != a { t.Errorf("expected d.prev == a") }
    if d.next != b { t.Errorf("expected d.next == b") }
    if b.prev != d { t.Errorf("expected b.prev == d") }
    if b.next != c { t.Errorf("expected b.next == c") }
    if c.prev != b { t.Errorf("expected c.prev == b") }
    if c.next != nil { t.Errorf("expected c.next == nil") }
}

func Test_remove_linked_list_node(t *testing.T) {
    a := NewLeaf(1)
    insert_linked_list_node(a, nil, nil)
    b := NewLeaf(2)
    insert_linked_list_node(b, a, nil)
    c := NewLeaf(3)
    insert_linked_list_node(c, b, nil)
    d := NewLeaf(4)
    insert_linked_list_node(d, a, b)
    if a.prev != nil { t.Errorf("expected a.prev == nil") }
    if a.next != d { t.Errorf("expected a.next == d") }
    if d.prev != a { t.Errorf("expected d.prev == a") }
    if d.next != b { t.Errorf("expected d.next == b") }
    if b.prev != d { t.Errorf("expected b.prev == d") }
    if b.next != c { t.Errorf("expected b.next == c") }
    if c.prev != b { t.Errorf("expected c.prev == b") }
    if c.next != nil { t.Errorf("expected c.next == nil") }
    remove_linked_list_node(d)
    if a.prev != nil { t.Errorf("expected a.prev == nil") }
    if a.next != b { t.Errorf("expected a.next == b") }
    if b.prev != a { t.Errorf("expected b.prev == a") }
    if b.next != c { t.Errorf("expected b.next == c") }
    if c.prev != b { t.Errorf("expected c.prev == b") }
    if c.next != nil { t.Errorf("expected c.next == nil") }
    remove_linked_list_node(a)
    if b.prev != nil { t.Errorf("expected b.prev == nil") }
    if b.next != c { t.Errorf("expected b.next == c") }
    if c.prev != b { t.Errorf("expected c.prev == b") }
    if c.next != nil { t.Errorf("expected c.next == nil") }
    remove_linked_list_node(c)
    if b.prev != nil { t.Errorf("expected b.prev == nil") }
    if b.next != nil { t.Errorf("expected b.next == nil") }
    remove_linked_list_node(b)
}
