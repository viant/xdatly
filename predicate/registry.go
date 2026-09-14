package predicate

import (
	"fmt"
	"sync"
	"time"
)

type (
	Registry struct {
		sync.Mutex
		notifications sync.Mutex
		notifiers     map[int64]NotifierFn
		registry      map[string]*Template
	}

	NotifierFn func(template *Template)
)

func (r *Registry) Lookup(name string) (*Template, error) {
	r.Lock()
	defer r.Unlock()
	ret, ok := r.registry[name]
	if !ok {
		return nil, fmt.Errorf("failed to lookup predicate template: %v", name)
	}
	return ret, nil
}

func New() *Registry {
	return &Registry{
		registry:  map[string]*Template{},
		notifiers: map[int64]NotifierFn{},
	}
}

var instance = New()

func RegisterTemplate(template *Template) {
	instance.register(template)
}

func Templates(callback NotifierFn) (map[string]*Template, func()) {
	return instance.templates(callback)
}

func (r *Registry) register(template *Template) {
	// Preserve serial notification delivery while permitting registry reads
	// and listener removal from a callback.
	r.notifications.Lock()
	defer r.notifications.Unlock()
	r.Lock()

	r.registry[template.Name] = template
	callbacks := make([]NotifierFn, 0, len(r.notifiers))
	for _, fn := range r.notifiers {
		callbacks = append(callbacks, fn)
	}
	r.Unlock()
	// Callbacks may look up templates or unsubscribe; never invoke them while
	// holding the registry lock. This registration uses a notifier snapshot.
	for _, fn := range callbacks {
		fn(template)
	}
}

func (r *Registry) templates(callback NotifierFn) (map[string]*Template, func()) {
	r.Lock()
	defer r.Unlock()

	result := map[string]*Template{}
	for _, template := range r.registry {
		result[template.Name] = template
	}

	var closer func()
	if callback != nil {
		key := r.key()
		r.notifiers[key] = callback
		closer = func() {
			r.Lock()
			defer r.Unlock()
			delete(r.notifiers, key)
		}
	}

	return result, closer
}

func (r *Registry) key() int64 {
	now := time.Now().Unix()
	for {
		if _, ok := r.notifiers[now]; ok {
			now++
			continue
		}
		return now
	}
}
