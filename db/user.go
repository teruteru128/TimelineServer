package db

import (
	"github.com/TinyKitten/TimelineServer/models"
	"gopkg.in/mgo.v2/bson"
)

const (
	UsersCol = "users"
)

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

func (m *MongoInstance) DeleteUser(userid string) error {
	conn, err := m.getConnection()
	if err != nil {
		return handleError(err)
	}
	return conn.C(UsersCol).Remove(bson.M{"userid": userid})
}

func (m *MongoInstance) SuspendUser(userid string, flag bool) error {
	conn, err := m.getConnection()
	if err != nil {
		return handleError(err)
	}
	return conn.C(UsersCol).Update(bson.M{"userid": userid}, bson.M{"$set": bson.M{"suspended": flag}})
}
