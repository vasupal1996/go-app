package memorystorage

import (
	"sync"
	"time"
)

type item struct {
	object     []byte
	expiration int64
}

// MemoryStore stores key value pair data inside memory
type MemoryStore struct {
	items       map[string]item
	mu          sync.RWMutex
	stopCleanup chan bool
}

// NewMemoryStorage returns a new MemoryStore instance, with a background cleanup goroutine that
// runs every minute to remove expired session data.
func NewMemoryStorage() *MemoryStore {
	mms := MemoryStore{}
	return &mms
}

// NewMemoryStorageWithCleanupInterval returns a new MemoryStore instance. The cleanupInterval
// parameter controls how frequently expired session data is removed by the
// background cleanup goroutine. Setting it to 0 prevents the cleanup goroutine
// from running (i.e. expired sessions will not be removed).
func NewMemoryStorageWithCleanupInterval(cleanupInterval time.Duration) *MemoryStore {
	m := &MemoryStore{
		items: make(map[string]item),
	}

	if cleanupInterval > 0 {
		go m.startCleanup(cleanupInterval)
	}

	return m
}

func (m *MemoryStore) startCleanup(interval time.Duration) {
	m.stopCleanup = make(chan bool)
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			m.deleteExpired()
		case <-m.stopCleanup:
			ticker.Stop()
			return
		}
	}
}

// Close implement storage interface method. Since memory store does not involve any external connection therefore this
// method does nothing.
func (m *MemoryStore) Close() {}

// StopCleanup terminates the background cleanup goroutine for the MemoryStore
// instance. It's rare to terminate this; generally MemoryStore instances and
// their cleanup goroutines are intended to be long-lived and run for the lifetime
// of your application.
//
// There may be occasions though when your use of the MemoryStore is transient.
// An example is creating a new MemoryStore instance in a test function. In this
// scenario, the cleanup goroutine (which will run forever) will prevent the
// MemoryStore object from being garbage collected even after the test function
// has finished. You can prevent this by manually calling StopCleanup.
func (m *MemoryStore) StopCleanup() {
	if m.stopCleanup != nil {
		m.stopCleanup <- true
	}
}

func (m *MemoryStore) deleteExpired() {
	now := time.Now().UnixNano()
	m.mu.Lock()
	for token, item := range m.items {
		if now > item.expiration {
			delete(m.items, token)
		}
	}
	m.mu.Unlock()
}

// Find returns the data for a given session token from the MemoryStore instance.
// If the session token is not found or is expired, the returned exists flag will
// be set to false.
func (m *MemoryStore) Find(token string) ([]byte, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, found := m.items[token]
	if !found {
		return nil, false, nil
	}

	if time.Now().UnixNano() > item.expiration {
		return nil, false, nil
	}

	return item.object, true, nil
}

// Commit adds a session token and data to the MemoryStore instance with the given
// expiry time. If the session token already exists, then the data and expiry
// time are updated.
func (m *MemoryStore) Commit(token string, b []byte, expiry time.Time) error {
	m.mu.Lock()
	m.items[token] = item{
		object:     b,
		expiration: expiry.UnixNano(),
	}
	m.mu.Unlock()

	return nil
}

// Delete removes a session token and corresponding data from the MemoryStore
// instance.
func (m *MemoryStore) Delete(token string) error {
	m.mu.Lock()
	delete(m.items, token)
	m.mu.Unlock()

	return nil
}

func (m *MemoryStore) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return reply, nil
}
