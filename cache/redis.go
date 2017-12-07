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

func newPool(conf config.CacheConfig) *redis.Pool {
	port := strconv.Itoa(conf.Port)
	host := conf.Server + ":" + port
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", host) },
	}
}

func NewRedisInstance(conf config.CacheConfig) RedisInstance {
	pool := newPool(conf)
	ins := RedisInstance{
		pool: pool,
	}
	return ins
}

func handleError(err error) error {
	logger := logger.GetLogger()
	if err == redis.ErrNil {
		logger.Info("Redis Info", zap.String("Error", err.Error()))
		return nil
	}
	logger.Error("Redis Error", zap.String("Error", err.Error()))
	return err
}
