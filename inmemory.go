package fastorm

type inMemory struct {
	kvs map[Key]Value
}

func NewInMemory() DataStore {
	return &inMemory{
		kvs: make(map[Key]Value),
	}
}

func (i *inMemory) Set(key Key, value Value) error {
	i.kvs[key] = value
	return nil
}

func (i *inMemory) Get(key Key) (Value, error) {
	val, ok := i.kvs[key]
	if !ok {
		return val, ErrNotFound
	}
	return val, nil
}

func (i *inMemory) Delete(key Key) error {
	delete(i.kvs, key)
	return nil
}

func (i *inMemory) Range(f func(key Key, value Value) error) error {
	for k, v := range i.kvs {
		if err := f(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (i *inMemory) Pairs() []Pair {
	list := make([]Pair, 0, len(i.kvs))
	for key, value := range i.kvs {
		list = append(list, NewPair(key, value))
	}
	return list
}

func (i *inMemory) Keys() []Key {
	list := make([]Key, 0, len(i.kvs))
	for key := range i.kvs {
		list = append(list, key)
	}
	return list
}

func (i *inMemory) Len() (int, error) {
	return len(i.kvs), nil
}

func (i *inMemory) Clear() error {
	i.kvs = nil
	return nil
}
