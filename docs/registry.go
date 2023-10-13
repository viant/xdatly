package docs

import "sync"

type Registry struct {
	byName map[string]Provider
	mux    sync.RWMutex
}

func (r *Registry) Register(name string, provider Provider) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.byName[name] = provider
}

func (r *Registry) ForEach(onElem func(name string, provider Provider)) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	for k, v := range r.byName {
		onElem(k, v)
	}
}

func (r *Registry) Lookup(name string) Provider {
	r.mux.RLock()
	defer r.mux.RUnlock()
	ret := r.byName[name]
	return ret
}

var register = &Registry{byName: map[string]Provider{}}

func Register(name string, provider Provider) {
	register.Register(name, provider)
}

func ForEach(onElem func(name string, provider Provider)) {
	register.ForEach(onElem)
}

func Lookup(name string) Provider {
	return register.Lookup(name)
}
