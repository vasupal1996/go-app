package redisstorage

import (
	"fmt"
	"go-app/server/config"
	"time"

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
	var client *redis.Pool
	client = &redis.Pool{
		MaxActive: 50,
		// Maximum number of idle connections in the pool.
		MaxIdle: 25,
		// max number of connections
		Dial: func() (redis.Conn, error) {
			redisConnStart := time.Now()
			fmt.Println(c.ConnectionURL())
			c, err := redis.Dial("tcp", c.ConnectionURL())
			redisConnDuration := time.Since(redisConnStart)
			fmt.Println("redis connection latency = s%", redisConnDuration)
			if err != nil {
				return nil, err
			}
			return c, err

		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	conn := client.Get()
	return &RedisStorage{Conn: conn}
}

// Close closes redis connection
func (rs *RedisStorage) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return rs.Conn.Do(commandName, args...)
}
