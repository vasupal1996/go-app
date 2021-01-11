package redisstorage

import (
	"go-app/server/config"
	"log"
	"os"

	"github.com/gomodule/redigo/redis"
)

// RedisStorage stores for redis client and redis config
type RedisStorage struct {
	Config *config.RedisConfig
	Conn   redis.Conn
}

// Close closes redis connection
func (rs *RedisStorage) Close() {
	rs.Conn.Close()
}

// NewRedisStorage returns new redis instance
func NewRedisStorage(c *config.RedisConfig) *RedisStorage {
	conn, err := redis.Dial(c.Network, c.ConnectionURL())
	if err != nil {
		log.Fatalf("failed to connect to redis: %s", err)
		os.Exit(1)
	}
	if _, err := conn.Do("PING"); err != nil {
		log.Fatalf("failed to ping redis: %s", err)
		os.Exit(1)
	}
	r := &RedisStorage{
		Config: c,
		Conn:   conn,
	}
	return r
}
