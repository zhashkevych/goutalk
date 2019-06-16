package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/goutalk/auth"
	"github.com/zhashkevych/goutalk/chat"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

const (
	httpHeaderAccessToken = "Authorization"
	ctxKeyUser            = "goutalkuser"
)

func AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader(httpHeaderAccessToken)
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &Response{
			"no auth header provided",
		})
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &Response{
			"invalid auth header structure",
		})
		return
	}

	if headerParts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &Response{
			"wrong auth header type",
		})
		return
	}

	claims, err := auth.VerifyAuthToken(headerParts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &Response{
			"access token is invalid",
		})
		return
	}

	user, err := toAuthUser(claims)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &Response{
			"failed to parse user data",
		})
		return
	}

	c.Set(ctxKeyUser, user)
}

func toAuthUser(c *auth.Claims) (*chat.User, error) {
	id, err := primitive.ObjectIDFromHex(c.UserID)
	if err != nil {
		return nil, err
	}

	return &chat.User{
		ID:       id,
		Username: c.Username,
	}, nil
}
