package handler

import (
	"context"
	"errors"
	"testing"
)

type testInitializer struct {
	calls int
}

func (i *testInitializer) Init(ctx context.Context) error {
	_ = ctx
	i.calls++
	return nil
}

type testFinalizer struct {
	calls int
}

func (f *testFinalizer) Finalize(ctx context.Context) error {
	_ = ctx
	f.calls++
	return nil
}

type testErrorFinalizer struct {
	calls int
	last  error
}

type testWriteLifecycle struct {
	initialized bool
	validated   bool
	err         error
}

func (v *testWriteLifecycle) InitWrite(context.Context) error {
	v.initialized = true
	return v.err
}

func (v *testWriteLifecycle) ValidateWrite(context.Context) error {
	v.validated = true
	return v.err
}

func (f *testErrorFinalizer) Finalize(ctx context.Context, err error) error {
	_ = ctx
	f.calls++
	f.last = err
	return nil
}

func TestLifecycleContractsCompile(t *testing.T) {
	var _ Initializer = &testInitializer{}
	var _ Finalizer = &testFinalizer{}
	var _ ErrorFinalizer = &testErrorFinalizer{}
	var _ WriteInitializer = &testWriteLifecycle{}
	var _ WriteValidator = &testWriteLifecycle{}
}

func TestWriteLifecycleShape(t *testing.T) {
	value := &testWriteLifecycle{}
	if err := WriteInitializer(value).InitWrite(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := WriteValidator(value).ValidateWrite(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !value.initialized || !value.validated {
		t.Fatalf("write lifecycle state = %+v", value)
	}
	expected := errors.New("write lifecycle failure")
	value.err = expected
	if err := WriteInitializer(value).InitWrite(context.Background()); !errors.Is(err, expected) {
		t.Fatalf("InitWrite() error = %v", err)
	}
}

func TestInitializerLifecycleShape(t *testing.T) {
	value := &testInitializer{}
	if err := value.Init(context.Background()); err != nil {
		t.Fatalf("unexpected init error: %v", err)
	}
	if value.calls != 1 {
		t.Fatalf("expected one init call, got %d", value.calls)
	}
}

func TestFinalizerLifecycleShape(t *testing.T) {
	value := &testFinalizer{}
	if err := value.Finalize(context.Background()); err != nil {
		t.Fatalf("unexpected finalize error: %v", err)
	}
	if value.calls != 1 {
		t.Fatalf("expected one finalize call, got %d", value.calls)
	}
}

func TestErrorFinalizerLifecycleShape(t *testing.T) {
	expected := errors.New("boom")
	value := &testErrorFinalizer{}
	if err := value.Finalize(context.Background(), expected); err != nil {
		t.Fatalf("unexpected error-finalize error: %v", err)
	}
	if value.calls != 1 {
		t.Fatalf("expected one error-finalize call, got %d", value.calls)
	}
	if !errors.Is(value.last, expected) {
		t.Fatalf("expected finalizer to receive original error, got %v", value.last)
	}
}
