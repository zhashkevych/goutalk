package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/goutalk/chat"
	"net/http"
)

type getUsersOutput struct {
	Users []*user `json:"users"`
}

type user struct {
	ID       string `json:"user_id"`
	Username string `json:"user_name"`
}

func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.chatter.GetAllUsers(c.Request.Context())
	if err != nil {
		// TODO: with response
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &getUsersOutput{
		Users: toUsers(users),
	})
}

func (h *Handler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.chatter.GetUserByID(c.Request.Context(), id)
	if err != nil {
		// TODO: 400 bad request
		c.AbortWithStatusJSON(http.StatusInternalServerError, &Response{
			"no user found",
		})
		return
	}

	c.JSON(http.StatusOK, toUser(user))
}

func toUsers(users []*chat.User) []*user {
	out := make([]*user, len(users))
	for i := range users {
		out[i] = toUser(users[i])
	}

	return out
}

func toUser(u *chat.User) *user {
	return &user{
		ID:       u.ID.Hex(),
		Username: u.Username,
	}
}
