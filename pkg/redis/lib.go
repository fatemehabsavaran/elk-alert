package redis

import (
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisConfig struct {
	Addr       string `json:"addr,omitempty" mapstructure:"addr"`
	Password   string `json:"password,omitempty" mapstructure:"password"`
	DB         int    `json:"db,omitempty" mapstructure:"db"`
	PoolSize   int    `json:"pool_size,omitempty" mapstructure:"pool_size"`
	Timeout    int    `json:"timeout,omitempty" mapstructure:"timeout"`
	ClientName string `json:"client_name,omitempty" mapstructure:"client_name"`
}

type RedisConnector struct {
	cfg    RedisConfig
	client *redis.Client
}

func NewRedisConnector(cfg RedisConfig) *RedisConnector {
	client := redis.NewClient(&redis.Options{
		Addr:        cfg.Addr,
		Password:    cfg.Password,
		DB:          cfg.DB,
		PoolSize:    cfg.PoolSize,
		DialTimeout: time.Duration(cfg.Timeout) * time.Second,
		ClientName:  cfg.ClientName,
	})

	return &RedisConnector{
		cfg:    cfg,
		client: client,
	}
}

func (r *RedisConnector) GetClient() *redis.Client {
	return r.client
}
