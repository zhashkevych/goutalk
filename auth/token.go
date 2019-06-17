package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"time"
)

const salt = "i1MjGcbuMX7DFEDL"
const tokenLifetime = time.Hour * 24

type Claims struct {
	jwt.StandardClaims
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func GenerateAuthToken(userID string, username, password string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(tokenLifetime).Unix(),
		},
		UserID:   userID,
		Username: username,
		Password: password,
	})

	return token.SignedString([]byte(salt))
}

// VerifyAuthToken returns masked profileID and email hash
func VerifyAuthToken(token string) (*Claims, error) {
	c := new(Claims)
	t, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("invalid signing method, expected HMAC, got: %s", t.Method.Alg())
		}

		return []byte(salt), nil
	})

	if err != nil {
		return nil, err
	}

	c, ok := t.Claims.(*Claims)
	if !ok {
		return nil, errors.New("failed to parse token")
	}

	return c, nil
}