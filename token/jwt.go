package token

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/TinyKitten/TimelineServer/config"
	jwt "github.com/dgrijalva/jwt-go"
)

type JwtClaim struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

func CreateToken(id bson.ObjectId) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["iss"] = "KittenTimeline"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	key := config.GetAPIConfig().Jwt

	signed, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return signed, nil
}
