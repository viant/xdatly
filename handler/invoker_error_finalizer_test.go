package handler

import "context"

type invokerErrorOutput struct{}

func (*invokerErrorOutput) Finalize(context.Context, ComponentInvoker, error) error { return nil }

type invokerOutput struct{}

func (*invokerOutput) Finalize(context.Context, ComponentInvoker) error { return nil }

type invokerInput struct{}

func (*invokerInput) Init(context.Context, ComponentInvoker) error { return nil }

type injectorOutput struct{}

func (*injectorOutput) Finalize(context.Context, InjectorLookup) error { return nil }

// Existing contracts and two-argument authored scaffolds remain source compatible.
var (
	_ InvokerErrorFinalizer = (*invokerErrorOutput)(nil)
	_ InvokerFinalizer      = (*invokerOutput)(nil)
	_ InvokerInitializer    = (*invokerInput)(nil)
	_ InjectorFinalizer     = (*injectorOutput)(nil)
	_ Initializer           = (*testInitializer)(nil)
	_ Finalizer             = (*testFinalizer)(nil)
	_ ErrorFinalizer        = (*testErrorFinalizer)(nil)
)
