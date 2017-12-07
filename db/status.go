package db

import (
	"encoding/json"
	"errors"

	"github.com/TinyKitten/TimelineServer/models"
	"github.com/garyburd/redigo/redis"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
)

const (
	PostsCol = "posts"
)

func (m *MongoInstance) FindPost(postID bson.ObjectId, cached bool) (*models.Post, error) {
	sess := m.session.Clone()
	defer sess.Close()

	if cached {
		data, err := m.cache.GetStruct(postID.Hex())
		if err != nil && err != redis.ErrNil {
			return nil, err
		}
		if data != nil {
			u := m.deserializePost(data)
			return u, nil
		}
	}

	post := new(models.Post)
	if err := sess.DB(m.db()).C(PostsCol).
		FindId(postID).One(&post); err != nil {
		return nil, err
	}

	err := m.updatePostCache(*post)
	if err != nil {
		m.logger.Debug("Redis Error", zap.String("Error", err.Error()))
		return nil, err
	}

	return post, nil
}

func (m *MongoInstance) UpdatePost(post models.Post) error {
	sess := m.session.Clone()
	defer sess.Close()

	err := m.Insert(PostsCol, post)
	if err != nil {
		return err
	}

	err = m.AppendUserPost(post.UserID, post.ID)
	return err
}

func (m *MongoInstance) GetAllPosts() (*[]models.Post, error) {
	sess := m.session.Clone()
	defer sess.Close()

	var posts []models.Post
	if err := sess.DB(m.db()).C(PostsCol).
		Find(nil).
		All(&posts); err != nil {
		return nil, err
	}

	return &posts, nil
}

func (m *MongoInstance) GetPosts(limit int) (*[]models.Post, error) {
	sess := m.session.Clone()
	defer sess.Close()

	var posts []models.Post
	if err := sess.DB(m.db()).C(PostsCol).
		Find(nil).
		Limit(limit).
		All(&posts); err != nil {
		return nil, err
	}

	return &posts, nil
}

// GetPostsByOIDArray ObjectIDの配列で投稿を一括検索し一致した投稿の配列を返す
func (m *MongoInstance) GetPostsByOIDArray(objectIds []bson.ObjectId) ([]models.Post, error) {
	sess := m.session.Clone()
	defer sess.Close()

	if &objectIds == nil {
		return nil, errors.New("empty array")
	}
	array := []models.Post{}

	for _, postID := range objectIds {
		data, err := m.cache.GetStruct(postID.Hex())
		if err != nil && err != redis.ErrNil {
			return nil, err
		}
		if data != nil {
			post := m.deserializePost(data)
			array = append(array, *post)
		} else {
			u := models.Post{}
			err := sess.DB(m.db()).C(PostsCol).
				Find(bson.M{"_id": postID}).One(&u)
			if err != nil {
				return nil, err
			}
			array = append(array, u)
		}
	}
	return array, nil
}

func (m *MongoInstance) CreateLike(postID, userID bson.ObjectId) error {
	sess := m.session.Clone()
	defer sess.Close()

	err := sess.DB(m.db()).C(PostsCol).
		Update(bson.M{"_id": postID}, bson.M{"$addToSet": bson.M{"favoritedIds": userID}})
	if err != nil {
		return handleError(err)
	}

	updated, err := m.FindPost(postID, false)
	if err != nil {
		return handleError(err)
	}

	err = m.updatePostCache(*updated)
	if err != nil {
		m.logger.Debug("Redis Error", zap.String("Error", err.Error()))
		return err
	}

	return nil
}

func (m *MongoInstance) DestroyLike(postID, userID bson.ObjectId) error {
	sess := m.session.Clone()
	defer sess.Close()

	err := sess.DB(m.db()).C(PostsCol).
		Update(bson.M{"_id": postID}, bson.M{"$pull": bson.M{"favoritedIds": userID}})
	if err != nil {
		return handleError(err)
	}

	updated, err := m.FindPost(postID, false)
	if err != nil {
		return handleError(err)
	}

	err = m.updatePostCache(*updated)
	if err != nil {
		m.logger.Debug("Redis Error", zap.String("Error", err.Error()))
		return err
	}

	return nil
}

func (m *MongoInstance) updatePostCache(p models.Post) (err error) {
	_, err = m.cache.SetStruct(p.ID.Hex(), p)
	if err != nil {
		return err
	}
	return
}

func (m *MongoInstance) deserializePost(serialized []byte) *models.Post {
	deserialized := new(models.Post)
	json.Unmarshal(serialized, &deserialized)
	return deserialized
}

func (m *MongoInstance) deserializePostArray(serialized []byte) *[]models.Post {
	deserialized := new([]models.Post)
	json.Unmarshal(serialized, &deserialized)
	return deserialized
}
