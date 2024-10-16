package fastorm

import (
	"sync"
	"sync/atomic"
)

type Database struct {
	sessions      map[uint64]Transaction
	commitedData  DataStore
	tempData      DataStore
	counter       uint64
	sessionLock   *sync.RWMutex
	dataStoreLock *sync.Mutex
}

func NewDatabase() *Database {
	return &Database{
		commitedData:  DataStoreBuilder(),
		sessions:      make(map[uint64]Transaction),
		sessionLock:   &sync.RWMutex{},
		dataStoreLock: &sync.Mutex{},
		counter:       0,
	}
}

func (tm *Database) Begin(level IsolationLevel) Transaction {
	txnID := atomic.AddUint64(&tm.counter, 1)
	txn := &txn{
		id:           txnID,
		manager:      tm,
		commitedData: tm.commitedData,
		tempData:     DataStoreBuilder(),
		deletedKeys:  DataStoreBuilder(),
		snapshot:     NewInMemory(),
		isolation:    level,
		active:       true,
	}
	tm.sessionLock.Lock()
	tm.sessions[txnID] = txn
	tm.sessionLock.Unlock()
	return txn
}
