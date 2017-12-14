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
	conf    config.DBConfig
}

func handleError(err error) error {
	logger := logger.NewLogger()
	logger.Error("MongoDB Error", zap.String("Reason", err.Error()))
	return err
}

func (m *MongoInstance) url() string {
	return m.conf.Server + "/" + m.db()
}

func (m *MongoInstance) db() string {
	return m.conf.Database
}

func setIndex(s *mgo.Database) (err error) {
	// users
	usersIndex := mgo.Index{
		Key:        []string{"userId", "email"},
		Unique:     true, // ユニーク
		DropDups:   true, // ユニークインデックスが付いているデータに対して上書きを許可しない
		Background: true, // バックグラウンドでインデックスを行う
		Sparse:     true, // nilのデータはインデックスしない
	}
	err = s.C("users").EnsureIndex(usersIndex)

	return
}

func NewMongoInstance(conf config.DBConfig, cacheConf config.CacheConfig) (*MongoInstance, error) {
	m := MongoInstance{}
	m.conf = conf
	session, err := mgo.Dial(m.url())
	if err != nil {
		return nil, handleError(err)
	}
	session.SetSafe(&mgo.Safe{})
	m.session = session
	err = setIndex(session.DB(m.conf.Database))
	if err != nil {
		return nil, err
	}
	logger := logger.NewLogger()
	m.logger = *logger

	redisInstance := cache.NewRedisInstance(cacheConf)
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
