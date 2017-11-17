package cache

import (
	"github.com/TinyKitten/DiscordBot/logger"
	"github.com/garyburd/redigo/redis"
	"go.uber.org/zap"
)

const loggerTopic = "Redis Error"

type RedisInstance struct {
	conn   redis.Conn
	logger *zap.Logger
}

func handleError(err error) error {
	logger := logger.GetLogger()
	logger.Debug("Redis Error", zap.String("Reason", err.Error()))
	return err
}

func (r *RedisInstance) getConnection() (redis.Conn, error) {
	logger := logger.GetLogger()

	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, handleError(err)
	}

	r.conn = c
	r.logger = logger

	return c, nil
}

func (r *RedisInstance) Do(cmd, key string, data interface{}) (interface{}, error) {
	c, err := r.getConnection()
	if err != nil {
		return nil, handleError(err)
	}
	defer c.Close()

	if data != nil {
		return c.Do(cmd, key, data)
	}
	return c.Do(cmd, key)
}
