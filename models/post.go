package models

// Post 投稿の構造体
type Post struct {
	UserID      uint   // 投稿したユーザのID
	PostID      uint   // 一意な投稿ID
	Text        string `redis:"text"`
	CreatedDate string `redis:"createdDate"`
}
