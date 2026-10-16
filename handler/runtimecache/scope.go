// Package runtimecache carries invocation-local reader cache decisions. Only a
// reader holding its original owner token can seal and consume a decision.
package runtimecache

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"sync"
)

type contextKey struct{}
type Identity struct{ Component, Root, Input any }
type Index struct{ Column, Parameter string }
type scope struct {
	mu               sync.Mutex
	owner            any
	id               Identity
	indexes          []Index
	phase            uint8
	selected, denied bool
	column           string
	values           []any
}
type Snapshot struct {
	s         *scope
	parameter string
}
type envelope struct {
	s                 *scope
	prepared          any
	parameter, column string
	values            []any
}
type frame struct {
	s        *scope
	identity *envelope
}

func comparable(v any) bool { return v != nil && reflect.TypeOf(v).Comparable() }
func InstallActual(ctx context.Context, owner any, id Identity, indexes []Index) (context.Context, error) {
	if !comparable(owner) || !comparable(id.Root) || !comparable(id.Input) || (id.Component != nil && !comparable(id.Component)) {
		return nil, fmt.Errorf("cache scope requires comparable owner and reader identity")
	}
	s := &scope{owner: owner, id: id, indexes: append([]Index(nil), indexes...)}
	return context.WithValue(ctx, contextKey{}, &frame{s: s}), nil
}
func current(ctx context.Context) *frame { f, _ := ctx.Value(contextKey{}).(*frame); return f }
func scalars(values []any) ([]any, error) {
	result := append([]any(nil), values...)
	for _, v := range result {
		if v == nil {
			return nil, fmt.Errorf("cache index requires immutable scalar values")
		}
		r := reflect.ValueOf(v)
		switch r.Kind() {
		case reflect.Bool, reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		case reflect.Float32, reflect.Float64:
			if math.IsNaN(r.Float()) || math.IsInf(r.Float(), 0) {
				return nil, fmt.Errorf("cache index requires finite scalar values")
			}
		default:
			return nil, fmt.Errorf("cache index requires immutable scalar values")
		}
	}
	return result, nil
}
func Select(ctx context.Context, column string, values []any) error {
	f := current(ctx)
	if f == nil {
		return nil
	}
	s := f.s
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.phase != 0 || f.identity != nil {
		return fmt.Errorf("cache selection outside actual reader build")
	}
	if column == "" || len(values) == 0 {
		return fmt.Errorf("cache selection requires column and values")
	}
	copied, err := scalars(values)
	if err != nil {
		return err
	}
	if s.selected && (s.column != column || !reflect.DeepEqual(s.values, copied)) {
		return fmt.Errorf("conflicting cache index selection")
	}
	s.selected = true
	s.column = column
	s.values = copied
	return nil
}
func Deny(ctx context.Context) error {
	f := current(ctx)
	if f == nil {
		return nil
	}
	s := f.s
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.phase != 0 || f.identity != nil {
		return fmt.Errorf("cache veto outside actual reader build")
	}
	s.denied = true
	return nil
}
func Seal(ctx context.Context, owner any, id Identity) (*Snapshot, error) {
	f := current(ctx)
	if f == nil {
		return nil, fmt.Errorf("missing cache scope")
	}
	s := f.s
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.owner != owner || s.id != id || s.phase != 0 || f.identity != nil {
		return nil, fmt.Errorf("cache scope ownership or phase mismatch")
	}
	parameter := ""
	if s.selected && !s.denied {
		count := 0
		for _, idx := range s.indexes {
			if idx.Column == s.column {
				count++
				parameter = idx.Parameter
			}
		}
		if count != 1 || parameter == "" {
			return nil, fmt.Errorf("cache selection requires exactly one canonical index")
		}
	}
	s.phase = 1
	return &Snapshot{s: s, parameter: parameter}, nil
}
func (p *Snapshot) Selection(owner any, id Identity) (participated, denied bool, column, parameter string, values []any, err error) {
	if p == nil || p.s == nil {
		err = fmt.Errorf("missing cache snapshot")
		return
	}
	s := p.s
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.owner != owner || s.id != id || s.phase != 1 {
		err = fmt.Errorf("cache snapshot ownership or phase mismatch")
		return
	}
	return s.selected || s.denied, s.denied, s.column, p.parameter, append([]any(nil), s.values...), nil
}
func EnterIdentity(ctx context.Context, owner any, p *Snapshot, id Identity, prepared any, parameter, column string, values []any) (context.Context, error) {
	participated, denied, wantColumn, wantParameter, wantValues, err := p.Selection(owner, id)
	if err != nil {
		return nil, err
	}
	f := current(ctx)
	if f == nil || f.s != p.s || f.identity != nil || !participated || denied || !comparable(prepared) || prepared == id.Input || parameter != wantParameter || column != wantColumn || !reflect.DeepEqual(values, wantValues) {
		return nil, fmt.Errorf("cache identity does not match sealed reader selection")
	}
	return context.WithValue(ctx, contextKey{}, &frame{s: p.s, identity: &envelope{s: p.s, prepared: prepared, parameter: parameter, column: column, values: wantValues}}), nil
}
func IsIdentity(ctx context.Context) (bool, error) {
	f := current(ctx)
	if f == nil {
		return false, nil
	}
	s := f.s
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.phase == 2 {
		return false, fmt.Errorf("closed cache scope")
	}
	if f.identity == nil {
		return false, nil
	}
	e := f.identity
	if s.phase != 1 || e.s != s || s.denied || !s.selected || e.column != s.column || !reflect.DeepEqual(e.values, s.values) {
		return false, fmt.Errorf("invalid cache identity envelope")
	}
	return true, nil
}
func Close(ctx context.Context, owner any) {
	f := current(ctx)
	if f == nil {
		return
	}
	s := f.s
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.owner == owner {
		s.phase = 2
	}
}

// CheckInput rejects an accidentally reused predicate input within this scope.
// The reader bridge itself owns identity; this check compares the input exposed
// by a native predicate with the corresponding actual or prepared object.
func CheckInput(ctx context.Context, input any) error {
	f := current(ctx)
	if f == nil {
		return nil
	}
	s := f.s
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.phase == 2 {
		return fmt.Errorf("closed cache scope")
	}
	expected := s.id.Input
	if f.identity != nil {
		expected = f.identity.prepared
		if input == nil {
			return fmt.Errorf("cache identity requires prepared input")
		}
	}
	if input != nil && (!comparable(input) || input != expected) {
		return fmt.Errorf("cache predicate input does not match reader scope")
	}
	return nil
}
