package storage

// Impl used by server to implement any storage interface such as mongodb, sql, redis etc.
type Impl interface {
	Close()
}
