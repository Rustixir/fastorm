package fastorm

import "sync"

type txn struct {
	id           uint64
	manager      *Database
	commitedData DataStore
	tempData     DataStore
	snapshot     DataStore
	deletedKeys  DataStore
	isolation    IsolationLevel
	active       bool
	txnLock      sync.Mutex
}

func (tx *txn) Commit() error {
	if !tx.active {
		return ErrTxnInactive
	}

	tx.txnLock.Lock()
	defer tx.txnLock.Unlock()

	if tx.isolation == Serializable {
		tx.manager.dataStoreLock.Lock()
		defer tx.manager.dataStoreLock.Unlock()
	}
	for _, key := range tx.deletedKeys.Keys() {
		if err := tx.commitedData.Delete(key); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	for _, pair := range tx.tempData.Pairs() {
		if err := tx.commitedData.Set(pair.GetKey(), pair.GetValue()); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	_ = tx.Clear()
	tx.active = false
	tx.manager.sessionLock.Lock()
	defer tx.manager.sessionLock.Unlock()
	delete(tx.manager.sessions, tx.id)
	return nil
}

func (tx *txn) Rollback() error {
	if !tx.active {
		return ErrTxnInactive
	}
	tx.txnLock.Lock()
	defer tx.txnLock.Unlock()
	_ = tx.Clear()
	tx.active = false
	tx.manager.sessionLock.Lock()
	defer tx.manager.sessionLock.Unlock()
	delete(tx.manager.sessions, tx.id)
	return nil
}

func (tx *txn) Set(key Key, value Value) error {
	if !tx.active {
		return ErrTxnInactive
	}
	return tx.tempData.Set(key, value)
}

func (tx *txn) Get(key Key) (Value, error) {
	if !tx.active {
		var empty Value
		return empty, ErrTxnInactive
	}

	if _, err := tx.deletedKeys.Get(key); err == nil {
		return nil, ErrNotFound
	}

	if val, err := tx.tempData.Get(key); err == nil {
		return val, nil
	}

	switch tx.isolation {
	case ReadUncommited:
		tx.manager.sessionLock.RLock()
		defer tx.manager.sessionLock.RUnlock()
		for _, session := range tx.manager.sessions {
			if val, err := session.Get(key); err == nil {
				return val, nil
			}
		}
		return tx.commitedData.Get(key)
	case ReadCommitted:
		return tx.commitedData.Get(key)
	case RepeatableRead:
		if val, err := tx.snapshot.Get(key); err == nil {
			return val, nil
		}
		val, err := tx.commitedData.Get(key)
		if err != nil {
			return nil, err
		}
		_ = tx.snapshot.Set(key, val)
		return val, nil
	case Serializable:
		tx.manager.dataStoreLock.Lock()
		defer tx.manager.dataStoreLock.Unlock()
		if val, err := tx.snapshot.Get(key); err == nil {
			return val, nil
		}
		val, err := tx.commitedData.Get(key)
		if err != nil {
			return nil, err
		}
		_ = tx.snapshot.Set(key, val)
		return val, nil
	}
	return nil, ErrInvalidIsolationLevel
}

func (tx *txn) Delete(key Key) error {
	if !tx.active {
		return ErrTxnInactive
	}
	if err := tx.deletedKeys.Set(key, nil); err != nil {
		return err
	}
	_ = tx.tempData.Delete(key)
	return nil
}

func (tx *txn) Range(f func(key Key, value Value) error) error {
	if !tx.active {
		return ErrTxnInactive
	}
	pairs := tx.commitedData.Pairs()
	var value any
	for _, p := range pairs {
		switch tx.isolation {
		case ReadUncommited:
			tx.manager.sessionLock.RLock()
			for _, session := range tx.manager.sessions {
				if val, err := session.Get(p.key); err == nil {
					value = val
					break
				}
			}
			tx.manager.sessionLock.RUnlock()
			if value == nil {
				value = p.GetValue()
			}
		case ReadCommitted:
			value = p.GetValue()
		case RepeatableRead:
			if val, err := tx.snapshot.Get(p.GetKey()); err == nil {
				value = val
			} else {
				value = p.GetValue()
				_ = tx.snapshot.Set(p.key, value)
			}
		case Serializable:
			tx.manager.dataStoreLock.Lock()
			value = p.GetValue()
			tx.manager.dataStoreLock.Unlock()
		}
		if err := f(p.key, value); err != nil {
			return err
		}
	}
	return nil
}

func (tx *txn) Pairs() []Pair {
	return tx.commitedData.Pairs()
}

func (tx *txn) Keys() []Key {
	return tx.commitedData.Keys()
}

func (tx *txn) Len() (int, error) {
	return tx.commitedData.Len()
}

func (tx *txn) Clear() error {
	_ = tx.tempData.Clear()
	_ = tx.deletedKeys.Clear()
	_ = tx.snapshot.Clear()
	return nil
}
