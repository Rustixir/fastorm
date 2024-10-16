package fastorm

type Transaction interface {
	Commit() error
	Rollback() error
	DataStore
}

type DataStore interface {

	// Set return err if key already exist
	Set(key Key, value Value) error

	// Get return value if not found return err
	Get(key Key) (Value, error)

	// Delete if occurred error return err
	Delete(key Key) error

	// Range iterate over kvs
	Range(f func(key Key, value Value) error) error

	Pairs() []Pair

	Keys() []Key

	Len() (int, error)

	Clear() error
}
