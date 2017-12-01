package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Post 投稿の構造体
type Post struct {
	FavoritedIds    []bson.ObjectId `bson:"favoritedIds" json:"favorited_ids"`
	CreatedAt       time.Time       `bson:"createdAt" json:"created_at"`
	ID          	bson.ObjectId   `json:"id" bson:"_id,omitempty"`   // BSON ObjectID
	MentionsID      []bson.ObjectId `bson:"mentionsId" json:"mentions_id"`
	URLs            []string        `bson:"urls" json:"urls"`
	Hashtags        []string        `bson:"hashtags" json:"hashtags"`
	InReplyToUserID bson.ObjectId   `bson:"in_reply_to_user_id,omitempty" json:"in_reply_to_user_id"`
	Text            string          `bson:"text" json:"text"`
	Shared          []bson.ObjectId `bson:"shared" json:"shared"`
	UserID          bson.ObjectId   `bson:"user_id" json:"user_id"`
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
