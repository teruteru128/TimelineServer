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
	InReplyToUserID bson.ObjectId   `bson:"in_reply_to_user_id"`
	Text            string          `bson:"text"`
	Shared          []bson.ObjectId `bson:"shared"`
	UserID          bson.ObjectId   `bson:"user_id"`
}

type PostEntity struct {
	URLs         []string `json:"urls"`
	Hashtags     []string `json:"hashtags"`
	UserMentions []Post   `json:"user_mentions"`
}

/*
	Favorited           bool        `json:"favorited"`
	CreatedAt           string      `json:"created_at"`
	ID                  string      `json:"id"`
	Entities            PostEntity  `json:"entities"`
	InReplyToUserID     string      `json:"in_reply_to_user_id"`
	Text                string      `json:"text"`
	Shared              bool        `json:"shared"`
	SharedCount         int         `json:"shared_count"`
	User                models.User `json:"user"`
	InReplyToScreenName string      `json:"in_reply_to_screen_name"`
*/

func NewPost(uid, text, inReplyToStatusID string) *Post {
	return &Post{
		UserID:          bson.ObjectId(uid),
		ID:              bson.NewObjectId(),
		Text:            text,
		CreatedAt:       time.Now(),
		InReplyToUserID: bson.ObjectId(inReplyToStatusID),
	}
}
