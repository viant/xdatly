package runtimecache

import (
	"context"
	"testing"
)

func TestSealUnwarmedAuthorizationSelectionAsDenied(t *testing.T) {
	for _, indexes := range [][]Index{nil, {{Column: "campaign_id", Parameter: "CampaignID"}}} {
		owner, root, input := new(int), new(int), new(int)
		id := Identity{Root: root, Input: input}
		ctx, err := InstallActual(context.Background(), owner, id, indexes)
		check.NoError(t, err)
		values := []any{77}
		check.NoError(t, Select(ctx, "order_id", values))
		values[0] = 88
		_, err = Seal(ctx, new(int), id)
		check.Error(t, err)
		_, err = Seal(ctx, owner, Identity{Root: root, Input: new(int)})
		check.Error(t, err)
		snapshot, err := Seal(ctx, owner, id)
		check.NoError(t, err)
		participated, denied, column, parameter, selected, err := snapshot.Selection(owner, id)
		check.NoError(t, err)
		check.True(t, participated)
		check.True(t, denied)
		check.Equal(t, "order_id", column)
		check.Equal(t, "", parameter)
		check.Equal(t, []any{77}, selected)
		selected[0] = 99
		_, _, _, _, again, err := snapshot.Selection(owner, id)
		check.NoError(t, err)
		check.Equal(t, []any{77}, again)
		_, err = EnterIdentity(ctx, owner, snapshot, id, new(int), "CampaignID", "campaign_id", []any{77})
		check.Error(t, err)
		check.Error(t, Select(ctx, "campaign_id", []any{77}))
		check.Error(t, Deny(ctx))
		_, err = Seal(ctx, owner, id)
		check.Error(t, err)
		_, _, _, _, _, err = snapshot.Selection(new(int), id)
		check.Error(t, err)
		identity, err := IsIdentity(ctx)
		check.NoError(t, err)
		check.False(t, identity)
		check.NoError(t, CheckInput(ctx, input))
		Close(ctx, owner)
		_, _, _, _, _, err = snapshot.Selection(owner, id)
		check.Error(t, err)
	}
}

func TestSealWithoutAuthorizationSelectionPreservesNonParticipation(t *testing.T) {
	for _, indexes := range [][]Index{nil, {{Column: "campaign_id", Parameter: "CampaignID"}}} {
		owner, root, input := new(int), new(int), new(int)
		id := Identity{Root: root, Input: input}
		ctx, err := InstallActual(context.Background(), owner, id, indexes)
		check.NoError(t, err)
		snapshot, err := Seal(ctx, owner, id)
		check.NoError(t, err)
		participated, denied, column, parameter, values, err := snapshot.Selection(owner, id)
		check.NoError(t, err)
		check.False(t, participated)
		check.False(t, denied)
		check.Equal(t, "", column)
		check.Equal(t, "", parameter)
		check.Equal(t, []any(nil), values)
		_, err = EnterIdentity(ctx, owner, snapshot, id, new(int), "CampaignID", "campaign_id", []any{77})
		check.Error(t, err)
		Close(ctx, owner)
	}
}
