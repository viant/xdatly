package core

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

var debugEnabled = os.Getenv("DATLY_DEBUG") != ""

func RegisterType(packageName, typeName string, rType reflect.Type, insertedAt time.Time) {
	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
	}
	instance.register(packageName, typeName, rType, insertedAt)
	registerDependencies(packageName, rType, insertedAt)
}

func registerDependencies(packageName string, rType reflect.Type, insertedAt time.Time) {
	if depTypes := getDependentTypes(rType); len(depTypes) > 0 { //dependent type are method call types that share the same package
		for depType := range depTypes {
			name := depType.Name()
			if index := strings.LastIndex(name, "."); index != -1 {
				name = name[index:]
			}
			instance.register(packageName, name, depType, insertedAt)
		}
	}
}

func getDependentTypes(rType reflect.Type) map[reflect.Type]bool {
	var depTypes = map[reflect.Type]bool{}
	if rType.Kind() == reflect.Struct {
		structPtr := reflect.PtrTo(rType)
		discoverDependentTypes(rType, rType, depTypes)
		discoverDependentTypes(structPtr, rType, depTypes)
	}
	return depTypes
}

func discoverDependentTypes(srcType reflect.Type, rType reflect.Type, dep map[reflect.Type]bool) {
	for i := 0; i < srcType.NumMethod(); i++ {
		method := srcType.Method(i)
		fType := method.Func.Type()
		for j := 0; j < fType.NumIn(); j++ {
			argType := fType.In(j)
			if argType.Kind() == reflect.Ptr {
				argType = argType.Elem()
			}
			if argType.PkgPath() != rType.PkgPath() || argType == rType {
				continue
			}
			dep[argType] = true
		}
	}
}

func Types(callback NotifierFn) (types map[string]map[string]reflect.Type, closeFn func()) {
	return instance.types(callback)
}

var instance *registry

func init() {
	instance = newRegistry()
}

func newRegistry() *registry {
	return &registry{
		packages:  map[string]*packageRegistry{},
		notifiers: map[int64]NotifierFn{},
	}
}

type (
	NotifierFn func(packageName, typeName string, rType reflect.Type, insertedAt time.Time)

	registry struct {
		sync.Mutex
		packages  map[string]*packageRegistry
		notifiers map[int64]NotifierFn
	}

	packageRegistry struct {
		index map[string]*entry
	}

	entry struct {
		inserted time.Time
		rType    reflect.Type
	}
)

func (r *packageRegistry) register(typeName string, rType reflect.Type, insertedAt time.Time) bool {
	anEntry, ok := r.index[typeName]
	if !ok {
		anEntry = &entry{
			inserted: insertedAt,
			rType:    rType,
		}
		r.index[typeName] = anEntry
		return true
	}

	if anEntry.inserted.After(insertedAt) {
		return false
	}

	anEntry.rType = rType
	anEntry.inserted = insertedAt
	return true
}

func (r *packageRegistry) getType(name string) (reflect.Type, bool) {
	e, ok := r.index[name]
	if !ok {
		return nil, false
	}

	return e.rType, true
}

func (r *registry) register(packageName, typeName string, rType reflect.Type, insertedAt time.Time) {
	r.Lock()
	defer r.Unlock()

	index := r.packageRegistry(packageName)
	if index.register(typeName, rType, insertedAt) {
		if debugEnabled {
			fmt.Printf("[DEBUG] overriding type %v, %v\n", strings.Join([]string{packageName, typeName}, "."), rType.String())
		}

		for _, notifier := range r.notifiers {
			notifier(packageName, typeName, rType, insertedAt)
		}
	} else {
		if debugEnabled {
			fmt.Printf("[DEBUG] ignoring type override %v, types inserted before datly built are ignored\n", strings.Join([]string{packageName, typeName}, "."))
		}
	}
}

func (r *registry) packageRegistry(name string) *packageRegistry {
	index, ok := r.packages[name]
	if ok {
		return index
	}

	index = &packageRegistry{index: map[string]*entry{}}
	r.packages[name] = index
	return index
}

func (r *registry) types(callback NotifierFn) (map[string]map[string]reflect.Type, func()) {
	r.Lock()
	defer r.Unlock()

	result := map[string]map[string]reflect.Type{}
	for packageName, pRegistry := range r.packages {
		tRegistry := map[string]reflect.Type{}
		result[packageName] = tRegistry
		for tName, tEntry := range pRegistry.index {
			tRegistry[tName] = tEntry.rType
		}
	}

	if callback == nil {
		return result, func() {}
	}

	now := time.Now().Unix()
	for {
		if _, ok := r.notifiers[now]; ok {
			now++
		} else {
			r.notifiers[now] = callback
			return result, func() {
				r.Lock()
				defer r.Unlock()
				delete(r.notifiers, now)
			}
		}
	}
}
