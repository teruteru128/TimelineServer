package db

import "github.com/TinyKitten/TimelineServer/models"

const (
	PostsCol = "posts"
)

func (m *MongoInstance) GetAllPosts(limit int) (*[]models.Post, error) {
	conn, err := m.getConnection()
	if err != nil {
		return nil, handleError(err)
	}
	var posts []models.Post
	if err := conn.C(PostsCol).
		Find(nil).
		Limit(limit).
		All(&posts); err != nil {
		return nil, err
	}

	return &posts, nil
}
