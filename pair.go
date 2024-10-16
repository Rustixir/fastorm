package fastorm

type Pair struct {
	key   Key
	value Value
}

func NewPair(key Key, value Value) Pair {
	return Pair{
		key:   key,
		value: value,
	}
}

func (p Pair) GetKey() Key {
	return p.key
}

func (p Pair) GetValue() Value {
	return p.value
}
