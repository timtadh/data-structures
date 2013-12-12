package types


func MakeKeysIterator(obj KVIterable) KIterator {
    kv_iterator := obj.Iterate()
    var k_iterator KIterator
    k_iterator = func() (key Equatable, next KIterator) {
        key, _, kv_iterator = kv_iterator()
        if kv_iterator == nil {
            return nil, nil
        }
        return key, k_iterator
    }
    return k_iterator
}

func MakeValuesIterator(obj KVIterable) Iterator {
    kv_iterator := obj.Iterate()
    var v_iterator Iterator
    v_iterator = func() (value interface{}, next Iterator) {
        _, value, kv_iterator = kv_iterator()
        if kv_iterator == nil {
            return nil, nil
        }
        return value, v_iterator
    }
    return v_iterator
}

