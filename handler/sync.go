package handler

import "sync"

type DataSync map[string]*sync.RWMutex

func (d DataSync) Wait(key string) {
	lock, ok := d[key]
	if !ok {
		return
	}
	lock.RLock()
	lock.RUnlock()
}
