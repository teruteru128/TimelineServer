package models

import (
	"gopkg.in/mgo.v2/bson"

	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// UserResponse GET /users/:id のためのレスポンス
type UserResponse struct {
	ID          string          `json:"id"`                // 恒久ID (bson ObjectID)
	UserID      string          `json:"screen_name"`       // ユーザ名(@kitten)
	DisplayName string          `json:"name"`              // 表示名(Kitten)
	PostsCount  uint            `json:"posts_count"`       // 投稿総数(0-)
	Location    string          `json:"location"`          // 居住地(グンマー)
	Following   []bson.ObjectId `json:"friends"`           // フォローしているユーザの恒久ID一覧
	Followers   []bson.ObjectId `json:"followers"`         // フォローされているユーザの恒久ID一覧
	WebsiteURL  string          `json:"url"`               // ウェブサイトのURL(http://example.com)
	AvatarURL   string          `json:"profile_image_url"` // プロフィール画像(http://static_cdn/profile_images/0.png)
	Official    bool            `json:"official"`          // 公式
	Description string          `json:"description"`
	jwt.StandardClaims
}

// LoginSuccessResponse POST /auth が成功したときのレスポンス
type LoginSuccessResponse struct {
	ID           string          `json:"id"`                // 恒久ID (bson ObjectID)
	UserID       string          `json:"screen_name"`       // ユーザ名(@kitten)
	DisplayName  string          `json:"name"`              // 表示名(Kitten)
	PostsCount   uint            `json:"posts_count"`       // 投稿総数(0-)
	Location     string          `json:"location"`          // 居住地(グンマー)
	Following    []bson.ObjectId `json:"friends"`           // フォローしているユーザの恒久ID一覧
	Followers    []bson.ObjectId `json:"followers"`         // フォローされているユーザの恒久ID一覧
	WebsiteURL   string          `json:"url"`               // ウェブサイトのURL(http://example.com)
	AvatarURL    string          `json:"profile_image_url"` // プロフィール画像(http://static_cdn/profile_images/0.png)
	Official     bool            `json:"official"`          // 公式
	SessionToken string          `json:"session_token"`     // JWTセッショントークン(RS256_JWT_TOKEN)
}

// ErrorResponse リクエストの処理中にエラーが発生したときのレスポンス
type ErrorResponse struct {
	Error string `json:"error"`
}

// PostResponse 投稿レスポンスの構造体
type PostResponse struct {
	FavoritedIds    []bson.ObjectId `json:"favorited_ids"`
	CreatedAt       time.Time       `json:"created_at"`
	ID              bson.ObjectId   `json:"id"` // BSON ObjectID
	MentionsID      []bson.ObjectId `json:"mentions_id"`
	URLs            []string        `json:"urls"`
	Hashtags        []string        `json:"hashtags"`
	InReplyToUserID bson.ObjectId   `json:"in_reply_to_user_id"`
	Text            string          `json:"text"`
	Shared          []bson.ObjectId `json:"shared"`
	User            UserResponse    `json:"user"`
}

func PostToPostResponse(post Post, user User) PostResponse {
	return PostResponse{
		FavoritedIds:    post.FavoritedIds,
		CreatedAt:       post.CreatedAt,
		ID:              post.ID,
		MentionsID:      post.MentionsID,
		URLs:            post.URLs,
		Hashtags:        post.Hashtags,
		InReplyToUserID: post.InReplyToUserID,
		Text:            post.Text,
		Shared:          post.Shared,
		User:            UserToUserResponse(user),
	}
}

func PostsToPostResponseArray(posts []Post, users []User, sameUser bool) []PostResponse {
	var arr []PostResponse
	for i, post := range posts {
		if !sameUser {
			resp := PostToPostResponse(post, users[i])
			arr = append(arr, resp)
		} else {
			resp := PostToPostResponse(post, users[0])
			arr = append(arr, resp)
		}

	}

	return arr
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
		Description: user.Description,
	}
}

func UserToLoginSucessResponse(user User, token string) LoginSuccessResponse {
	return LoginSuccessResponse{
		ID:           user.ID.Hex(),
		UserID:       user.UserID,
		DisplayName:  user.DisplayName,
		PostsCount:   uint(len(user.Posts)),
		Location:     user.Location,
		Following:    user.Following,
		Followers:    user.Followers,
		WebsiteURL:   user.WebsiteURL,
		AvatarURL:    user.AvatarURL,
		Official:     user.Official,
		SessionToken: token,
	}
}

// UsersToUserResponseArray User配列をAPI用ユーザ配列構造体に変換する
func UsersToUserResponseArray(users []User) []UserResponse {
	var arr []UserResponse
	for _, user := range users {
		resp := UserToUserResponse(user)
		arr = append(arr, resp)
	}

	return arr
}
