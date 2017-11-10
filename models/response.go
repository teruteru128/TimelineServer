package models

import "time"

// UserResponse GET /users/:id のためのレスポンス
type UserResponse struct {
	ID             string `json:"id"`             // 恒久ID (bson ObjectID)
	UserID         string `json:"userID"`         // ユーザ名(@kitten)
	DisplayName    string `json:"displayName"`    // 表示名(Kitten)
	PostsCount     uint   `json:"postsCount"`     // 投稿総数(0-)
	Location       string `json:"location"`       // 居住地(グンマー)
	FollowingCount uint   `json:"followingCount"` // フォローしている数(0-)
	FollowersCount uint   `json:"followersCount"` // フォローされている数(0-)
	WebsiteURL     string `json:"websiteUrl"`     // ウェブサイトのURL(http://example.com)
	AvatarURL      string `json:"avatarUrl"`      // プロフィール画像(http://static_cdn/profile_images/0.png)
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

func UserToUserResponse(user User) UserResponse {
	return UserResponse{
		ID:             user.ID.Hex(),
		UserID:         user.UserID,
		DisplayName:    user.DisplayName,
		PostsCount:     uint(len(user.Posts)),
		Location:       user.Location,
		FollowingCount: uint(len(user.Following)),
		FollowersCount: uint(len(user.Followers)),
		WebsiteURL:     user.WebsiteURL,
		AvatarURL:      user.AvatarURL,
	}
}
