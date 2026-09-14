package destination

type Cache struct {
	Cache          *string `json:"cache,omitempty"`
	CacheKey       *string `json:"cacheKey,omitempty"`
	CacheSet       *string `json:"cacheSet,omitempty"`
	CacheNamespace *string `json:"cacheNamespace,omitempty"`
}
