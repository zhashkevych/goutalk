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


func (h *Handler) Authorize(c *gin.Context) {
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

	user, err := h.chatter.GetUserByID(c.Request.Context(), claims.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, &Response{
			"no user found",
		})
		return
	}

	c.Set(ctxKeyUser, user)
}
