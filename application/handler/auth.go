package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/goutalk/auth"
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
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if headerParts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := auth.VerifyAuthToken(headerParts[1])
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set(ctxKeyUser, user)
}