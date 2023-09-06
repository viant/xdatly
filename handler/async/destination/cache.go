package destination

type Cache struct {
	Cache          *string `json:",omitempty"`
	CacheKey       *string `json:",omitempty"`
	CacheSet       *string `json:",omitempty"`
	CacheNamespace *string `json:",omitempty"`
}
