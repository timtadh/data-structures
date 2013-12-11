package types

type String string


func (self String) Equals(other Equatable) bool {
    if o, ok := other.(String); ok {
        return self == o
    } else {
        return false
    }
}

func (self String) Less(other Sortable) bool {
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
