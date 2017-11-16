package db

import (
	"github.com/TinyKitten/TimelineServer/models"
	"gopkg.in/mgo.v2/bson"
)

const (
	// UsersCol DB上のUser用カラム
	UsersCol = "users"
)

// FindUserByOID ObjectIDでユーザを検索する
func (m *MongoInstance) FindUserByOID(objectID bson.ObjectId) (*models.User, error) {
	conn, err := m.getConnection()
	if err != nil {
		return nil, handleError(err)
	}
	u := new(models.User)
	if err := conn.C(UsersCol).
		Find(bson.M{"_id": objectID}).One(&u); err != nil {
		return nil, err
	}
	return u, nil
}

// FindUserByOIDArray ObjectIDの配列でユーザーを一括検索し一致したユーザの配列を返す
func (m *MongoInstance) FindUserByOIDArray(objectIds []bson.ObjectId) ([]models.User, error) {
	conn, err := m.getConnection()
	if err != nil {
		return nil, handleError(err)
	}
	u := []models.User{}
	if err := conn.C(UsersCol).
		Find(bson.M{"_id": bson.M{"$in": objectIds}}).All(&u); err != nil {
		return nil, err
	}
	return u, nil
}

// FindUser userid(displayName)でユーザーを検索する
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

// DeleteUser userid(displayName)に一致したユーザを削除する
func (m *MongoInstance) DeleteUser(userid string) error {
	conn, err := m.getConnection()
	if err != nil {
		return handleError(err)
	}
	return conn.C(UsersCol).Remove(bson.M{"userid": userid})
}

// SuspendUser ObjectIDに一致したユーザを凍結する
func (m *MongoInstance) SuspendUser(objectID bson.ObjectId, flag bool) error {
	conn, err := m.getConnection()
	if err != nil {
		return handleError(err)
	}
	return conn.C(UsersCol).
		Update(bson.M{"_id": objectID}, bson.M{"$set": bson.M{"suspended": flag}})
}

// FollowUser fromOIDのユーザからtoOIDのユーザをフォローする
func (m *MongoInstance) FollowUser(fromOID, toOID bson.ObjectId) error {
	conn, err := m.getConnection()
	if err != nil {
		return handleError(err)
	}
	err = conn.C(UsersCol).
		Update(bson.M{"_id": fromOID}, bson.M{"$push": bson.M{"following": toOID}})
	if err != nil {
		return handleError(err)
	}
	err = conn.C(UsersCol).
		Update(bson.M{"_id": toOID}, bson.M{"$push": bson.M{"followers": fromOID}})
	if err != nil {
		handleError(err)
	}
	return nil
}

// UnfollowUser fromOIDのユーザがフォローしているユーザからtoOIDのユーザのフォローを解除する
func (m *MongoInstance) UnfollowUser(fromOID, toOID bson.ObjectId) error {
	conn, err := m.getConnection()
	if err != nil {
		return handleError(err)
	}
	err = conn.C(UsersCol).
		Update(bson.M{"_id": fromOID}, bson.M{"$pull": bson.M{"following": toOID}})
	if err != nil {
		return handleError(err)
	}
	err = conn.C(UsersCol).
		Update(bson.M{"_id": toOID}, bson.M{"$pull": bson.M{"followers": fromOID}})
	if err != nil {
		handleError(err)
	}
	return nil
}

// SetOfficial ユーザにを公式アカウントに設定するか、剥奪する
func (m *MongoInstance) SetOfficial(objectID bson.ObjectId, flag bool) error {
	conn, err := m.getConnection()
	if err != nil {
		return handleError(err)
	}
	return conn.C(UsersCol).
		Update(bson.M{"_id": objectID}, bson.M{"$set": bson.M{"official": flag}})
}
