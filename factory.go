package fastorm

func DataStoreBuilder() DataStore {
	return NewInMemory()
}
