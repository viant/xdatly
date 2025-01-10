package handler

import "sync"

// DataSync represents data sync
type DataSync struct {
	aMap map[string]*sync.RWMutex
	rw   sync.RWMutex
}

func (d *DataSync) Put(key string) {
	d.rw.Lock()
	defer d.rw.Unlock()
	if _, ok := d.aMap[key]; !ok {
		d.aMap[key] = &sync.RWMutex{}
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
	if _, ok := d.aMap[key]; ok {
		delete(d.aMap, key)

	}
}

func (d *DataSync) Wait(key string) {
	d.rw.RLock()
	lock, ok := d.aMap[key]
	d.rw.RUnlock()
	if !ok {
		return
	}
	lock.RLock()
	lock.RUnlock()
}
