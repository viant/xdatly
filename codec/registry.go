package codec

import (
	"sync"
	"time"
)

type (
	registry struct {
		sync.Mutex
		registry  map[string]*Codec
		notifiers map[int64]func(codec *Codec)
	}

	Codec struct {
		Name               string
		Instance           Instance
		instanceInsertedAt time.Time
		Factory            Factory
		factoryInsertedAt  time.Time
	}
)

var registryInstance = newRegistryInstance()

func Codecs(notifier func(codec *Codec)) (codecs map[string]*Codec, closer func()) {
	return registryInstance.codecs(notifier)
}

func RegisterCodec(name string, codec Instance, at time.Time) {
	registryInstance.registerInstance(name, codec, at)
}

func RegisterFactory(name string, factory Factory, at time.Time) {
	registryInstance.registerFactory(name, factory, at)
}

func newRegistryInstance() *registry {
	return &registry{
		registry:  map[string]*Codec{},
		notifiers: map[int64]func(codec *Codec){},
	}
}

func (r *registry) codecs(notifier func(codec *Codec)) (map[string]*Codec, func()) {
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

func (r *registry) registerFactory(name string, factory Factory, at time.Time) {
	r.Lock()
	defer r.Unlock()

	codec := r.getCodec(name)
	codec.Factory = factory
	r.notify(codec)
}

func (r *registry) registerInstance(name string, instance Instance, at time.Time) {
	r.Lock()
	defer r.Unlock()

	codec := r.getCodec(name)
	codec.Instance = instance
	r.notify(codec)
}

func (r *registry) getCodec(name string) *Codec {
	codec, ok := r.registry[name]
	if !ok {
		codec = &Codec{}
		r.registry[name] = codec
	}

	return codec
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

func (r *registry) notify(codec *Codec) {
	for _, f := range r.notifiers {
		f(codec)
	}
}
