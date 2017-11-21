package cache

import (
	"strconv"
	"time"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/logger"
	"github.com/garyburd/redigo/redis"
	"go.uber.org/zap"
)

type RedisInstance struct {
	pool *redis.Pool
}

func newPool() *redis.Pool {
	cfg := config.GetCacheConfig()
	port := strconv.Itoa(cfg.Port)
	host := ":" + port
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", host) },
	}
}

func NewRedisInstance() RedisInstance {
	pool := newPool()
	ins := RedisInstance{
		pool: pool,
	}
	return ins
}

func handleError(err error) error {
	logger := logger.GetLogger()
	logger.Debug("Redis Error", zap.String("Error", err.Error()))
	return err
}
