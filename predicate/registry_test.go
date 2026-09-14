package predicate

import (
	"sync"
	"testing"
	"time"
)

func TestRegistryLookupAndTemplates(t *testing.T) {
	reg := New()
	template := &Template{
		Name:   "by_id",
		Source: "ID = ?",
		Args: []*NamedArgument{
			{Name: "id", Position: 0},
		},
	}

	reg.register(template)

	lookup, err := reg.Lookup("by_id")
	if err != nil {
		t.Fatalf("expected lookup to succeed, got %v", err)
	}
	if lookup != template {
		t.Fatalf("expected lookup to return registered template")
	}

	snapshot, closer := reg.templates(nil)
	if closer != nil {
		t.Fatalf("expected nil closer without callback")
	}
	if len(snapshot) != 1 || snapshot["by_id"] != template {
		t.Fatalf("expected templates snapshot to preserve registry content")
	}
}

func TestRegistryNotifierLifecycle(t *testing.T) {
	reg := New()
	var notified []string

	snapshot, closer := reg.templates(func(template *Template) {
		if template != nil {
			notified = append(notified, template.Name)
		}
	})
	if len(snapshot) != 0 {
		t.Fatalf("expected empty snapshot on new registry")
	}
	if closer == nil {
		t.Fatalf("expected closer when callback is provided")
	}

	reg.register(&Template{Name: "first"})
	closer()
	reg.register(&Template{Name: "second"})

	if len(notified) != 1 || notified[0] != "first" {
		t.Fatalf("expected notifier to stop after closer, got %v", notified)
	}
}

func TestRegistryConcurrentLookupAndRegistration(t *testing.T) {
	reg := New()
	reg.register(&Template{Name: "shared", Source: "id = ?"})
	notifications := 0
	_, closeListener := reg.templates(func(*Template) { notifications++ })
	defer closeListener()
	var group sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		group.Add(1)
		go func(write bool) {
			defer group.Done()
			for i := 0; i < 100; i++ {
				if write {
					reg.register(&Template{Name: "shared", Source: "id = ?"})
				} else if _, err := reg.Lookup("shared"); err != nil {
					t.Error(err)
				}
			}
		}(worker%2 == 0)
	}
	group.Wait()
	if notifications != 400 {
		t.Fatalf("notifications=%d", notifications)
	}
}

func TestRegistryCallbackMayLookupAndUnsubscribe(t *testing.T) {
	reg := New()
	var closeListener func()
	_, closeListener = reg.templates(func(template *Template) {
		if _, err := reg.Lookup(template.Name); err != nil {
			t.Error(err)
		}
		closeListener()
	})
	done := make(chan struct{})
	go func() { reg.register(&Template{Name: "one"}); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("callback deadlocked registry")
	}
	reg.register(&Template{Name: "two"})
}
