package db

import (
	"github.com/TinyKitten/TimelineServer/models"
)

const (
	PostsCol = "posts"
)

/*
func (m *MongoInstance) createPostCollection(col string) error {
	conn, err := m.getConnection()
	if err != nil {
		return handleError(err)
	}
	colExist := false
	names, err := conn.CollectionNames()
	if err != nil {
		return handleError(err)
	}
	for _, c := range names {
		if c == col {
			colExist = true
		}
	}
	if !colExist {
		info := mgo.CollectionInfo{
			Capped:   true,
			MaxBytes: 104857600,
		}
		return conn.C(col).Create(&info)
	}
	return nil
}
*/

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
