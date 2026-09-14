package plugin

import (
	xcodec "github.com/viant/xdatly/codec"
	xdocs "github.com/viant/xdatly/docs"
	xpredicate "github.com/viant/xdatly/predicate"
)

// Registry is the reduced public plugin/provider registration seam.
// It only covers the curated contract families that already exist in
// `xdatly`, and it is intentionally declarative. It must not regrow into a
// hidden global extension/bootstrap system.
type Registry interface {
	RegisterCodec(name string, instance xcodec.Instance)
	RegisterCodecFactory(name string, factory xcodec.Factory)
	RegisterPredicate(template *xpredicate.Template)
	RegisterDocs(name string, provider xdocs.Provider)
}
