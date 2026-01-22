package tool

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Exp     int64  `json:"exp"`
	Iat     int64  `json:"iat"`
	Payload string `json:"payload"`
}

type Pl struct {
	Uid int64 `json:"uid"`
}

func GetJwtToken(secretKey string, iat, seconds, uid int64) (string, error) {
	var p = map[string]int64{
		"uid": uid,
	}

	payload, _ := json.Marshal(p)
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["payload"] = payload
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func VerifyTokenHS256(tokenString string, secretKey []byte) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
}

func GetJwtUuid(tokenString string) (int64, error) {
	jwtParts := strings.Split(tokenString, ".")

	if len(jwtParts) != 3 {
		return 0, errors.New("invalid token")
	}

	encodedPayload := jwtParts[1]

	dstPayload, _ := base64.RawStdEncoding.DecodeString(encodedPayload)
	var c Claims
	json.Unmarshal(dstPayload, &c)

	if c.Exp < time.Now().Unix() {
		return 0, errors.New("invalid token")
	}

	var p Pl
	dp, _ := base64.URLEncoding.DecodeString(c.Payload)

	json.Unmarshal(dp, &p)

	return p.Uid, nil
}
