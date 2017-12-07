package token

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/TinyKitten/TimelineServer/config"
	jwt "github.com/dgrijalva/jwt-go"
)

// JWTClaim JWTのクレーム
type JWTClaim struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

// CreateToken JWTトークンを生成する
func CreateToken(id bson.ObjectId, adminFlag bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["iss"] = "KittenTimeline"
	claims["admin"] = adminFlag
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	key := config.GetAPIConfig().Jwt

	signed, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return signed, nil
}
