package cache

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
)

func (r *RedisInstance) serialize(data interface{}) ([]byte, error) {
	serialized, err := json.Marshal(data)
	return serialized, err
}

func (r *RedisInstance) SetStruct(key string, data interface{}) (interface{}, error) {
	conn := r.pool.Get()
	defer conn.Close()

	serialized, err := r.serialize(data)
	if err != nil {
		return "", err
	}
	return conn.Do("SET", key, serialized)
}

func (r *RedisInstance) SetStructArray(key string, data []interface{}) (interface{}, error) {
	conn := r.pool.Get()
	defer conn.Close()

	serialized, err := r.serialize(data)
	if err != nil {
		return "", handleError(err)
	}
	return conn.Do("SET", key, serialized)
}

func (r *RedisInstance) GetStruct(key string) ([]byte, error) {
	conn := r.pool.Get()
	defer conn.Close()

	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, handleError(err)
	}
	return data, nil
}
