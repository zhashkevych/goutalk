package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"time"
)

const salt = "qwertyuiop1234567890"

type Claims struct {
	jwt.StandardClaims
	UserID   uuid.UUID `json:"user_id"`
	Name     string    `json:"name"`
	Nickname string    `json:"nickname"`
	Password string    `json:"password"`
}

func GenerateAuthToken(userID uuid.UUID, name, nickname, password string, lifetime time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(lifetime).Unix(),
		},
		UserID:   userID,
		Name:     name,
		Nickname: nickname,
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
