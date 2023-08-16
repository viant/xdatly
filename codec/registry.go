package codec

import (
	"fmt"
	"sync"
	"time"
)

type (
	Registry struct {
		sync.Mutex
		registry  map[string]*Codec
		notifiers map[int64]func(codec *Codec)
	}

	Codec struct {
		Name               string
		Instance           Instance
		InstanceInsertedAt time.Time
		Factory            Factory
		FactoryInsertedAt  time.Time
	}

	RegistryOption func(registry *Registry)
)

func WithFactory(name string, factory Factory, at time.Time) RegistryOption {
	return func(registry *Registry) {
		registry.RegisterFactory(name, factory, at)
	}
}

func WithCodec(name string, codec Instance, at time.Time) RegistryOption {
	return func(registry *Registry) {
		registry.RegisterInstance(name, codec, at)
	}
}

var registryInstance = NewRegistry()

func Codecs(notifier func(codec *Codec)) (codecs map[string]*Codec, closer func()) {
	return registryInstance.Codecs(notifier)
}

func RegisterCodec(name string, codec Instance, at time.Time) {
	registryInstance.RegisterInstance(name, codec, at)
}

func RegisterFactory(name string, factory Factory, at time.Time) {
	registryInstance.RegisterFactory(name, factory, at)
}

func Register(name string, codec *Codec) {
	registryInstance.RegisterCodec(name, codec)
}

func NewRegistry(opts ...RegistryOption) *Registry {
	r := &Registry{
		registry:  map[string]*Codec{},
		notifiers: map[int64]func(codec *Codec){},
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Registry) Codecs(notifier func(codec *Codec)) (map[string]*Codec, func()) {
	r.Lock()
	defer r.Unlock()

	result := map[string]*Codec{}
	for key, codec := range r.registry {
		result[key] = codec
	}

	var closer func()
	if notifier != nil {
		key := r.key()
		r.notifiers[key] = notifier
		closer = func() {
			r.Lock()
			defer r.Unlock()
			delete(r.notifiers, key)
		}
	}

	return result, closer
}

func (r *Registry) RegisterCodec(name string, codec *Codec) {
	r.Lock()
	defer r.Unlock()
	r.registry[name] = codec
	r.notify(codec)
}

func (r *Registry) RegisterFactory(name string, factory Factory, at time.Time) {
	r.Lock()
	defer r.Unlock()

	codec := r.getCodec(name)
	codec.Factory = factory
	r.notify(codec)
}

func (r *Registry) RegisterInstance(name string, instance Instance, at time.Time) {
	r.Lock()
	defer r.Unlock()

	codec := r.getCodec(name)
	codec.Instance = instance
	r.notify(codec)
}

func (r *Registry) getCodec(name string) *Codec {
	codec, ok := r.registry[name]
	if !ok {
		codec = &Codec{}
		r.registry[name] = codec
	}

	return codec
}

func (r *Registry) Lookup(name string) (*Codec, error) {
	r.Lock()
	defer r.Unlock()

	codec, ok := r.registry[name]
	if !ok {
		return nil, fmt.Errorf("not found codec %v", name)
	}

	return codec, nil
}

func (r *Registry) key() int64 {
	now := time.Now().Unix()
	for {
		if _, ok := r.notifiers[now]; ok {
			now++
		} else {
			return now
		}
	}
}

func (r *Registry) notify(codec *Codec) {
	for _, f := range r.notifiers {
		f(codec)
	}
}
