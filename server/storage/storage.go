package storage

// DB used by server to implement any storage interface by redis client.
type DB interface {
	Close()
}

// Redis used by server to implement any storage interface by redis client.
type Redis interface {
	Close()
}
