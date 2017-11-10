package db

import (
	"github.com/TinyKitten/Timeline/config"
	"github.com/TinyKitten/Timeline/logger"
	"github.com/TinyKitten/Timeline/models"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	UsersCol = "users"
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

func (m *MongoInstance) FindUser(userid string) (*models.User, error) {
	conn, err := m.getConnection()
	if err != nil {
		return nil, handleError(err)
	}
	u := new(models.User)
	if err := conn.C(UsersCol).
		Find(bson.M{"userid": userid}).One(&u); err != nil {
		return nil, err
	}
	return u, nil
}
