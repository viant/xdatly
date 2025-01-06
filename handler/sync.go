package handler

import "sync"

type DataSync map[string]*sync.RWMutex
