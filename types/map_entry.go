package types

type MapEntry struct {
	Key Hashable
	Value interface{}
}

func (m *MapEntry) Equals(other Equatable) bool {
	if o, ok := other.(*MapEntry); ok {
		return m.Key.Equals(o.Key)
	} else {
		return false
	}
}

func (m *MapEntry) Less(other Sortable) bool {
	if o, ok := other.(*MapEntry); ok {
		return m.Key.Less(o.Key)
	} else {
		return false
	}
}

func (m *MapEntry) Hash() int {
	return m.Key.Hash()
}

