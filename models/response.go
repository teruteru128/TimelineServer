package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	jwt "github.com/dgrijalva/jwt-go"
)

// UserResponse GET /users/:id のためのレスポンス
type UserResponse struct {
	ID          string          `json:"id"`          // 恒久ID (bson ObjectID)
	UserID      string          `json:"userId"`      // ユーザ名(@kitten)
	DisplayName string          `json:"displayName"` // 表示名(Kitten)
	PostsCount  uint            `json:"postsCount"`  // 投稿総数(0-)
	Location    string          `json:"location"`    // 居住地(グンマー)
	Following   []bson.ObjectId `json:"following"`   // フォローしているユーザの恒久ID一覧
	Followers   []bson.ObjectId `json:"followers"`   // フォローされているユーザの恒久ID一覧
	WebsiteURL  string          `json:"websiteUrl"`  // ウェブサイトのURL(http://example.com)
	AvatarURL   string          `json:"avatarUrl"`   // プロフィール画像(http://static_cdn/profile_images/0.png)
	Official    bool            `json:"official"`    // 公式
	jwt.StandardClaims
}

// LoginSuccessResponse POST /auth が成功したときのレスポンス
type LoginSuccessResponse struct {
	ID           string    `json:"id"`           // 恒久ID (bson ObjectID)
	UserID       string    `json:"userId"`       // ユーザ名(@kitten)
	CreatedDate  time.Time `json:"createdDate"`  // ユーザ登録日時(2017-08-28T07:46:09.801Z)
	UpdatedDate  time.Time `json:"updatedDate"`  // 最終更新日(2017-09-28T07:46:09.801Z)
	SessionToken string    `json:"sessionToken"` // JWTセッショントークン(RS256_JWT_TOKEN)
}

// ErrorResponse リクエストの処理中にエラーが発生したときのレスポンス
type ErrorResponse struct {
	Error string `json:"error"`
}

// UserToUserResponse UserをAPI用ユーザ構造体に変換する
func UserToUserResponse(user User) UserResponse {
	return UserResponse{
		ID:          user.ID.Hex(),
		UserID:      user.UserID,
		DisplayName: user.DisplayName,
		PostsCount:  uint(len(user.Posts)),
		Location:    user.Location,
		Following:   user.Following,
		Followers:   user.Followers,
		WebsiteURL:  user.WebsiteURL,
		AvatarURL:   user.AvatarURL,
		Official:    user.Official,
	}
}

// UsersToUserResponseArray User配列をAPI用ユーザ配列構造体に変換する
func UsersToUserResponseArray(users []User) []UserResponse {
	var arr []UserResponse
	for _, user := range users {
		resp := UserResponse{
			ID:          user.ID.Hex(),
			UserID:      user.UserID,
			DisplayName: user.DisplayName,
			PostsCount:  uint(len(user.Posts)),
			Location:    user.Location,
			Following:   user.Following,
			Followers:   user.Followers,
			WebsiteURL:  user.WebsiteURL,
			AvatarURL:   user.AvatarURL,
			Official:    user.Official,
		}
		arr = append(arr, resp)
	}

	return arr
}
