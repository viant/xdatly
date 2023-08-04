package predicate

import (
	"fmt"
	"sync"
	"time"
)

type (
	registry struct {
		sync.Mutex
		notifiers map[int64]NotifierFn
		registry  map[string]*Template
	}
	NotifierFn func(template *Template)
)

func (r *registry) Lookup(name string) (*Template, error) {
	ret, ok := r.registry[name]
	if !ok {
		return nil, fmt.Errorf("failed to lookup predicate template: %v", name)
	}
	return ret, nil
}

var instance *registry

func init() {
	instance = &registry{registry: map[string]*Template{}, notifiers: map[int64]NotifierFn{}}
}

func RegisterTemplate(template *Template) {
	instance.register(template)
}

func Templates(callback NotifierFn) (types map[string]*Template, closer func()) {
	return instance.templates(callback)
}

func (r *registry) register(template *Template) {
	r.Lock()
	defer r.Unlock()

	r.registry[template.Name] = template
	for _, fn := range r.notifiers {
		fn(template)
	}
}

func (r *registry) templates(callback NotifierFn) (map[string]*Template, func()) {
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

func (r *registry) key() int64 {
	now := time.Now().Unix()
	for {
		if _, ok := r.notifiers[now]; ok {
			now++
		} else {
			return now
		}
	}
}
