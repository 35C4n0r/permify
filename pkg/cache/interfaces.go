package cache

// Cache - Defines an interface for a generic cache.
type Cache interface {
	Get(key interface{}) (interface{}, bool)
	Set(key interface{}, entry interface{}, cost int64) bool
	Wait()
	Close()
}
