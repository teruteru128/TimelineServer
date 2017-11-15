package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Post 投稿の構造体
type Post struct {
	UserID      string        `json:"userId"` // 投稿したユーザのID
	PostID      bson.ObjectId `json:"postId"` // 一意な投稿ID
	Text        string        `json:"text"`
	CreatedDate time.Time     `json:"createdDate"`
}

func NewPost(uid, text string) *Post {
	return &Post{
		UserID:      uid,
		PostID:      bson.NewObjectId(),
		Text:        text,
		CreatedDate: time.Now(),
	}
}
