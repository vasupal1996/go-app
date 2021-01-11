package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config struct stores entire project configurations
type Config struct {
	ServerConfig     ServerConfig     `mapstructure:"server"`
	APIConfig        APIConfig        `mapstructure:"api"`
	KafkaConfig      KafkaConfig      `mapstructure:"kafka"`
	LoggerConfig     LoggerConfig     `mapstructure:"logger"`
	DatabaseConfig   DatabaseConfig   `mapstructure:"database"`
	RedisConfig      RedisConfig      `mapstructure:"redis"`
	MiddlewareConfig MiddlewareConfig `mapstructure:"middleware"`

	// CacheConfig    CacheConfig    `mapstructure:"cache"`

	// ZipkinConfig   ZipkinConfig   `mapstructure:"zipkin"`
	TokenAuthConfig TokenAuthConfig `mapstructure:"token"`
}

// ServerConfig has only server specific configuration
type ServerConfig struct {
	ListenAddr      string        `mapstructure:"listenAddr"`
	Port            string        `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"readTimeout"`
	WriteTimeout    time.Duration `mapstructure:"writeTimeout"`
	CloseTimeout    time.Duration `mapstructure:"closeTimeout"`
	Env             string        `mapstructure:"env"`
	UserMemoryStore bool          `mapstructure:"useMemoryStore"`
}

// APIConfig contains api package related configurations
type APIConfig struct {
	Mode               string `mapstructure:"mode"`
	EnableTestRoute    bool   `mapstructure:"enableTestRoute"`
	EnableMediaRoute   bool   `mapstructure:"enableMediaRoute"`
	EnableStaticRoute  bool   `mapstructure:"enableStaticRoute"`
	MaxRequestDataSize int    `mapstructure:"maxRequestDataSize"`
}

// TokenAuthConfig contains token authentication related configuration
type TokenAuthConfig struct {
	JWTSignKey string `mapstructure:"jwtSignKey"`
}

// KafkaConfig has kafka cluster specific configuration
type KafkaConfig struct {
	BrokerDial string `mapstructure:"brokerDial"`
	BrokerURL  string `mapstructure:"brokerUrl"`
	BrokerPort string `mapstructure:"brokerPort"`
}

// LoggerConfig contains logger specific configuration
type LoggerConfig struct {
	EnableKafkaLog   bool   `mapstructure:"enableKafkaLog"`
	EnableConsoleLog bool   `mapstructure:"enableConsoleLog"`
	KafkaTopic       string `mapstructure:"kafkaTopic"`
	KafkaPartition   string `mapstructure:"kafkaPartition"`
}

// DatabaseConfig contains mongodb related configuration
type DatabaseConfig struct {
	Scheme string `mapstructure:"scheme"`
	Host   string `mapstructure:"host"`
	// Name     string `mapstructure:"name"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	ReplicaSet string `mapstructure:"replicaSet"`
}

// ConnectionURL returns connection string to of mongodb storage
func (d *DatabaseConfig) ConnectionURL() string {
	url := fmt.Sprintf("%s://", d.Scheme)
	if d.Username != "" && d.Password != "" {
		url += fmt.Sprintf("%s:%s@", d.Username, d.Password)
	}
	url += fmt.Sprintf("%s", d.Host)
	if d.ReplicaSet != "" {
		url += fmt.Sprintf("?replicaSet=%s", d.ReplicaSet)
	}
	return url
}

// RedisConfig has cache related configuration.
type RedisConfig struct {
	Network  string `mapstructure:"network"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// ConnectionURL returns connection string to of mongodb storage
func (r *RedisConfig) ConnectionURL() string {
	var url string
	if r.Username != "" {
		url += fmt.Sprintf("%s", r.Username)
	}
	if r.Password != "" {
		url += fmt.Sprintf(":%s@", r.Password)
	}
	url += fmt.Sprintf("%s", r.Host)
	if r.Port != "" {
		url += fmt.Sprintf(":%s", r.Port)
	}
	return url
}

// MiddlewareConfig has middlewares related configuration
type MiddlewareConfig struct {
	EnableRequestLog bool `mapstructure:"enableRequestLog"`
}

// GetConfig returns entire project configuration
func GetConfig() *Config {
	return getConfigFromFile("default")
}

func getConfigFromFile(fileName string) *Config {
	if fileName == "" {
		fileName = "default"
	}

	// looking for filename `default` inside `src/server` dir with `.toml` extension
	viper.SetConfigName(fileName)
	viper.AddConfigPath("./conf/")
	viper.SetConfigType("toml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("couldn't load config: %s", err)
		os.Exit(1)
	}
	config := &Config{}
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("couldn't read config: %s", err)
		os.Exit(1)
	}
	return config
}
