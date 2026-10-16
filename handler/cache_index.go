package handler

import (
	"context"
	"github.com/viant/xdatly/handler/runtimecache"
)

// SelectCacheIndex retains the authorized values for this reader's indexed replay.
func SelectCacheIndex(ctx context.Context, column string, values []any) error {
	if err := runtimecache.CheckInput(ctx, ctx.Value(InputKey)); err != nil {
		return err
	}
	return runtimecache.Select(ctx, column, values)
}

// DenyCacheIndex vetoes indexed replay without disabling exact query caching.
func DenyCacheIndex(ctx context.Context) error {
	if err := runtimecache.CheckInput(ctx, ctx.Value(InputKey)); err != nil {
		return err
	}
	return runtimecache.Deny(ctx)
}

// CacheIndexIdentity reports the reader's isolated, validated cache key render.
func CacheIndexIdentity(ctx context.Context) (bool, error) {
	if err := runtimecache.CheckInput(ctx, ctx.Value(InputKey)); err != nil {
		return false, err
	}
	return runtimecache.IsIdentity(ctx)
}
