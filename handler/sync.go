package handler

import "sync"

type syncKey string

// DataSync coordinates cross-view hook ordering. A relation registers its
// Holder key before fetching begins (Put), releases it when data is ready
// (Delete), and dependents block in Wait until the lock is released.
// Ported faithfully from xdatly/handler/sync.go.
type DataSync struct {
	aMap map[string]*sync.RWMutex
	rw   sync.RWMutex
}

func (d *DataSync) Put(key string) {
	d.rw.Lock()
	defer d.rw.Unlock()
	if _, ok := d.aMap[key]; !ok {
		lock := &sync.RWMutex{}
		lock.Lock()
		d.aMap[key] = lock
	}
}

func (d *DataSync) Get(key string) *sync.RWMutex {
	d.rw.Lock()
	defer d.rw.Unlock()
	if _, ok := d.aMap[key]; !ok {
		return nil
	}
	return d.aMap[key]
}

func (d *DataSync) Delete(key string) {
	d.rw.Lock()
	defer d.rw.Unlock()
	if lock, ok := d.aMap[key]; ok {
		lock.Unlock()
		delete(d.aMap, key)
	}
}

// Wait blocks until the named relation's data is available.
// Returns false when the key was never Put (no-op is safe).
func (d *DataSync) Wait(key string) bool {
	d.rw.RLock()
	lock, ok := d.aMap[key]
	d.rw.RUnlock()
	if !ok {
		return false
	}
	lock.RLock()
	lock.RUnlock()
	return true
}

func NewDataSync() *DataSync {
	return &DataSync{
		aMap: make(map[string]*sync.RWMutex),
	}
}

// DataSyncKey is the context key used to carry *DataSync through hooks that
// coordinate sibling view readiness.
const DataSyncKey = syncKey("dataSync")
