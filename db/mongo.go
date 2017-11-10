package db

import (
	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/logger"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
)

type MongoInstance struct {
	Conf    config.DBConfig
	session *mgo.Session
}

func handleError(err error) error {
	logger := logger.GetLogger()
	logger.Error("MongoDB Error", zap.String("Reason", err.Error()))
	return err
}

func (m *MongoInstance) url() string {
	if m.Conf.User == "" || m.Conf.Password == "" {
		return m.Conf.Server
	}
	return m.Conf.User + ":" + m.Conf.Password + "@" + m.Conf.Server
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

func (m *MongoInstance) getConnection() (*mgo.Database, error) {
	db := config.GetDBConfig().Database
	if m.session != nil {
		return m.session.DB(db), nil
	}
	session, err := mgo.Dial(m.url())
	if err != nil {
		return nil, handleError(err)
	}
	session.SetMode(mgo.Monotonic, true)
	session.SetSafe(&mgo.Safe{})
	m.session = session
	err = setIndex(session.DB(db))
	if err != nil {
		return nil, err
	}
	return session.DB(db), nil
}

func (m *MongoInstance) Ping() (err error) {
	conn, err := m.getConnection()
	err = conn.Session.Ping()
	return nil
}

func (m *MongoInstance) Create(key string, data interface{}) error {
	conn, err := m.getConnection()
	if err != nil {
		return handleError(err)
	}
	return conn.C(key).Insert(data)
}
