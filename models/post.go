package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Post 投稿の構造体
type Post struct {
	UserID      string        `bson:"userId"` // 投稿したユーザのID
	PostID      bson.ObjectId `bson:"postId"` // 一意な投稿ID
	Text        string        `bson:"text"`
	CreatedDate time.Time     `bson:"createdDate"`
}

func NewPost(uid, text string) *Post {
	return &Post{
		UserID:      uid,
		PostID:      bson.NewObjectId(),
		Text:        text,
		CreatedDate: time.Now(),
	}
}
