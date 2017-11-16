package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// User ユーザの構造体
type User struct {
	ID          bson.ObjectId   `json:"_id" bson:"_id,omitempty"` // BSON ObjectID
	UserID      string          `json:"userId"`                   // ユーザ名(@kitten)
	DisplayName string          `json:"displayName"`              // 表示名(Kitten)
	Password    string          `json:"password"`                 // 暗号化済みパスワード
	EMail       string          `json:"email"`                    // メールアドレス
	Location    string          `json:"location"`                 // 居住地(グンマー)
	Following   []bson.ObjectId `json:"following"`                // フォローしているユーザーのセット
	Followers   []bson.ObjectId `json:"followers"`                // フォローされているユーザーのセット
	Posts       []bson.ObjectId `json:"posts"`                    // 投稿のセット
	WebsiteURL  string          `json:"websiteUrl"`               // ウェブサイトのURL(http://example.com)
	AvatarURL   string          `json:"avatarUrl"`                // プロフィール画像(http://static_cdn/profile_images/0.png)
	Suspended   bool            `json:"suspended"`                // 凍結フラグ(TRUE/FALSE)
	CreatedDate time.Time       `json:"createdDate"`              // ユーザ登録日時
	UpdatedDate time.Time       `json:"updatedDate"`              // 最終更新日
}

// NewUser 初期化されたUser構造体を返す
func NewUser(id, password, mail string) *User {
	return &User{
		ID:          bson.NewObjectId(),
		UserID:      id,
		DisplayName: "New User",
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
	}
}
