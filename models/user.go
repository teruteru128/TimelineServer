package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// User ユーザの構造体
type User struct {
	ID          bson.ObjectId   `json:"id" bson:"_id,omitempty"`   // BSON ObjectID
	UserID      string          `json:"screen_name" bson:"userId"` // ユーザ名(@kitten)
	DisplayName string          `json:"name" bson:"displayName"`   // 表示名(Kitten)
	Description string          `json:"description" bson:"description"`
	Password    string          `json:"password" bson:"password"`           // 暗号化済みパスワード
	EMail       string          `json:"email" bson:"email"`                 // メールアドレス
	Location    string          `json:"location" bson:"location"`           // 居住地(グンマー)
	Following   []bson.ObjectId `json:"friends" bson:"following"`           // フォローしているユーザーのセット
	Followers   []bson.ObjectId `json:"followers" bson:"followers"`         // フォローされているユーザーのセット
	Posts       []bson.ObjectId `json:"posts" bson:"posts"`                 // 投稿のセット
	WebsiteURL  string          `json:"url" bson:"websiteUrl"`              // ウェブサイトのURL(http://example.com)
	AvatarURL   string          `json:"profile_image_url" bson:"avatarUrl"` // プロフィール画像(http://static_cdn/profile_images/0.png)
	Suspended   bool            `json:"suspended" bson:"suspended"`         // 凍結フラグ(TRUE/FALSE)
	CreatedDate time.Time       `json:"created_at" bson:"createdDate"`      // ユーザ登録日時
	UpdatedDate time.Time       `json:"updated_at" bson:"updatedDate"`      // 最終更新日
	Official    bool            `json:"official" bson:"official"`           // 公式
}

// NewUser 初期化されたUser構造体を返す
func NewUser(id, password, mail string, isOfficial bool) *User {
	return &User{
		ID:          bson.NewObjectId(),
		UserID:      id,
		DisplayName: id,
		Password:    password,
		EMail:       mail,
		Location:    "",
		Following:   []bson.ObjectId{},
		Followers:   []bson.ObjectId{},
		Posts:       []bson.ObjectId{},
		WebsiteURL:  "",
		AvatarURL:   "",
		Suspended:   false,
		CreatedDate: time.Now(),
		UpdatedDate: time.Now(),
		Official:    isOfficial,
	}
}
