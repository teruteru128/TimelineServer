package db

import (
	"encoding/json"
	"errors"

	"github.com/garyburd/redigo/redis"
	"go.uber.org/zap"

	"github.com/TinyKitten/TimelineServer/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// UsersCol DB上のUser用カラム
	UsersCol = "users"
)

// FindUserByOID ObjectIDでユーザを検索する
func (m *MongoInstance) FindUserByOID(objectID bson.ObjectId) (*models.User, error) {
	sess := m.session.Clone()
	defer sess.Close()

	data, err := m.cache.GetStruct(objectID.Hex())
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	if data != nil {
		u := m.deserializeUser(data)
		return u, nil
	}

	u := new(models.User)
	if err := sess.DB(m.db()).C(UsersCol).
		FindId(objectID).One(&u); err != nil {
		return nil, err
	}

	err = m.updateUserCache(*u)
	if err != nil {
		m.logger.Debug("Redis Error", zap.String("Error", err.Error()))
		return nil, err
	}

	return u, nil
}

// FindUserByOIDArray ObjectIDの配列でユーザーを一括検索し一致したユーザの配列を返す
func (m *MongoInstance) FindUserByOIDArray(objectIds []bson.ObjectId) ([]models.User, error) {
	sess := m.session.Clone()
	defer sess.Close()

	if &objectIds == nil {
		return nil, errors.New("empty array")
	}
	array := []models.User{}

	for _, objectID := range objectIds {
		data, err := m.cache.GetStruct(objectID.Hex())
		if err != nil && err != redis.ErrNil {
			return nil, err
		}
		if data != nil {
			u := m.deserializeUser(data)
			array = append(array, *u)
		} else {
			u := models.User{}
			err := sess.DB(m.db()).C(UsersCol).
				Find(bson.M{"_id": objectID}).One(&u)
			if err != nil {
				return nil, err
			}
			array = append(array, u)
		}
	}
	return array, nil
}

// FindUser userid(displayName)でユーザーを検索する
func (m *MongoInstance) FindUser(userid string) (*models.User, error) {
	sess := m.session.Clone()
	defer sess.Close()

	data, err := m.cache.GetStruct(userid)
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	if data != nil {
		u := m.deserializeUser(data)
		return u, nil
	}

	u := new(models.User)
	if err := sess.DB(m.db()).C(UsersCol).
		Find(bson.M{"userId": userid}).One(&u); err != nil {
		return nil, err
	}

	err = m.updateUserCache(*u)
	if err != nil {
		m.logger.Debug("Redis Error", zap.String("Error", err.Error()))
		return nil, err
	}

	return u, nil
}

// DeleteUser userid(displayName)に一致したユーザを削除する
func (m *MongoInstance) DeleteUser(userid string) error {
	sess := m.session.Clone()
	defer sess.Close()

	return sess.DB(m.db()).C(UsersCol).Remove(bson.M{"userId": userid})
}

// SuspendUser ObjectIDに一致したユーザを凍結する
func (m *MongoInstance) SuspendUser(objectID bson.ObjectId, flag bool) error {
	sess := m.session.Clone()
	defer sess.Close()

	return sess.DB(m.db()).C(UsersCol).
		Update(bson.M{"_id": objectID}, bson.M{"$set": bson.M{"suspended": flag}})
}

// FollowUser fromOIDのユーザからtoOIDのユーザをフォローする
func (m *MongoInstance) FollowUser(fromOID, toOID bson.ObjectId) error {
	sess := m.session.Clone()
	defer sess.Close()

	user, err := m.FindUserByOID(fromOID)
	if err != nil {
		return handleError(err)
	}
	_, err = m.FindUserByOIDArray(user.Following)
	if err != nil {
		if err == mgo.ErrNotFound {
			return errors.New("already followed")
		} else {
			return handleError(err)
		}
	}

	err = sess.DB(m.db()).C(UsersCol).
		Update(bson.M{"_id": fromOID}, bson.M{"$push": bson.M{"following": toOID}})
	if err != nil {
		return handleError(err)
	}
	err = sess.DB(m.db()).C(UsersCol).
		Update(bson.M{"_id": toOID}, bson.M{"$push": bson.M{"followers": fromOID}})
	if err != nil {
		return handleError(err)
	}
	return nil
}

// UnfollowUser fromOIDのユーザがフォローしているユーザからtoOIDのユーザのフォローを解除する
func (m *MongoInstance) UnfollowUser(fromOID, toOID bson.ObjectId) error {
	sess := m.session.Clone()
	defer sess.Close()

	err := sess.DB(m.db()).C(UsersCol).
		Update(bson.M{"_id": fromOID}, bson.M{"$pull": bson.M{"following": toOID}})
	if err != nil {
		return handleError(err)
	}
	err = sess.DB(m.db()).C(UsersCol).
		Update(bson.M{"_id": toOID}, bson.M{"$pull": bson.M{"followers": fromOID}})
	if err != nil {
		return handleError(err)
	}
	return nil
}

// SetOfficial ユーザにを公式アカウントに設定するか、剥奪する
func (m *MongoInstance) SetOfficial(objectID bson.ObjectId, flag bool) error {
	sess := m.session.Clone()
	defer sess.Close()

	return sess.DB(m.db()).C(UsersCol).
		Update(bson.M{"_id": objectID}, bson.M{"$set": bson.M{"official": flag}})
}

func (m *MongoInstance) deserializeUser(serialized []byte) *models.User {
	deserialized := new(models.User)
	json.Unmarshal(serialized, &deserialized)
	return deserialized
}

func (m *MongoInstance) deserializeUserArray(serialized []byte) *[]models.User {
	deserialized := new([]models.User)
	json.Unmarshal(serialized, &deserialized)
	return deserialized
}

func (m *MongoInstance) UpdateUser(objectID bson.ObjectId, key string, value interface{}) error {
	sess := m.session.Clone()
	defer sess.Close()

	err := sess.DB(m.db()).C(UsersCol).
		Update(bson.M{"_id": objectID}, bson.M{"$set": bson.M{key: value}})
	if err != nil {
		m.logger.Debug("MongoDB Error", zap.String("Error", err.Error()))
		return err
	}

	var u models.User
	if err := sess.DB(m.db()).C(UsersCol).
		FindId(objectID).One(&u); err != nil {
	}
	if err != nil {
		m.logger.Debug("MongoDB Error", zap.String("Error", err.Error()))
		return err
	}

	err = m.updateUserCache(u)
	if err != nil {
		m.logger.Debug("Redis Error", zap.String("Error", err.Error()))
		return err
	}

	return nil
}

func (m *MongoInstance) updateUserCache(u models.User) (err error) {
	_, err = m.cache.SetStruct(u.ID.Hex(), u)
	if err != nil {
		return err
	}
	_, err = m.cache.SetStruct(u.UserID, u)
	if err != nil {
		return err
	}
	return
}
