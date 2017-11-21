package db

import (
	"github.com/TinyKitten/TimelineServer/models"
)

const (
	PostsCol = "posts"
)

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
