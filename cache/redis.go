package cache

import (
	"time"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/logger"
	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
	"go.uber.org/zap"
)

type RedisInstance struct {
	pool *redis.Pool
}

func newPool(conf config.CacheConfig) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redisurl.ConnectToURL(conf.Server)
		},
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
