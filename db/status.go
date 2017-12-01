package db

import (
	"github.com/TinyKitten/TimelineServer/models"
	"gopkg.in/mgo.v2/bson"
	"github.com/garyburd/redigo/redis"
	"errors"
	"encoding/json"
)

const (
	PostsCol = "posts"
)

func (m *MongoInstance) UpdatePost(post models.Post) (error) {
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

// FindUserByOIDArray ObjectIDの配列で投稿を一括検索し一致した投稿の配列を返す
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

func (m *MongoInstance) deserializePost(serialized []byte) *models.Post {
	deserialized := new(models.Post)
	json.Unmarshal(serialized, &deserialized)
	return deserialized
}
