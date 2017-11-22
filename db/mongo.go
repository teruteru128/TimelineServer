package db

import (
	"github.com/TinyKitten/TimelineServer/cache"
	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/logger"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
)

const loggerTopic = "MongoDB Error"

type MongoInstance struct {
	session *mgo.Session
	logger  zap.Logger
	cache   cache.RedisInstance
}

func handleError(err error) error {
	logger := logger.GetLogger()
	logger.Debug("MongoDB Error", zap.String("Reason", err.Error()))
	return err
}

func (m *MongoInstance) url() string {
	cfg := config.GetDBConfig()
	if cfg.User == "" || cfg.Password == "" {
		return cfg.Server
	}
	return cfg.User + ":" + cfg.Password + "@" + cfg.Server
}

func (m *MongoInstance) db() string {
	return config.GetDBConfig().Database
}

func setIndex(s *mgo.Database) error {
	// users
	usersIndex := mgo.Index{
		Key:        []string{"userId", "email"},
		Unique:     true, // ユニーク
		DropDups:   true, // ユニークインデックスが付いているデータに対して上書きを許可しない
		Background: true, // バックグラウンドでインデックスを行う
		Sparse:     true, // nilのデータはインデックスしない
	}
	err := s.C("users").EnsureIndex(usersIndex)

	return err
}

func NewMongoInstance() (*MongoInstance, error) {
	db := config.GetDBConfig().Database
	m := MongoInstance{}
	session, err := mgo.Dial(m.url())
	if err != nil {
		return nil, handleError(err)
	}
	session.SetSafe(&mgo.Safe{})
	m.session = session
	err = setIndex(session.DB(db))
	if err != nil {
		return nil, err
	}
	logger := logger.GetLogger()
	m.logger = *logger

	redisInstance := cache.NewRedisInstance()
	m.cache = redisInstance
	return &m, nil
}

func (m *MongoInstance) GetCollection(key string) (*mgo.Collection, error) {
	sess := m.session.Clone()
	defer sess.Close()

	col := sess.DB(m.db()).C(key)
	return col, nil
}

func (m *MongoInstance) Ping() (err error) {
	sess := m.session.Clone()
	defer sess.Close()
	err = sess.Ping()
	return nil
}

func (m *MongoInstance) Insert(key string, data interface{}) error {
	sess := m.session.Clone()
	defer sess.Close()

	return sess.DB(m.db()).C(key).Insert(data)
}

func (m *MongoInstance) InsertWithCache(key string, data interface{}) error {
	sess := m.session.Clone()
	defer sess.Close()

	_, err := m.cache.SetStruct(key, data)
	if err != nil {
		m.logger.Debug("Redis Error", zap.String("Error", err.Error()))
	}

	return sess.DB(m.db()).C(key).Insert(data)
}
