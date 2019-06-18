package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var tests = []struct {
	userID   string
	username string
	password string
}{
	{"test_id_1", "user1", "pass"},
	{"test_id_2", "user2", "pass123"},
	{"test_id_3", "user3", "qwerty123"},
}

func TestAuthTokenOK(t *testing.T) {
	for _, test := range tests {
		token, err := GenerateAuthToken(test.userID, test.username, test.password)
		assert.NoError(t, err)

		_, err = VerifyAuthToken(token)
		assert.NoError(t, err)
	}
}
