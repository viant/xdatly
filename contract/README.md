# Caller-dependent contracts

An original component Input or Output can implement `Provider`:

```go
Contract(context.Context, bindly.Binder) (*contract.Policy, error)
```

The binder is supplied by the invocation host. Bind application-defined state
inside the method; no state type name or `handler.Session` is prescribed. Resolve
security decisions once, then return a policy whose name and callback describe
those decisions. Return `nil, nil` when the linked contract should remain intact.

```go
func (*Output) Contract(ctx context.Context, binder bindly.Binder) (*contract.Policy, error) {
    var state struct {
        Auth *auth.Output `parameter:",kind=component,in=GET:/auth"`
    }
    if err := binder.Bind(ctx, &state); err != nil {
        return nil, err
    }
    // Application policy; the registered auth component owns verification.
    if state.Auth == nil || state.Auth.Context == nil {
        return nil, errors.New("authorization context required")
    }
    if state.Auth.Context.IsAdmin() {
        return nil, nil
    }
    return &contract.Policy{
        Name: "financial/restricted/v1",
        Field: func(parentPath string, parentType reflect.Type, field *contract.Field) (*contract.Field, error) {
            if field.Tag.Get("access") != "financial" {
                return nil, nil
            }
            return field.Clone().Hide(), nil
        },
    }, nil
}
```

The contract uses the shared Bindly base interface, not Xdatly's Lookup extension.
V0 can adapt only this contract-adjustment path without changing its existing
session or handler binder APIs. V1's Xdatly Binder satisfies the shared base.

The example assumes a host which installs the component provider. It does not
claim arbitrary `codec` tags are automatically compiled by every Binder. Datly's
ordinary and canonical/scoped binding paths have different admission and codec
semantics. This SDK addition does not implement a metadata invocation scope.

## Host responsibilities

- Preserve the exact linked type at every unchanged node. A nil callback result
  or equal Name/Tag retains that field. Reconstruct only changed ancestors; keep
  unaffected siblings and embedded structure. `viant/x/shape.CompileRewrite` owns
  direct-field traversal; do not use flattened visible fields to rebuild types.
- Supply original parent type and canonical parent path. Only Name and Tag are
  editable. Maintain canonical binding/name mappings and original method authority.
- Check `Field.Validate(originalFieldType)` before accepting a replacement. Helpers
  record errors so fluent calls can be checked once. Typed enum normalization uses
  `tagly/choice`; enum metadata does not install runtime membership validation.
- `Hide` upserts `json:"-" internal:"true" mcp:"-"`. Required/Optional upsert explicit
  booleans; the consuming schema and validator must interpret those booleans, not
  mere tag presence. SetEnum replaces the ordered repeated-choice group.
- Cache immutable effective types/plans by generation, component/route, direction,
  original type, provider configuration and policy name. Empty names disable retained
  caching. Never cache callbacks or captured request state. Authenticate and compute
  current policy selection before reuse. A stable policy name must mean stable edits.
- Preserve native binding, authorization, scope lifetime, codec and presence behavior.
  A new shared interface is not permission to bypass a scoped binder's restrictions.

The package contains contracts and field helpers only. HTTP/MCP/OpenAPI integration,
v0 compatibility releases, codec adapters and the policy cache remain host work.
