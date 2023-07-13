package predicate

import (
	"fmt"
	"sync"
)

type (
	registry struct {
		lock     sync.Mutex
		notifier NotifierFn
		registry map[string]*Template
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
	instance = &registry{registry: map[string]*Template{}}
	RegisterTemplate(&Template{Name: "exists",
		Source: `
 EXISTS (SELECT 1 FROM $Table t WHERE t.$Column = $FilterValue AND t.$JoinColumn = p.$ParentColumn) 
`,
		Args: []*NamedArgument{
			{Name: "Table", Position: 0},
			{Name: "Column", Position: 1},
			{Name: "JoinColumn", Position: 2},
			{Name: "ParentColumn", Position: 3},
		}})
}

func RegisterTemplate(template *Template) {
	instance.registry[template.Name] = template
	if instance.notifier != nil {
		instance.notifier(template)
	}
}

func Templates(callback NotifierFn) (types map[string]*Template) {
	instance.lock.Lock()
	defer instance.lock.Unlock()
	instance.notifier = callback
	result := map[string]*Template{}
	for _, template := range instance.registry {
		result[template.Name] = template
	}

	return instance.registry
}
