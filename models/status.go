package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Post 投稿の構造体
type Post struct {
	FavoritedIds    []bson.ObjectId `bson:"favoritedIds"`
	CreatedAt       time.Time       `bson:"createdAt"`
	ID              bson.ObjectId   `bson:"id"`
	MentionsID      []bson.ObjectId `bson:"mentionsId"`
	URLs            []string        `bson:"urls"`
	Hashtags        []string        `bson:"hashtags"`
	InReplyToUserID bson.ObjectId   `bson:"in_reply_to_user_id,omitempty"`
	Text            string          `bson:"text"`
	Shared          []bson.ObjectId `bson:"shared"`
	UserID          bson.ObjectId   `bson:"user_id"`
}

type PostEntity struct {
	URLs         []string `json:"urls"`
	Hashtags     []string `json:"hashtags"`
	UserMentions []Post   `json:"user_mentions"`
}

func NewPost(uid, inReplyToStatusID bson.ObjectId, text string) *Post {
	return &Post{
		UserID:          uid,
		ID:              bson.NewObjectId(),
		Text:            text,
		CreatedAt:       time.Now(),
		InReplyToUserID: inReplyToStatusID,
	}
}
