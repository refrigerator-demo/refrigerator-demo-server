package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

var JwtSecretKey = []byte(os.Getenv("JWT_SECRET"))

type claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(id uint) (string, error) {
	return generateToken(id, time.Now())
}

func generateToken(userId uint, now time.Time) (string, error) {
	claims := &claims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: now.Add(time.Hour * 36).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES512, claims)
	accessToken, err := token.SignedString(JwtSecretKey)

	if nil != err {
		return "", err
	}

	return accessToken, err
}

func ExtractUserId(ctx context.Context) (uint, error) {
	tokenString, err := grpc_auth.AuthFromMD(ctx, "Token")
	if nil != err {
		return 0, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecretKey, nil
	})

	if false == token.Valid {

	}

	c, ok := token.Claims.(*claims)
	if false == ok {
		return 0, errors.New("invalid token: cannot map token to claims")
	}

	if c.ExpiresAt < time.Now().Unix() {
		return 0, errors.New("token expired")
	}

	return c.UserID, nil
}
