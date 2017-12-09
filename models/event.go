package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// EventType イベントの種類
type EventType int

const (
	// FollowEvent フォローされた
	FollowEvent EventType = iota
	// UnfollowEvent フォローを外された
	UnfollowEvent
	// LikedEvent いいねされた
	LikedEvent
	// DislikedEvent いいねが取り消された
	DislikedEvent
	// SharedEvent シェアされた
	SharedEvent
	// ReceivedReplyEvent ポストに返信された
	ReceivedReplyEvent
)

// Event イベント
type Event struct {
	// ID 識別用ID
	ID bson.ObjectId `bson:"_id,omitempty" json:"id"`
	// Type イベントタイプ
	Type EventType `bson:"type" json:"type"`
	// FromUserID 通知元ID
	FromUserID bson.ObjectId `bson:"from_user_id" json:"from_user_id"`
	// ToUserID 通知送信先ID
	ToUserID bson.ObjectId `bson:"to_user_id" json:"to_user_id"`
	// AlreadyRead 既読
	AlreadyRead bool `bson:"already_read" json:"already_read"`
	// CreatedAt イベントが発生した日時
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	// TargetPostID イベント対象のポストID
	TargetPostID bson.ObjectId `bson:"post_id,omitempty" json:"post_id"`
}
