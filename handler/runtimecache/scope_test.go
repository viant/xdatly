package runtimecache

import (
	"context"
	"reflect"
	"testing"
)

func TestSealedReaderScope(t *testing.T) {
	type input struct{ ID int }
	original, prepared := &input{1}, &input{}
	root, owner := new(int), new(int)
	id := Identity{Root: root, Input: original}
	ctx, err := InstallActual(context.Background(), owner, id, []Index{{"campaign_id", "CampaignId"}})
	check.NoError(t, err)
	values := []any{7}
	check.NoError(t, Select(ctx, "campaign_id", values))
	values[0] = 99
	check.NoError(t, Select(ctx, "campaign_id", []any{7}))
	check.Error(t, Select(ctx, "order_id", []any{7}))
	snapshot, err := Seal(ctx, owner, id)
	check.NoError(t, err)
	participated, denied, column, parameter, selected, err := snapshot.Selection(owner, id)
	check.NoError(t, err)
	check.True(t, participated)
	check.False(t, denied)
	check.Equal(t, []any{7}, selected)
	selected[0] = 88
	check.Error(t, Select(ctx, column, []any{7}))
	check.Error(t, Deny(ctx))
	identity, err := EnterIdentity(ctx, owner, snapshot, id, prepared, parameter, column, []any{7})
	check.NoError(t, err)
	check.NoError(t, CheckInput(identity, prepared))
	check.Error(t, CheckInput(identity, original))
	check.Error(t, CheckInput(identity, nil))
	check.NoError(t, CheckInput(ctx, original))
	check.Error(t, CheckInput(ctx, prepared))
	yes, err := IsIdentity(identity)
	check.NoError(t, err)
	check.True(t, yes)
	yes, err = IsIdentity(ctx)
	check.NoError(t, err)
	check.False(t, yes)
	_, err = EnterIdentity(ctx, owner, snapshot, id, original, parameter, column, []any{7})
	check.Error(t, err)
	_, err = EnterIdentity(ctx, owner, snapshot, id, prepared, "OrderId", column, []any{7})
	check.Error(t, err)
	_, err = EnterIdentity(ctx, owner, snapshot, id, prepared, parameter, column, []any{99})
	check.Error(t, err)
	child, err := InstallActual(identity, new(int), Identity{Root: root, Input: prepared}, nil)
	check.NoError(t, err)
	yes, err = IsIdentity(child)
	check.NoError(t, err)
	check.False(t, yes)
	_, _, _, _, _, err = snapshot.Selection(new(int), id)
	check.Error(t, err)
	Close(ctx, owner)
	_, err = IsIdentity(identity)
	check.Error(t, err)
	_, _, _, _, _, err = snapshot.Selection(owner, id)
	check.Error(t, err)
}
func TestVetoAndInvalidSelection(t *testing.T) {
	owner, root, input := new(int), new(int), new(int)
	id := Identity{Root: root, Input: input}
	ctx, err := InstallActual(context.Background(), owner, id, []Index{{"campaign_id", "CampaignId"}})
	check.NoError(t, err)
	check.Error(t, Select(ctx, "campaign_id", []any{[]int{1}}))
	check.NoError(t, Deny(ctx))
	check.NoError(t, Select(ctx, "campaign_id", []any{1}))
	p, err := Seal(ctx, owner, id)
	check.NoError(t, err)
	participated, denied, _, _, _, err := p.Selection(owner, id)
	check.NoError(t, err)
	check.True(t, participated)
	check.True(t, denied)
	_, err = EnterIdentity(ctx, owner, p, id, new(int), "CampaignId", "campaign_id", []any{1})
	check.Error(t, err)
	for _, indexes := range [][]Index{nil, {{"campaign_id", "CampaignId"}, {"campaign_id", "Other"}}} {
		ctx, err := InstallActual(context.Background(), owner, id, indexes)
		check.NoError(t, err)
		check.NoError(t, Select(ctx, "campaign_id", []any{1}))
		_, err = Seal(ctx, owner, id)
		check.Error(t, err)
	}
	var empty *Snapshot
	_, _, _, _, _, err = empty.Selection(owner, id)
	check.Error(t, err)
	_, err = EnterIdentity(ctx, owner, nil, id, new(int), "CampaignId", "campaign_id", []any{1})
	check.Error(t, err)
	_, _, _, _, _, err = new(Snapshot).Selection(owner, id)
	check.Error(t, err)
	check.NoError(t, Select(context.Background(), "campaign_id", []any{1}))
	check.NoError(t, Deny(context.Background()))
}

type cacheAssertions struct{}

var check cacheAssertions

func (cacheAssertions) NoError(t *testing.T, e error) {
	t.Helper()
	if e != nil {
		t.Fatal(e)
	}
}
func (cacheAssertions) Error(t *testing.T, e error) {
	t.Helper()
	if e == nil {
		t.Fatal("expected error")
	}
}
func (cacheAssertions) True(t *testing.T, v bool) {
	t.Helper()
	if !v {
		t.Fatal("expected true")
	}
}
func (cacheAssertions) False(t *testing.T, v bool) {
	t.Helper()
	if v {
		t.Fatal("expected false")
	}
}
func (cacheAssertions) Equal(t *testing.T, want, got any) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v got %v", want, got)
	}
}
