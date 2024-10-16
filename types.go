package fastorm

type Key = string
type Value = any

type IsolationLevel int

const (
	ReadUncommited IsolationLevel = 1
	ReadCommitted  IsolationLevel = 2
	RepeatableRead IsolationLevel = 3
	Serializable   IsolationLevel = 4
)

type StorageType string

const (
	DiscCopies StorageType = "discCopies"
	RamCopies  StorageType = "ramCopies"
)
