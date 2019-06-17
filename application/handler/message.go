package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/goutalk/chat"
	"net/http"
)

type messageInput struct {
	RoomID  string `json:"room_id"`
	Message string `json:"message"`
}

func (h *Handler) SendMessage(c *gin.Context) {
	var inp messageInput
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &Response{
			"wrong input body",
		})
		return
	}

	user := c.MustGet(ctxKeyUser).(*chat.User)
	m := &chat.Message{
		UserID:  user.ID.Hex(),
		RoomID:  inp.RoomID,
		Message: inp.Message,
	}

	if err := h.chatter.SendMessage(m); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &Response{
			"failed to broadcast message",
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		"message successfully sent",
	})
}
